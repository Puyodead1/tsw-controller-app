package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"tsw_controller_app/logger"
)

const DEFAULT_TSWAPI_SUBSCRIPTION_ID_START = 83211
const DEFAULT_PREFERRED_CONTROL_MODE = PreferredControlMode_DirectControl

type Config_ProgramConfig struct {
	LastInstalledModVersion   string               `json:"last_instalaled_mod_version,omitempty"`
	TSWAPIKeyLocation         string               `json:"tsw_api_key_location,omitempty"`
	TSWAPISubscriptionIDStart int                  `json:"tsw_api_subscription_id_start,omitempty"`
	PreferredControlMode      PreferredControlMode `json:"preferred_control_mode,omitempty"`
}

func LoadProgramConfigFromFile(filepath string) *Config_ProgramConfig {
	file_bytes, err := os.ReadFile(filepath)
	if err != nil {
		logger.Logger.Error("[Config_ProgramConfig] could not read config file", "filepath", filepath)
		return &Config_ProgramConfig{
			LastInstalledModVersion:   "",
			TSWAPIKeyLocation:         "",
			TSWAPISubscriptionIDStart: DEFAULT_TSWAPI_SUBSCRIPTION_ID_START,
			PreferredControlMode:      DEFAULT_PREFERRED_CONTROL_MODE,
		}
	}

	var pc Config_ProgramConfig = Config_ProgramConfig{
		LastInstalledModVersion:   "",
		TSWAPIKeyLocation:         "",
		TSWAPISubscriptionIDStart: DEFAULT_TSWAPI_SUBSCRIPTION_ID_START,
		PreferredControlMode:      DEFAULT_PREFERRED_CONTROL_MODE,
	}
	if err := json.Unmarshal(file_bytes, &pc); err != nil {
		logger.Logger.Error("[Config_ProgramConfig] failed to parse json", "filepath", filepath)
		return &Config_ProgramConfig{
			LastInstalledModVersion:   "",
			TSWAPIKeyLocation:         "",
			TSWAPISubscriptionIDStart: DEFAULT_TSWAPI_SUBSCRIPTION_ID_START,
			PreferredControlMode:      DEFAULT_PREFERRED_CONTROL_MODE,
		}
	}

	return &pc
}

func (c *Config_ProgramConfig) AutoDetectTSWAPIKeyLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	tsw6_path := filepath.Join(home, "Documents/My Games/TrainSimWorld6/Saved/Config/CommAPIKey.txt")
	tsw5_path := filepath.Join(home, "Documents/My Games/TrainSimWorld6/Saved/Config/CommAPIKey.txt")
	if _, err := os.Stat(tsw6_path); err == nil {
		return tsw6_path
	}
	if _, err := os.Stat(tsw5_path); err == nil {
		return tsw5_path
	}
	return ""
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
