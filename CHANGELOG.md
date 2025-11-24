# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-23

### Added

- Initial release of the Petitorium Request Logger Plugin
- Dual logging support for both raw template variables and expanded environment variables
- Request logging with method, URL, headers, and body
- Response logging with status code, status text, and response body
- Configurable log file location with home directory expansion
- Support for all major platforms (Linux, macOS, Windows) and architectures (AMD64, ARM64)
- Comprehensive documentation and examples
- GitHub Actions workflows for CI/CD and releases
- MIT License

### Features

- **PreRequest Hook**: Logs outgoing requests with raw template variables (e.g., `{{protocol}}{{domain}}`)
- **PostVariableSubstitution Hook**: Logs outgoing requests with expanded environment variables (e.g., `https://api.example.com`)
- **PostReceive Hook**: Logs incoming responses with status and body information
- **Configurable Output**: Customizable log file path via plugin configuration
- **Thread-Safe**: Safe for concurrent request processing
- **RFC3339 Timestamps**: All log entries include standardized timestamps

### Technical Details

- Built with Go 1.23
- Uses reflection for safe response field access
- Supports all Petitorium hook types for future extensibility
- Standalone types package for independence from Petitorium core
- Comprehensive error handling and logging
