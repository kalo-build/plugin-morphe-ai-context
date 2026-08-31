package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morphe-ai-context/pkg/compile"
)

type PluginConfig struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	Verbose    bool   `json:"verbose,omitempty"`
}

const (
	ErrMissingConfig      = 3
	ErrInvalidConfig      = 4
	ErrInputPathRequired  = 12
	ErrOutputPathRequired = 13
	ErrCompileFailed      = 1
)

func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-ai-context <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON with inputPath and outputPath")
		os.Exit(ErrMissingConfig)
	}

	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(os.Args[1]), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ErrInvalidConfig)
	}

	if pluginConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: inputPath is required")
		os.Exit(ErrInputPathRequired)
	}
	if pluginConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: outputPath is required")
		os.Exit(ErrOutputPathRequired)
	}

	if abs, err := filepath.Abs(pluginConfig.InputPath); err == nil {
		pluginConfig.InputPath = abs
	}
	if abs, err := filepath.Abs(pluginConfig.OutputPath); err == nil {
		pluginConfig.OutputPath = abs
	}

	logInfo(pluginConfig.Verbose, "Reading Morphe registry from: '%s'", pluginConfig.InputPath)
	logInfo(pluginConfig.Verbose, "Writing AI context to: '%s'", pluginConfig.OutputPath)

	config := compile.DefaultCompileConfig(pluginConfig.InputPath, pluginConfig.OutputPath)

	if err := compile.MorpheToAIContext(config); err != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", err)
		os.Exit(ErrCompileFailed)
	}

	logInfo(pluginConfig.Verbose, "AI context generation completed successfully")
}
