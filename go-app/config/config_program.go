package config

import (
	"encoding/json"
	"os"
	"tsw_controller_app/logger"
)

type Config_ProgramConfig struct {
	LastInstalledModVersion string `json:"last_instalaled_mod_version"`
}

func LoadProgramConfigFromFile(filepath string) *Config_ProgramConfig {
	file_bytes, err := os.ReadFile(filepath)
	if err != nil {
		logger.Logger.Error("[Config_ProgramConfig] could not read config file", "filepath", filepath)
		return &Config_ProgramConfig{
			LastInstalledModVersion: "",
		}
	}

	var pc Config_ProgramConfig
	if err := json.Unmarshal(file_bytes, &pc); err != nil {
		logger.Logger.Error("[Config_ProgramConfig] failed to parse json", "filepath", filepath)
		return &Config_ProgramConfig{
			LastInstalledModVersion: "",
		}
	}

	return &pc
}

func (c *Config_ProgramConfig) Save(filepath string) error {
	json_bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath, json_bytes, 0755); err != nil {
		return err
	}
	return nil
}
