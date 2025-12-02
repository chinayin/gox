// Package config provides unified configuration loading based on Viper.
//
// This package is designed as a reusable component for loading YAML
// configuration files with automatic default value injection, environment
// variable substitution, and local configuration overrides.
//
// # Basic Usage
//
// Create a loader and load configuration from a YAML file:
//
//	import "github.com/chinayin/gox/config"
//
//	loader := config.NewLoader()
//	var cfg AppConfig
//	if err := loader.Load("config.yaml", &cfg); err != nil {
//	    return fmt.Errorf("failed to load config: %w", err)
//	}
//
// # Configuration Structure
//
// Define your configuration using struct tags:
//
//	type AppConfig struct {
//	    Port     int    `default:"8080" validate:"required,min=1,max=65535"`
//	    LogLevel string `default:"info" validate:"oneof=debug info warn error"`
//	    Timeout  int    `default:"30"`
//	}
//
// # Default Values
//
// Two ways to set default values (in priority order from low to high):
//
// 1. Using struct tags (recommended for simple cases):
//
//	type Config struct {
//	    Port     int    `default:"8080"`
//	    LogLevel string `default:"info"`
//	}
//
// 2. Implementing Defaultable interface (for complex cases):
//
//	func (c *AppConfig) SetDefaults(set config.DefaultOption) {
//	    // Override struct tag defaults
//	    set("port", 9090)
//
//	    // Dynamic defaults
//	    set("cache_expiry", time.Hour * 24)
//
//	    // Environment-specific defaults
//	    if os.Getenv("ENV") == "production" {
//	        set("database.host", "prod-db.example.com")
//	    } else {
//	        set("database.host", "localhost")
//	    }
//	}
//
// Default value priority: struct tag < SetDefaults < config file < environment variables
//
// # Configuration Validation
//
// Implement the Validatable interface to enable automatic validation:
//
//	import "github.com/chinayin/gox/validator"
//
//	func (c *AppConfig) Validate() error {
//	    return validator.Validate(c)
//	}
//
// The Validate() method will be automatically called after loading configuration.
//
// For validation rules and custom validators, see github.com/chinayin/gox/validator package.
//
// # Mapstructure Tags
//
// Following the project's viper configuration standards, only add mapstructure
// tags when field names don't match YAML keys:
//
//	type Config struct {
//	    Enabled     bool   // matches "enabled" automatically (no tag needed)
//	    MaxAttempts int    `mapstructure:"max_attempts"` // required for snake_case
//	}
//
// # Local Configuration Override
//
// The loader automatically merges .local.yaml files for local overrides:
//
//	config.yaml        # Main configuration (committed to git)
//	config.local.yaml  # Local overrides (gitignored)
//
// Example:
//
//	# config.yaml
//	database:
//	  host: production.example.com
//
//	# config.local.yaml (overrides for local development)
//	database:
//	  host: localhost
//
// # Environment Variables
//
// Environment variables are automatically read and can override configuration:
//
//	# Configuration key: app.log_level
//	# Environment variable: APP_LOG_LEVEL=debug
//
// The loader uses dot-to-underscore conversion (app.log_level → APP_LOG_LEVEL).
//
// To disable environment variable reading (e.g., to avoid conflicts):
//
//	loader := config.NewLoader(config.WithoutEnv())
//	if err := loader.Load("config.yaml", &cfg); err != nil {
//	    return err
//	}
//
// To set a custom environment variable prefix:
//
//	loader := config.NewLoader(config.WithEnvPrefix("MYAPP"))
//	// Configuration key: port
//	// Environment variable: MYAPP_PORT=8080
//
// # Loading Multiple Configurations
//
// Load all configuration files from a directory:
//
//	loader := config.NewLoader()
//	configs, err := loader.LoadDirectory("configs/suppliers", &SupplierConfig{})
//	if err != nil {
//	    return err
//	}
//
//	for _, cfg := range configs {
//	    supplier := cfg.(*SupplierConfig)
//	    // process supplier configuration
//	}
//
// # Complete Example
//
//	package main
//
//	import (
//	    "fmt"
//	    "os"
//	    "github.com/chinayin/gox/config"
//	    "github.com/chinayin/gox/validator"
//	)
//
//	type AppConfig struct {
//	    Port      int    `default:"8080" validate:"required,min=1,max=65535"`
//	    LogLevel  string `default:"info" validate:"oneof=debug info warn error"`
//	    NacosAddr string `validate:"required,nacos_addr"`
//	}
//
//	func (c *AppConfig) SetDefaults(set config.DefaultOption) {
//	    // Environment-specific defaults
//	    if os.Getenv("ENV") == "production" {
//	        set("port", 80)
//	        set("nacos_addr", "prod-nacos.example.com:8848")
//	    } else {
//	        set("nacos_addr", "127.0.0.1:8848")
//	    }
//	}
//
//	func (c *AppConfig) Validate() error {
//	    return validator.Validate(c)
//	}
//
//	func main() {
//	    loader := config.NewLoader()
//	    var cfg AppConfig
//	    if err := loader.Load("config.yaml", &cfg); err != nil {
//	        panic(err)
//	    }
//	    fmt.Printf("Config loaded: %+v\n", cfg)
//	}
//
// # Error Handling
//
// All errors are wrapped with context information using fmt.Errorf with %w.
// This allows for error inspection using errors.Is and errors.As:
//
//	if err := loader.Load("config.yaml", &cfg); err != nil {
//	    if os.IsNotExist(err) {
//	        // handle missing file
//	    }
//	    return fmt.Errorf("config load failed: %w", err)
//	}
//
// # Thread Safety
//
// The Loader type is safe for concurrent use after initialization.
//
// # Related Packages
//
// This package builds on:
//   - github.com/spf13/viper - Configuration management
//   - github.com/creasty/defaults - struct tag defaults
//   - github.com/chinayin/gox/validator - Unified validation
//
// For validation rules and custom validators, see github.com/chinayin/gox/validator package.
package config
