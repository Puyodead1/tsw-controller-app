package config_loader

import (
	"fmt"
	"os"
	"path"
	"strings"
	"tsw_controller_app/config"
)

type ConfigLoader struct{}

func New() *ConfigLoader {
	return &ConfigLoader{}
}

func (c *ConfigLoader) FromDirectory(dir string) ([]config.Config_Controller_SDLMap, []config.Config_Controller_Calibration, []config.Config_Controller_Profile, []error) {
	var errors []error

	calibration_files_dir := path.Join(dir, "calibration")
	sdl_mapping_files_dir := path.Join(dir, "sdl_mappings")
	profiles_files_dir := path.Join(dir, "profiles")

	calibration_file_entries, err := os.ReadDir(calibration_files_dir)
	var parsed_calibration_files []config.Config_Controller_Calibration
	if err != nil {
		errors = append(errors, fmt.Errorf("could not read calibration directory (%s)", calibration_files_dir))
	} else {
		for _, entry := range calibration_file_entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				file_bytes, err := os.ReadFile(path.Join(calibration_files_dir, entry.Name()))
				if err != nil {
					errors = append(errors, fmt.Errorf("could not read calibration file (%s)", entry.Name()))
					continue
				}
				calibration, err := config.ControllerCalibrationFromJSON(string(file_bytes))
				if err != nil {
					errors = append(errors, fmt.Errorf("could not parse calibration file (%s)", entry.Name()))
					continue
				}
				parsed_calibration_files = append(parsed_calibration_files, *calibration)
			}
		}
	}

	sdl_mappings_file_entries, err := os.ReadDir(sdl_mapping_files_dir)
	var parsed_sdl_mappings_files []config.Config_Controller_SDLMap
	if err != nil {
		errors = append(errors, fmt.Errorf("could not read SDL mappings directory (%s)", sdl_mapping_files_dir))
	} else {
		for _, entry := range sdl_mappings_file_entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				file_bytes, err := os.ReadFile(path.Join(sdl_mapping_files_dir, entry.Name()))
				if err != nil {
					errors = append(errors, fmt.Errorf("could not read SDL mapping file (%s)", entry.Name()))
					continue
				}
				sdl_mapping, err := config.ControllerSDLMapFromJSON(string(file_bytes))
				if err != nil {
					errors = append(errors, fmt.Errorf("could not parse SDL mapping file (%s)", entry.Name()))
					continue
				}
				parsed_sdl_mappings_files = append(parsed_sdl_mappings_files, *sdl_mapping)
			}
		}
	}

	profiles_file_entries, err := os.ReadDir(profiles_files_dir)
	var parsed_profile_files []config.Config_Controller_Profile
	if err != nil {
		errors = append(errors, fmt.Errorf("could not read profiles directory (%s)", profiles_files_dir))
	} else {
		for _, entry := range profiles_file_entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				file_bytes, err := os.ReadFile(path.Join(profiles_files_dir, entry.Name()))
				if err != nil {
					errors = append(errors, fmt.Errorf("could not read profile file (%s)", entry.Name()))
					continue
				}
				profile, err := config.ControllerProfileFromJSON(string(file_bytes))
				if err != nil {
					errors = append(errors, fmt.Errorf("could not parse profile (%s)", entry.Name()))
					continue
				}
				parsed_profile_files = append(parsed_profile_files, *profile)
			}
		}
	}

	return parsed_sdl_mappings_files, parsed_calibration_files, parsed_profile_files, errors
}
