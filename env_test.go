package env

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	type Config struct {
		Port int    `env:"PORT"`
		Dsn  string `env:"DSN"`
	}

	var config Config

	os.Setenv("PORT", "8080")
	os.Setenv("DSN", "user:pass@tcp(localhost:3306)/db")

	Parse(&config)

	if config.Port != 8080 {
		t.Errorf("Expected config.Port to be 8080, got %d", config.Port)
	}

	if config.Dsn != "user:pass@tcp(localhost:3306)/db" {
		t.Errorf("Expected config.Dsn to be user:pass@tcp(localhost:3306)/db, got %s", config.Dsn)
	}

	os.Setenv("PORT", "")

	Parse(&config)

	if config.Port != 0 {
		t.Errorf("Expected config.Port to be 8080, got %d", config.Port)
	}

	if config.Dsn != "user:pass@tcp(localhost:3306)/db" {
		t.Errorf("Expected config.Dsn to be user:pass@tcp(localhost:3306)/db, got %s", config.Dsn)
	}
}

func TestParseNested(t *testing.T) {
	type Config struct {
		Port int `env:"PORT"`
		Sub  struct {
			Dsn string `env:"DSN"`
		}
	}

	var config Config

	os.Setenv("PORT", "8080")
	os.Setenv("DSN", "user:pass@tcp(localhost:3306)/db")

	Parse(&config)

	if config.Port != 8080 {
		t.Errorf("Expected config.Port to be 8080, got %d", config.Port)
	}

	if config.Sub.Dsn != "user:pass@tcp(localhost:3306)/db" {
		t.Errorf("Expected config.Sub.Dsn to be user:pass@tcp(localhost:3306)/db, got %s", config.Sub.Dsn)
	}
}
