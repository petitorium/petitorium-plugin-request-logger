// Package main provides a request logger plugin for Petitorium.
// This plugin logs outgoing requests and incoming responses with support for both
// raw template variables and expanded environment variables.
//
// To build: go build -buildmode=plugin -o request-logger.so .
package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/go-homedir"
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

// HookFuncs returns the hook functions
func (rl *RequestLogger) HookFuncs() map[types.HookType]types.PluginHook {
	return map[types.HookType]types.PluginHook{
		types.PreRequest:               rl.logRawRequest,
		types.PostVariableSubstitution: rl.logExpandedRequest,
		types.PreSend:                  rl.logFinalRequest,
		types.PostReceive:              rl.logResponse,
	}
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

// logResponse logs the incoming response
func (rl *RequestLogger) logResponse(ctx *types.HookContext) error {
	logEntry := fmt.Sprintf("[%s] [RESPONSE]\n", time.Now().Format(time.RFC3339))

	// Try to get response details if available
	if ctx.Response != nil {
		// Use reflection to safely access response fields
		respVal := reflect.ValueOf(ctx.Response)
		if respVal.Kind() == reflect.Ptr {
			respVal = respVal.Elem()
		}

		if respVal.IsValid() && respVal.Kind() == reflect.Struct {
			// Try to get StatusCode field
			if statusField := respVal.FieldByName("StatusCode"); statusField.IsValid() {
				if statusField.Kind() == reflect.Int {
					logEntry += fmt.Sprintf("Status: %d\n", statusField.Int())
				}
			}

			// Try to get Status field
			if statusField := respVal.FieldByName("Status"); statusField.IsValid() {
				if statusField.Kind() == reflect.String {
					logEntry += fmt.Sprintf("Status Text: %s\n", statusField.String())
				}
			}

			// Try to get Body field
			if bodyField := respVal.FieldByName("Body"); bodyField.IsValid() {
				if bodyField.Kind() == reflect.String {
					body := bodyField.String()
					if body != "" {
						logEntry += fmt.Sprintf("Body:\n%s\n", body)
					}
				}
			}
		}
	} else {
		logEntry += "No response data available\n"
	}

	logEntry += "\n"
	return rl.writeToLog(logEntry, ctx)
}

// writeToLog writes a log entry to the log file
func (rl *RequestLogger) writeToLog(entry string, ctx *types.HookContext) error {
	// Get log file path from context config, default to "petitorium.log" if not specified
	logFile := "petitorium.log"
	if ctx.Config != nil {
		if configLogFile, ok := ctx.Config["logFile"].(string); ok && configLogFile != "" {
			logFile = configLogFile
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
