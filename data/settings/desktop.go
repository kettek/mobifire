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

	if err := os.MkdirAll(appdir, 0755); err != nil {
		return err
	}

	cfgbytes, err := os.ReadFile(cfgfile)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not read file: %s", err.Error())
		} else {
			if err := os.WriteFile(cfgfile, nil, 0755); err != nil {
				return err
			}
		}
	} else {
		if err := yaml.Unmarshal(cfgbytes, &settings); err != nil {
			return err
		}
	}
	return nil
}

func save() error {
	b, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cfgfile, b, 0755); err != nil {
		return err
	}
	return nil
}
