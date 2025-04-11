package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Version 变量将在构建时通过 -ldflags 设置
var Version = "development"

type MockAPI struct {
	Method   string
	Path     string
	Response string
}

func parseMockFiles(directory string) ([]MockAPI, error) {
	var mockAPIs []MockAPI

	// Ensure mock directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", directory)
	}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".api" {
			log.Printf("Parsing file: %s", path)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Create a scanner with increased buffer size (10MB)
			scanner := bufio.NewScanner(file)
			const maxCapacity = 10 * 1024 * 1024 // 10MB
			buf := make([]byte, maxCapacity)
			scanner.Buffer(buf, maxCapacity)

			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
					continue
				}

				parts := strings.SplitN(line, " ", 3)
				if len(parts) < 3 {
					log.Printf("Warning: %s:%d - Invalid line format: %s", path, lineNum, line)
					continue
				}

				method := strings.TrimSpace(parts[0])
				path := strings.TrimSpace(parts[1])
				response := strings.TrimSpace(parts[2])

				// Check and remove quotes from response
				if strings.HasPrefix(response, "'") && strings.HasSuffix(response, "'") {
					response = response[1 : len(response)-1]
				} else if strings.HasPrefix(response, "\"") && strings.HasSuffix(response, "\"") {
					response = response[1 : len(response)-1]
				} else if strings.HasPrefix(response, "`") && strings.HasSuffix(response, "`") {
					response = response[1 : len(response)-1]
				}

				mockAPIs = append(mockAPIs, MockAPI{
					Method:   method,
					Path:     path,
					Response: response,
				})
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error scanning file %s: %v", path, err)
			}
		}
		return nil
	})

	return mockAPIs, err
}

// Fix JSON response by adding quotes to keys and values
func fixJSONResponse(jsonStr string) string {
	// Use regex to find keys without quotes
	reKey := regexp.MustCompile(`([{,])(\s*)([a-zA-Z0-9_]+)(\s*):`)
	jsonStr = reKey.ReplaceAllString(jsonStr, `$1$2"$3"$4:`)

	// Find values without quotes (non-numeric)
	reValue := regexp.MustCompile(`:(\s*)([a-zA-Z][a-zA-Z0-9_]*)(\s*[,}])`)
	jsonStr = reValue.ReplaceAllString(jsonStr, `:$1"$2"$3`)

	return jsonStr
}

func setupMockServer(mockAPIs []MockAPI) {
	for _, api := range mockAPIs {
		// Use closure to preserve response for each route
		func(method, path, response string) {
			http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					http.Error(w, fmt.Sprintf("Method not allowed, expected: %s", method), http.StatusMethodNotAllowed)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				// Fix JSON format
				validJSON := fixJSONResponse(response)

				// Validate JSON
				if !json.Valid([]byte(validJSON)) {
					log.Printf("Error: Cannot fix invalid JSON: %s", validJSON)
					http.Error(w, "Server error: Invalid JSON response", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, validJSON)
				log.Printf("Request successful: %s %s", method, path)
			})
		}(api.Method, api.Path, api.Response)
	}
}

// 尝试从多个目录加载mock文件
func tryLoadMockFiles(directories []string) ([]MockAPI, string, error) {
	var lastErr error
	var mockAPIs []MockAPI
	var usedDir string

	for _, dir := range directories {
		log.Printf("Trying to load mock files from: %s", dir)
		apis, err := parseMockFiles(dir)
		if err != nil {
			log.Printf("Failed to load from %s: %v", dir, err)
			lastErr = err
			continue
		}

		if len(apis) > 0 {
			log.Printf("Successfully loaded mock files from: %s", dir)
			mockAPIs = apis
			usedDir = dir
			break
		} else {
			log.Printf("No mock files found in: %s", dir)
		}
	}

	if len(mockAPIs) == 0 {
		if lastErr != nil {
			return nil, "", fmt.Errorf("failed to load mock files from any directory: %v", lastErr)
		}
		return nil, "", fmt.Errorf("no mock files found in any of the directories: %v", directories)
	}

	return mockAPIs, usedDir, nil
}

func main() {
	// 定义命令行参数
	mockDir := flag.String("mock", "", "Directory containing mock API definitions")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Go API Mock Server v%s...", Version)
	log.Printf("Port: %d", *port)

	var mockAPIs []MockAPI
	var usedDir string
	var err error

	// 如果指定了mock目录，直接使用
	if *mockDir != "" {
		log.Printf("Using specified mock directory: %s", *mockDir)
		mockAPIs, err = parseMockFiles(*mockDir)
		usedDir = *mockDir
	} else {
		// 否则尝试多个目录
		log.Printf("No mock directory specified, trying multiple locations...")
		directories := []string{"mock", "../mock"}
		mockAPIs, usedDir, err = tryLoadMockFiles(directories)
	}

	if err != nil {
		log.Fatalf("Error parsing mock files: %v", err)
	}

	if len(mockAPIs) == 0 {
		log.Fatal("No mock API configurations found, please check if there are .api files in the mock directory")
	}

	log.Printf("Successfully loaded %d mock APIs from %s", len(mockAPIs), usedDir)
	for _, api := range mockAPIs {
		log.Printf("Registered API: %s %s -> %s", api.Method, api.Path, api.Response)
	}

	setupMockServer(mockAPIs)

	log.Printf("HTTP server started, listening on port %d", *port)
	log.Printf("APIs can be accessed via http://localhost:%d/", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
