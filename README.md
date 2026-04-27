# Petitorium Request Logger Plugin

A comprehensive request and response logging plugin for [Petitorium](https://github.com/petitorium/petitorium) that provides detailed logging of HTTP requests and responses with support for both raw template variables and expanded environment variables.

## Features

- **Dual Logging**: Logs both raw template variables (`{{protocol}}{{domain}}`) and expanded environment variables (`https://api.example.com`)
- **Request Logging**: Captures HTTP method, URL, headers, and body
- **Response Logging**: Records status codes, status text, and response body
- **Configurable Output**: Customizable log file location
- **Timestamped Entries**: All log entries include RFC3339 timestamps
- **Thread-Safe**: Safe for concurrent request processing

## Installation

### Prerequisites

- Go 1.23 or later
- Petitorium with plugin support

### Building from Source

1. Clone this repository:

```bash
git clone https://github.com/petitorium/petitorium-plugin-request-logger.git
cd petitorium-plugin-request-logger
```

2. Build the plugin:

```bash
go build -o request-logger .
# For cross-compilation:
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o request-logger .
```

3. Copy the plugin to your Petitorium plugins directory:

```bash
cp request-logger ~/.config/petitorium/plugins/
```

## Configuration

Add the plugin to your Petitorium configuration file (`~/.config/petitorium/config.yaml`):

```yaml
plugins:
  enabled:
    - request-logger
  config:
    request-logger:
      log_file: ~/petitorium.log # Optional: global log file (defaults to ~/petitorium.log)
      workspaces: # Optional: workspace-specific overrides
        "My Workspace Name":
          log_file: ~/my-workspace.log
```

### Configuration Options

| Option       | Type   | Default            | Description                                      |
| ------------ | ------ | ------------------ | ------------------------------------------------ |
| `log_file`   | string | `~/petitorium.log` | Path to the log file (supports `~` expansion)    |
| `workspaces` | object | `null`             | Map of workspace names to specific overrides     |

## Usage

Once installed and configured, the plugin automatically logs all HTTP requests and responses:

### Example Log Output

```
[2025-11-23T10:30:45Z] [REQUEST] [RAW] GET {{protocol}}{{domain}}/users/1
Headers:
  Authorization: Bearer {{token}}
  Content-Type: application/json

[2025-11-23T10:30:45Z] [REQUEST] [EXPANDED] GET https://api.example.com/users/1
Headers:
  Authorization: Bearer abc123def456
  Content-Type: application/json

[2025-11-23T10:30:45Z] [RESPONSE]
Status: 200
Status Text: OK
Body:
{"id": 1, "name": "John Doe", "email": "john@example.com"}
```

## Plugin Hooks

This plugin implements the following Petitorium hooks:

- `PreRequest`: Logs requests with raw template variables
- `PostVariableSubstitution`: Logs requests with expanded environment variables
- `PostReceive`: Logs responses

## Development

### Project Structure

```
petitorium-plugin-request-logger/
├── request-logger.go       # Main plugin implementation
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
└── README.md               # This file
```

### Building for Development

```bash
# Download dependencies
go mod tidy

# Build the plugin executable
go build -o request-logger .

# Run tests (if any)
go test ./...
```

### Plugin Architecture

The plugin follows Petitorium's plugin architecture (using HashiCorp go-plugin):

1. **SDK Imports**: Uses `github.com/petitorium/petitorium-plugin-sdk`
2. **Main Plugin**: Implements the `types.Plugin` interface
3. **Hook Functions**: Process requests/responses at different stages via `ExecuteHook`
4. **gRPC Server**: Runs a local gRPC server via `plugin.Serve()` for communication with Petitorium

## Troubleshooting

### Plugin Not Loading

1. Verify the plugin executable exists: `ls -la ~/.config/petitorium/plugins/request-logger`
2. Ensure the file has execution permissions: `chmod +x ~/.config/petitorium/plugins/request-logger`
2. Check Petitorium logs for plugin loading errors
3. Ensure the plugin is enabled in your configuration

### Log File Not Created

1. Check file permissions in the target directory
2. Verify the log file path in your configuration
3. Ensure the directory exists (the plugin won't create parent directories)

### Missing Log Entries

1. Verify the plugin is properly loaded: `petitorium --list-plugins`
2. Check that requests are actually being sent through Petitorium
3. Look for errors in Petitorium's main log

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and test thoroughly
4. Commit your changes: `git commit -am 'Add new feature'`
5. Push to the branch: `git push origin feature-name`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Petitorium](https://github.com/petitorium/petitorium) - The main API client application
- [Petitorium Auth Injector Plugin](https://github.com/petitorium/petitorium-plugin-auth-injector) - An authentication injection plugin for [Petitorium](https://github.com/petitorium/petitorium) that automatically injects authentication headers into HTTP requests and captures authentication tokens from responses.

## Support

For issues and questions:

- Create an issue in this repository
- Check the [Petitorium documentation](https://github.com/petitorium/petitorium)
- Review existing issues for similar problems
