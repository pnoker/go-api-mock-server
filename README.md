# Go API Mock Server

A simple yet powerful Go API Mock server that automatically reads `.api` files from the `mock` directory and creates HTTP API endpoints.

## Features

- Automatically loads API definitions from the `mock` directory
- Supports multiple response formats (single quotes, double quotes, backticks)
- Automatically fixes JSON format (e.g., unquoted `ok` values)
- Supports custom port and mock directory location
- Auto-detects mock directory location
- Detailed logging output
- Supports Windows and macOS platforms

## Installation

### Method 1: Download Executable

Download the appropriate executable for your platform from the [releases page](https://github.com/pnoker/go-api-mock-server/releases).

### Method 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/pnoker/go-api-mock-server.git
cd go-api-mock-server

# Build executables
./build.sh
```

The build script will create executables for the following platforms:
- Windows (64-bit)
- macOS (Intel/AMD 64-bit)
- macOS (Apple Silicon/ARM64)

## Usage

### 1. Create API Definition Files

Create `.api` files in the `mock` directory with the following format:

```
METHOD PATH 'RESPONSE'
```

Three quote formats are supported:
```
GET api/v3/device/list '{"code":"ok","message":"xxx"}'
POST api/v3/user/login "{"code":"ok","message":"xxx"}"
POST api/v3/user/logout `{"code":ok,"message":"xxx"}`
```

### 2. Run the Server

#### Basic Usage

```bash
# Auto-detect mock directory
./go-api-mock-server-windows-amd64.exe

# Specify mock directory
./go-api-mock-server-windows-amd64.exe -mock="../mock"

# Specify port
./go-api-mock-server-windows-amd64.exe -port=9090

# Specify both directory and port
./go-api-mock-server-windows-amd64.exe -mock="../mock" -port=9090
```

#### Command Line Arguments

- `-mock`: Directory containing mock API definitions (default: auto-detect)
- `-port`: Port to listen on (default: 8080)

### 3. Access APIs

Once the server is running, you can access the configured API endpoints via HTTP requests:

```bash
# Example request
curl http://localhost:8080/api/v3/device/list
```

Response example:
```json
{"code":"ok","message":"xxx"}
```

## Auto-detection of Mock Directory

If you don't specify the `-mock` parameter, the server will automatically try the following locations (in order):
1. `./mock` (current directory)
2. `../mock` (parent directory)

This eliminates the need to specify the mock directory location every time.

## Important Notes

- Responses must be valid JSON format
- If the response contains unquoted `ok` values, the server will automatically add quotes
- The server logs all API access
- Comments are supported: lines starting with `#` or `//` are ignored

## Examples

### API Definition File (mock/example.api)

```
# This is a comment
GET api/v3/device/list `{"code":ok,"message":"Device List"}`

# This is another comment
POST api/v3/user/login "{"code":"ok","message":"Login successful","data":{"token":"abc123"}}"
```

### Accessing APIs

```bash
# Get device list
curl http://localhost:8080/api/v3/device/list

# User login
curl -X POST http://localhost:8080/api/v3/user/login
```

## Troubleshooting

If you encounter issues, check the following:

1. Does the mock directory exist and contain `.api` files?
2. Is the `.api` file format correct?
3. Is the JSON response format valid?
4. Is the port already in use by another application?