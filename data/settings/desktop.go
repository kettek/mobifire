//go:build !android

package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfgfile string

func load() error {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	appdir := filepath.Join(cdir, "mobifire")
	cfgfile = filepath.Join(appdir, "settings.yaml")

	if err := os.MkdirAll(appdir, 0o755); err != nil {
		return err
	}

	lsettings := make(map[string]RawMessage)

	cfgbytes, err := os.ReadFile(cfgfile)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not read file: %s", err.Error())
		} else {
			if err := os.WriteFile(cfgfile, nil, 0o755); err != nil {
				return err
			}
		}
	} else {
		if err := yaml.Unmarshal(cfgbytes, &lsettings); err != nil {
			return err
		}
	}
	for k, v := range lsettings {
		settings[k] = Setting{
			raw: v,
		}
	}
	return nil
}

func save() error {
	lsettings := make(map[string]any)
	for k, v := range settings {
		lsettings[k] = v.value
	}
	b, err := yaml.Marshal(lsettings)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cfgfile, b, 0o755); err != nil {
		return err
	}
	return nil
}
