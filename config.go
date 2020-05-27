package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml"
)

var (
	DefaultConfigPath string
)

const (
	DefaultListenAddress = "127.0.0.1:8080"
)

type Config struct {
	// listen address for tcp
	ListenAddress string `toml:"listen-address"`
}

func newConfig() Config {
	var cfg Config
	// apply defaults
	cfg.ListenAddress = DefaultListenAddress
	return cfg
}

func LoadConfig(r io.Reader) (Config, error) {
	cfg := newConfig()

	td := toml.NewDecoder(r)
	if err := td.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func LoadConfigFromFile(name string) (Config, error) {
	rc, err := os.Open(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newConfig(), nil
		} else {
			return Config{}, err
		}
	}
	defer rc.Close()

	return LoadConfig(rc)
}

func getUserHome() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	} else {
		return os.Getenv("HOME")
	}
}

func init() {
	homeDir := os.Getenv("XDG_CONFIG_HOME")
	if homeDir == "" {
		homeDir = filepath.Join(getUserHome(), ".config", progName)
	}
	DefaultConfigPath = filepath.Join(homeDir, "config.toml")
}

// func init() {
// 	homeDir := os.Getenv("XDG_CACHE_HOME")
// 	if homeDir == "" {
// 		homeDir = filepath.Join(getUserHome(), ".cache", progName)
// 	}
// 	DefaultCacheDir = homeDir
// }
