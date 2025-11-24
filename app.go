package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ConfigFileInfo holds config file metadata
type ConfigFileInfo struct {
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	Exists       bool   `json:"exists"`
}

// GetConfigFileInfo returns file metadata for the given path
func (a *App) GetConfigFileInfo(path string) ConfigFileInfo {
	info := ConfigFileInfo{Path: path, Exists: false}

	stat, err := os.Stat(path)
	if err != nil {
		return info
	}

	info.Exists = true
	info.Size = stat.Size()
	info.LastModified = stat.ModTime().Format(time.RFC3339)
	return info
}

// ReadConfigFile reads and returns the content of a config file
func (a *App) ReadConfigFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// DownloadConfigResult holds the result of downloading a config
type DownloadConfigResult struct {
	TempPath string `json:"tempPath"`
	Size     int64  `json:"size"`
	Content  string `json:"content"`
}

// DownloadConfig downloads config from URL and returns temp file path and content
func (a *App) DownloadConfig(url string) (*DownloadConfigResult, error) {
	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", "zashboard-config-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() { _ = tempFile.Close() }()

	if _, err := tempFile.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	return &DownloadConfigResult{
		TempPath: tempFile.Name(),
		Size:     int64(len(content)),
		Content:  string(content),
	}, nil
}

// ValidateConfig is a placeholder for config validation
// Returns true if valid, error message if invalid
// This can be extended to perform actual validation
func (a *App) ValidateConfig(content string) (bool, string) {
	if content == "" {
		return false, "config content is empty"
	}
	// Placeholder: add actual validation logic here
	// e.g., YAML parsing, schema validation, etc.
	return true, ""
}

// UpdateConfigFile replaces the target config file with content from temp file
// Only proceeds if validation passes
func (a *App) UpdateConfigFile(tempPath string, targetPath string) error {
	// Read temp file content
	content, err := os.ReadFile(tempPath)
	if err != nil {
		return fmt.Errorf("failed to read temp file: %w", err)
	}

	// Validate config
	valid, errMsg := a.ValidateConfig(string(content))
	if !valid {
		return fmt.Errorf("validation failed: %s", errMsg)
	}

	// Write to target path
	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Clean up temp file
	_ = os.Remove(tempPath)

	return nil
}
