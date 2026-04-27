// Package main provides a request logger plugin for Petitorium.
// This plugin logs outgoing requests and incoming responses with support for both
// raw template variables and expanded environment variables.
//
// To build: go build -o request-logger .
// For cross-compilation: GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o request-logger .
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-homedir"
	"github.com/petitorium/petitorium-plugin-sdk/shared"
	"github.com/petitorium/petitorium-plugin-sdk/types"
)

// RequestLogger is a plugin that logs requests and responses
type RequestLogger struct{}

// Name returns the plugin name
func (rl *RequestLogger) Name() string {
	return "request-logger"
}

// Version returns the plugin version
func (rl *RequestLogger) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (rl *RequestLogger) Description() string {
	return "Logs outgoing requests and incoming responses with raw and expanded variable support"
}

// Hooks returns the hook types this plugin implements
func (rl *RequestLogger) Hooks() []types.HookType {
	return []types.HookType{types.PreRequest, types.PostVariableSubstitution, types.PreSend, types.PostReceive}
}

// logRawRequest logs the outgoing request with raw template variables
func (rl *RequestLogger) logRawRequest(ctx *types.HookContext) error {
	req := ctx.Request
	logEntry := fmt.Sprintf("[%s] [REQUEST] [RAW] %s %s\n", time.Now().Format(time.RFC3339), req.Method, req.URL)

	// Add headers if present
	if len(req.Headers) > 0 {
		logEntry += "Headers:\n"
		for key, value := range req.Headers {
			logEntry += fmt.Sprintf("  %s: %s\n", key, value)
		}
	}

	// Add body if present
	if req.Body != "" {
		logEntry += fmt.Sprintf("Body:\n%s\n", req.Body)
	}

	logEntry += "\n"
	return rl.writeToLog(logEntry, ctx)
}

// logExpandedRequest logs the outgoing request with expanded environment variables
func (rl *RequestLogger) logExpandedRequest(ctx *types.HookContext) error {
	req := ctx.Request
	logEntry := fmt.Sprintf("[%s] [REQUEST] [EXPANDED] %s %s\n", time.Now().Format(time.RFC3339), req.Method, req.URL)

	// Add headers if present
	if len(req.Headers) > 0 {
		logEntry += "Headers:\n"
		for key, value := range req.Headers {
			logEntry += fmt.Sprintf("  %s: %s\n", key, value)
		}
	}

	// Add body if present
	if req.Body != "" {
		logEntry += fmt.Sprintf("Body:\n%s\n", req.Body)
	}

	logEntry += "\n"
	return rl.writeToLog(logEntry, ctx)
}

// logFinalRequest logs the final request with all modifications (including auth headers)
func (rl *RequestLogger) logFinalRequest(ctx *types.HookContext) error {
	req := ctx.Request
	logEntry := fmt.Sprintf("[%s] [REQUEST] [FINAL] %s %s\n", time.Now().Format(time.RFC3339), req.Method, req.URL)

	// Add headers if present (this will include injected auth headers)
	if len(req.Headers) > 0 {
		logEntry += "Headers (including injected):\n"
		for key, value := range req.Headers {
			logEntry += fmt.Sprintf("  %s: %s\n", key, value)
		}
	} else {
		logEntry += "No headers found\n"
	}

	// Add body if present
	if req.Body != "" {
		logEntry += fmt.Sprintf("Body:\n%s\n", req.Body)
	} else {
		logEntry += "No body found\n"
	}

	logEntry += "\n"
	return rl.writeToLog(logEntry, ctx)
}

// ExecuteHook executes a specific hook with the given context.
func (rl *RequestLogger) ExecuteHook(hookType types.HookType, ctx *types.HookContext) (*types.HookContext, error) {
	var err error
	switch hookType {
	case types.PreRequest:
		err = rl.logRawRequest(ctx)
	case types.PostVariableSubstitution:
		err = rl.logExpandedRequest(ctx)
	case types.PreSend:
		err = rl.logFinalRequest(ctx)
	case types.PostReceive:
		err = rl.logResponse(ctx)
	}
	return ctx, err
}

// logResponse logs the incoming response
func (rl *RequestLogger) logResponse(ctx *types.HookContext) error {
	logEntry := fmt.Sprintf("[%s] [RESPONSE]\n", time.Now().Format(time.RFC3339))

	// Try to get response details if available
	if ctx.Response != nil {
		logEntry += fmt.Sprintf("Status: %d\n", ctx.Response.StatusCode)
		logEntry += fmt.Sprintf("Status Text: %s\n", ctx.Response.Status)

		if len(ctx.Response.Headers) > 0 {
			logEntry += "Headers:\n"
			for key, values := range ctx.Response.Headers {
				logEntry += fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", "))
			}
		}

		if ctx.Response.Body != "" {
			logEntry += fmt.Sprintf("Body:\n%s\n", ctx.Response.Body)
		}
	} else {
		logEntry += "No response data available\n"
	}

	logEntry += "\n"
	return rl.writeToLog(logEntry, ctx)
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"request-logger": &shared.PetitoriumPlugin{Impl: &RequestLogger{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// writeToLog writes a log entry to the log file
func (rl *RequestLogger) writeToLog(entry string, ctx *types.HookContext) error {
	// Get log file path from context config, default to "petitorium.log" if not specified
	logFile := "petitorium.log"
	if ctx.Config != nil {
		if pluginConfig, ok := ctx.Config["request-logger"].(map[string]interface{}); ok {
			// Global fallback
			if configLogFile, ok := pluginConfig["log_file"].(string); ok && configLogFile != "" {
				logFile = configLogFile
			}

			// Workspace specific overrides
			if ctx.Workspace != "" {
				if workspacesConfig, ok := pluginConfig["workspaces"].(map[string]interface{}); ok {
					if wsConfig, ok := workspacesConfig[ctx.Workspace].(map[string]interface{}); ok {
						if configLogFile, ok := wsConfig["log_file"].(string); ok && configLogFile != "" {
							logFile = configLogFile
						}
					}
				}
			}
		}
	}

	// Expand home directory if path contains ~
	if expandedPath, err := homedir.Expand(logFile); err == nil {
		logFile = expandedPath
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

// Plugin is the exported plugin instance
var Plugin types.Plugin = &RequestLogger{}
