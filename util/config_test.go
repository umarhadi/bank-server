package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// create a temporary directory
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// prepare sample config content
	content := strings.Join([]string{
		"ENVIRONMENT=development",
		"DB_SOURCE=postgres://user:pass@localhost:5432/dbname",
		"MIGRATION_URL=file://migrations",
		"REDIS_ADDRESS=localhost:6379",
		"HTTP_SERVER_ADDRESS=:8080",
		"GRPC_SERVER_ADDRESS=:9090",
		"TOKEN_SYMMETRIC_KEY=mysecrettoken",
		"ACCESS_TOKEN_DURATION=15m",
		"REFRESH_TOKEN_DURATION=30m",
		"EMAIL_SENDER_NAME=Bank",
		"EMAIL_SENDER_ADDRESS=bank@example.com",
		"EMAIL_SENDER_PASSWORD=password123",
	}, "\n")

	// write the config file (app.env) into the temporary directory
	configPath := filepath.Join(tmpDir, "app.env")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// load configuration from the temporary directory
	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// assertions
	if cfg.Environment != "development" {
		t.Errorf("expected Environment 'development', got '%s'", cfg.Environment)
	}
	if cfg.DBSource != "postgres://user:pass@localhost:5432/dbname" {
		t.Errorf("unexpected DBSource value")
	}
	if cfg.MigrationURL != "file://migrations" {
		t.Errorf("unexpected MigrationURL value")
	}
	if cfg.RedisAddress != "localhost:6379" {
		t.Errorf("unexpected RedisAddress value")
	}
	if cfg.HTTPServerAddress != ":8080" {
		t.Errorf("unexpected HTTPServerAddress value")
	}
	if cfg.GRPCServerAddress != ":9090" {
		t.Errorf("unexpected GRPCServerAddress value")
	}
	if cfg.TokenSymmetricKey != "mysecrettoken" {
		t.Errorf("unexpected TokenSymmetricKey value")
	}
	// parse durations to compare
	expectedAccess, _ := time.ParseDuration("15m")
	if cfg.AccessTokenDuration != expectedAccess {
		t.Errorf("expected AccessTokenDuration %v, got %v", expectedAccess, cfg.AccessTokenDuration)
	}
	expectedRefresh, _ := time.ParseDuration("30m")
	if cfg.RefreshTokenDuration != expectedRefresh {
		t.Errorf("expected RefreshTokenDuration %v, got %v", expectedRefresh, cfg.RefreshTokenDuration)
	}
	if cfg.EmailSenderName != "Bank" {
		t.Errorf("unexpected EmailSenderName value")
	}
	if cfg.EmailSenderAddress != "bank@example.com" {
		t.Errorf("unexpected EmailSenderAddress value")
	}
	if cfg.EmailSenderPassword != "password123" {
		t.Errorf("unexpected EmailSenderPassword value")
	}
}

func TestLoadConfig_FileMissing(t *testing.T) {
	// create a temporary directory without creating a config file
	tmpDir, err := os.MkdirTemp("", "config_missing")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// call LoadConfig which should fail due to missing config file
	_, err = LoadConfig(tmpDir)
	if err == nil {
		t.Errorf("expected error when config file is missing, got nil")
	}
}
