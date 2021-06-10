// Copyright © 2021 The Sanuscoin Team

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"sanus/sanus-sdk/misc/disk"

	dotenv "github.com/joho/godotenv"
	json "github.com/json-iterator/go"
	"github.com/kelseyhightower/envconfig"
)

const (
	dotenvFileName = ".env"
)

type (
	// Config wraps all configuration parts.
	Config struct {
		App *appConfig
		Dst *dstConfig
		Net *netConfig
	}
)

// NewConfig returns config instance.
func NewConfig() *Config {
	return &Config{
		App: newAppConfig(),
		Dst: newDstConfig(),
		Net: newNetConfig(),
	}
}

// Init inits config.
func Init(cfg *Config) error {
	// check if configuration already initialized
	if cfg.App.Name == "" {
		cfg.App.Name = appDefaultName
	}
	return cfg.write()
}

// Load loads config from .env file.
func Load(cfg *Config) error {
	if err := dotenv.Load(filepath.Join(appRootPath, dotenvFileName)); err != nil {
		if os.IsNotExist(err) {
			cfg = NewConfig()
			err = cfg.write()
		}
		if err != nil {
			return err
		}
	}

	return cfg.prepare()
}

// Make makes config globals.
func Make(cfg *Config) error {
	InitPaths(cfg)
	return disk.MakeDirs(
		appRootPath,
		appLogsPath,
		appDataPath,
	)
}

// prepare prepares environment config.
func (c *Config) prepare() error {
	cfg := reflect.ValueOf(c).Elem()
	for idx, n := 0, cfg.NumField(); idx < n; idx++ {
		if fld := cfg.Field(idx); fld.CanInterface() {
			if err := envconfig.Process(
				cfg.Type().Field(idx).Name,
				fld.Interface(),
			); err != nil {
				return fmt.Errorf("prepare config: %w", err)
			}
		}
	}

	return nil
}

// write writes configuration to .env file.
func (c *Config) write() (err error) {
	rows := make([]string, 0)
	elem := reflect.ValueOf(c).Elem()
	for idx, n := 0, elem.NumField(); idx < n; idx++ {
		pref := elem.Type().Field(idx).Name
		part := elem.Field(idx).Elem()
		for idx, n := 0, part.NumField(); idx < n; idx++ {
			val := ""
			fld := part.Field(idx)
			typ := fld.Type()
			if fld.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
				val = fmt.Sprintf(`"%s"`, time.Duration(fld.Int()).String())
			} else if fld.Kind() == reflect.Slice {
				if ss, ok := fld.Interface().([]string); ok {
					val = fmt.Sprintf(`"%s"`, strings.Join(ss, ","))
				} else if val, err = json.MarshalToString(fld.Interface()); err != nil {
					return fmt.Errorf("config option marshal: %w", err)
				}
			} else if fld.Kind() == reflect.Bool || !fld.IsZero() {
				if val, err = json.MarshalToString(fld.Interface()); err != nil {
					return fmt.Errorf("config option marshal: %w", err)
				}
			}
			key := strings.ToUpper(pref + "_" + part.Type().Field(idx).Name)
			rows = append(rows, fmt.Sprintf("%s=%v", key, val))
		}
		rows = append(rows, "") // empty line separate config groups
	}

	file, err := disk.Create(filepath.Join(appRootPath, dotenvFileName))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	_, err = file.WriteString(strings.Join(rows, "\n"))

	return err
}

// IsTestnet function returns true if service worked in testnet mode
func (c *Config) IsTestnet() bool {
	return c.Net.Testnet
}
