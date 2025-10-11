package main

import (
	"context"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	go_runtime "runtime"
	"sort"
	"strings"
	"tsw_controller_app/action_sequencer"
	"tsw_controller_app/config"
	"tsw_controller_app/config_loader"
	"tsw_controller_app/controller_mgr"
	"tsw_controller_app/logger"
	"tsw_controller_app/profile_runner"
	"tsw_controller_app/sdl_mgr"
	"tsw_controller_app/string_utils"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed mod_assets/*
var mod_assets embed.FS

type AppEventType = string

const (
	AppEventType_JoyDevicesUpdated AppEventType = "joydevices_updated"
	AppEventType_ProfilesUpdated   AppEventType = "profiles_updated"
	AppEventType_RawEvent          AppEventType = "rawevent"
	AppEventType_SyncControlState  AppEventType = "synccontrolstate"
	AppEventType_Log               AppEventType = "log"
)

type ModAssets_Manifest struct {
	Manifest []string `json:"manifest"`
}

type AppRawSubscriber struct {
	Channel   chan controller_mgr.ControllerManager_RawEvent
	Cancel    func()
	LastEvent *controller_mgr.ControllerManager_RawEvent
}

type AppConfig struct {
	GlobalConfigDir string
	LocalConfigDir  string
}

type App struct {
	ctx                context.Context
	config             AppConfig
	program_config     *config.Config_ProgramConfig
	config_loader      *config_loader.ConfigLoader
	sdl_manager        *sdl_mgr.SDLMgr
	controller_manager *controller_mgr.ControllerManager
	action_sequencer   *action_sequencer.ActionSequencer
	socket_connection  *profile_runner.SocketConnection
	direct_controller  *profile_runner.DirectController
	sync_controller    *profile_runner.SyncController
	profile_runner     *profile_runner.ProfileRunner

	raw_subscriber *AppRawSubscriber
}

func NewApp(
	appconfig AppConfig,
) *App {
	sdl_manager := sdl_mgr.New()
	sdl_manager.PanicInit()

	controller_manager := controller_mgr.New(sdl_manager)
	action_sequencer := action_sequencer.New()
	socket_connection := profile_runner.NewSocketConnection()
	direct_controller := profile_runner.NewDirectController(socket_connection)
	sync_controller := profile_runner.NewSyncController(socket_connection)

	return &App{
		config:             appconfig,
		program_config:     config.LoadProgramConfigFromFile(path.Join(appconfig.GlobalConfigDir, "program.json")),
		config_loader:      config_loader.New(),
		sdl_manager:        sdl_manager,
		controller_manager: controller_manager,
		action_sequencer:   action_sequencer,
		socket_connection:  socket_connection,
		direct_controller:  direct_controller,
		sync_controller:    sync_controller,
		profile_runner: profile_runner.New(
			action_sequencer,
			controller_manager,
			direct_controller,
			sync_controller,
		),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.LoadConfiguration()

	go func() {
		channel, unsubscribe := logger.Logger.Listen()
		defer unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-channel:
				runtime.EventsEmit(ctx, AppEventType_Log, msg)
			}
		}
	}()

	go func() {
		a.socket_connection.Start()
	}()

	go func() {
		cancel := a.controller_manager.Attach(a.ctx)
		defer cancel()
		<-ctx.Done()
	}()

	go func() {
		cancel := a.profile_runner.Run(ctx)
		defer cancel()
		<-ctx.Done()
	}()

	go func() {
		cancel := a.action_sequencer.Run(ctx)
		defer cancel()
		<-ctx.Done()
	}()

	go func() {
		cancel := a.direct_controller.Run(ctx)
		defer cancel()
		<-ctx.Done()
	}()

	go func() {
		cancel := a.sync_controller.Run(ctx)
		defer cancel()

		<-ctx.Done()
	}()

	go func() {
		channel, unsubscribe := a.sync_controller.Subscribe()
		defer unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case <-channel:
				runtime.EventsEmit(ctx, AppEventType_SyncControlState)
			}
		}
	}()

	go func() {
		channel, cancel := a.controller_manager.SubscribeJoyDevicesUpdated()
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case <-channel:
				runtime.EventsEmit(a.ctx, AppEventType_JoyDevicesUpdated)
			}
		}
	}()
}

func (a *App) shutdown(ctx context.Context) {
}

func (a *App) GetVersion() string {
	return VERSION
}

func (a *App) GetLastInstalledModVersion() string {
	return a.program_config.LastInstalledModVersion
}

func (a *App) SetLastInstalledModVersion(version string) {
	a.program_config.LastInstalledModVersion = version
	a.program_config.Save(path.Join(a.config.GlobalConfigDir, "program.json"))
}

func (a *App) LoadConfiguration() {
	/* load config from relative config directory */
	dirs_to_load := []string{
		a.config.GlobalConfigDir,
		a.config.LocalConfigDir,
	}

	for _, dir := range dirs_to_load {
		sdl_mappings, calibrations, profiles, errors := a.config_loader.FromDirectory(dir)

		for _, err := range errors {
			logger.Logger.Error("[App] encountered error while reading configuration files", "error", err)
		}

		for _, sdl_mapping := range sdl_mappings {
			var calibration *config.Config_Controller_Calibration
			for _, c := range calibrations {
				if c.UsbID == sdl_mapping.UsbID {
					calibration = &c
					break
				}
			}
			if calibration != nil {
				logger.Logger.Info("[App] registering SDL map and calibration for controller", "name", sdl_mapping.Name, "usb_id", sdl_mapping.UsbID)
				a.controller_manager.RegisterConfig(sdl_mapping, *calibration)
			}
		}
		for _, profile := range profiles {
			logger.Logger.Info("[App] registering profile", "profile", profile.Name)
			a.profile_runner.RegisterProfile(profile)
		}
	}

	runtime.EventsEmit(a.ctx, AppEventType_ProfilesUpdated)
}

func (a *App) GetControllers() []Interop_GenericController {
	var controllers []Interop_GenericController
	a.controller_manager.ConfiguredControllers.ForEach(func(c controller_mgr.ControllerManager_ConfiguredController, _ controller_mgr.JoystickGUIDString) bool {
		controllers = append(controllers, Interop_GenericController{
			GUID:         c.Joystick.GUID,
			UsbID:        c.Joystick.ToString(),
			Name:         c.Joystick.Name,
			IsConfigured: true,
		})
		return true
	})
	a.controller_manager.UnconfiguredControllers.ForEach(func(c controller_mgr.ControllerManager_UnconfiguredController, key controller_mgr.JoystickGUIDString) bool {
		controllers = append(controllers, Interop_GenericController{
			GUID:         c.Joystick.GUID,
			UsbID:        c.Joystick.ToString(),
			Name:         c.Joystick.Name,
			IsConfigured: false,
		})
		return true
	})
	sort.Slice(controllers, func(i, j int) bool {
		return controllers[i].GUID < controllers[j].GUID
	})
	return controllers
}

func (a *App) GetProfiles() []Interop_Profile {
	var profiles []Interop_Profile
	a.profile_runner.Profiles.ForEach(func(profile config.Config_Controller_Profile, key string) bool {
		profiles = append(profiles, Interop_Profile{
			Name: profile.Name,
		})
		return true
	})
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})
	return profiles
}

func (a *App) GetSelectedProfiles() map[controller_mgr.JoystickGUIDString]string {
	selected_profiles := map[controller_mgr.JoystickGUIDString]string{}
	a.profile_runner.Settings.GetSelectedProfiles().ForEach(func(value *config.Config_Controller_Profile, key controller_mgr.JoystickGUIDString) bool {
		selected_profiles[key] = value.Name
		return true
	})
	return selected_profiles
}

func (a *App) GetControllerConfiguration(guid controller_mgr.JoystickGUIDString) *Interop_ControllerConfiguration {
	if controller, has_controller := a.controller_manager.ConfiguredControllers.Get(guid); has_controller {
		/* when configured the SDL map and calibration always exist */
		sdl_mapping, _ := controller.Manager.Config.SDLMappingsByUsbID.Get(controller.Joystick.ToString())
		interop_calibration := Interop_ControllerCalibration{
			Name:     sdl_mapping.Name,
			UsbId:    sdl_mapping.UsbID,
			Controls: []Interop_ControllerCalibration_Control{},
		}
		controller.Controls.ForEach(func(control controller_mgr.ControllerManager_Controller_Control, key string) bool {
			calibration := Interop_ControllerCalibration_Control{
				Kind:        control.SDLMapping.Kind,
				Index:       control.SDLMapping.Index,
				Name:        control.Name,
				Min:         control.Calibration.Min,
				Max:         control.Calibration.Max,
				Idle:        0,
				Deadzone:    0,
				Invert:      false,
				EasingCurve: []float64{0.0, 0.0, 1.0, 1.0},
			}
			if control.Calibration.Idle != nil {
				calibration.Idle = *control.Calibration.Idle
			}
			if control.Calibration.Deadzone != nil {
				calibration.Deadzone = *control.Calibration.Deadzone
			}
			if control.Calibration.Invert != nil {
				calibration.Invert = *control.Calibration.Invert
			}
			if control.Calibration.EasingCurve != nil {
				calibration.EasingCurve = *control.Calibration.EasingCurve
			}
			interop_calibration.Controls = append(interop_calibration.Controls, calibration)
			return true
		})
		return &Interop_ControllerConfiguration{
			SDLMapping:  sdl_mapping,
			Calibration: interop_calibration,
		}
	}
	return nil
}

func (a *App) GetSyncControlState() []Interop_SyncController_ControlState {
	control_states := []Interop_SyncController_ControlState{}
	a.sync_controller.ControlState.ForEach(func(value profile_runner.SyncController_ControlState, key string) bool {
		control_states = append(control_states, Interop_SyncController_ControlState{
			Identifier:             value.Identifier,
			PropertyName:           value.PropertyName,
			CurrentValue:           value.CurrentValue,
			CurrentNormalizedValue: value.CurrentNormalizedValue,
			TargetValue:            value.TargetValue,
			Moving:                 value.Moving,
		})
		return true
	})
	return control_states
}

// https://github.com/LiamMartens/tsw-controller-app/releases/download/v0.2.6/beta.package.zip
func (a *App) GetLatestReleaseVersion() string {
	resp, err := http.Get("https://raw.githubusercontent.com/LiamMartens/tsw-controller-app/refs/heads/main/RELEASE_VERSION")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.Split(string(body), "\n")[0]
}

func (a *App) SelectProfile(guid controller_mgr.JoystickGUIDString, name string) error {
	if err := a.profile_runner.SetProfile(guid, name); err != nil {
		logger.Logger.Error("selected profile", "profile", name)
		return err
	}
	return nil
}

func (a *App) ClearProfile(guid controller_mgr.JoystickGUIDString) {
	a.profile_runner.ClearProfile(guid)
}

func (a *App) LastRawEvent() *Interop_RawEvent {
	if a.raw_subscriber != nil {
		switch e := a.raw_subscriber.LastEvent.Event.(type) {
		case *sdl.JoyAxisEvent:
			return &Interop_RawEvent{
				GUID:      a.raw_subscriber.LastEvent.Joystick.GUID,
				UsbID:     a.raw_subscriber.LastEvent.Joystick.ToString(),
				Kind:      sdl_mgr.SDLMgr_Control_Kind_Axis,
				Index:     int(e.Axis),
				Value:     float64(e.Value),
				Timestamp: int(e.Timestamp),
			}
		case *sdl.JoyButtonEvent:
			return &Interop_RawEvent{
				GUID:      a.raw_subscriber.LastEvent.Joystick.GUID,
				UsbID:     a.raw_subscriber.LastEvent.Joystick.ToString(),
				Kind:      sdl_mgr.SDLMgr_Control_Kind_Button,
				Index:     int(e.Button),
				Value:     float64(e.State),
				Timestamp: int(e.Timestamp),
			}
		case *sdl.JoyHatEvent:
			return &Interop_RawEvent{
				GUID:      a.raw_subscriber.LastEvent.Joystick.GUID,
				UsbID:     a.raw_subscriber.LastEvent.Joystick.ToString(),
				Kind:      sdl_mgr.SDLMgr_Control_Kind_Hat,
				Index:     int(e.Hat),
				Value:     float64(e.Value),
				Timestamp: int(e.Timestamp),
			}
		}
	}
	return nil
}

func (a *App) UnsubscribeRaw() {
	if a.raw_subscriber != nil {
		a.raw_subscriber.Cancel()
		a.raw_subscriber = nil
	}
}

func (a *App) SubscribeRaw(guid string) error {
	if a.raw_subscriber != nil {
		logger.Logger.Error("already listening")
		return fmt.Errorf("already listening")
	}

	var joystick *sdl_mgr.SDLMgr_Joystick
	if j, has_unconfigured_joystick := a.controller_manager.UnconfiguredControllers.Get(guid); has_unconfigured_joystick {
		joystick = j.Joystick
	} else if j, has_configured_joystick := a.controller_manager.ConfiguredControllers.Get(guid); has_configured_joystick {
		joystick = j.Joystick
	}

	if joystick == nil {
		logger.Logger.Error("joystick not found")
		return fmt.Errorf("joystick not found")
	}

	channel, cancel := a.controller_manager.SubscribeRaw()
	raw_subscriber := AppRawSubscriber{
		Channel: channel,
		Cancel:  cancel,
	}
	go func() {
		for e := range channel {
			if e.Joystick.GUID == joystick.GUID {
				raw_subscriber.LastEvent = &e
				if event := a.LastRawEvent(); event != nil {
					runtime.EventsEmit(a.ctx, AppEventType_RawEvent, event)
				}
			}
		}
	}()
	a.raw_subscriber = &raw_subscriber

	return nil
}

func (a *App) OpenProfileBuilder(name string) {
	if len(name) == 0 {
		runtime.BrowserOpenURL(a.ctx, "https://tsw-controller-app.vercel.app/")
	} else if profile, has_profile := a.profile_runner.Profiles.Get(name); has_profile {
		profile_json, _ := json.Marshal(profile)
		encoded := base64.StdEncoding.EncodeToString(profile_json)
		runtime.BrowserOpenURL(a.ctx, fmt.Sprintf("https://tsw-controller-app.vercel.app/?profile=%s", encoded))
	}
}

func (a *App) SaveCalibration(data Interop_ControllerCalibration) error {
	sdl_mapping := config.Config_Controller_SDLMap{
		Name:  data.Name,
		UsbID: data.UsbId,
		Data:  []config.Config_Controller_SDLMap_Control{},
	}
	calibration := config.Config_Controller_Calibration{
		UsbID: data.UsbId,
		Data:  []config.Config_Controller_CalibrationData{},
	}
	for _, control := range data.Controls {
		if control.Name != "" {
			sdl_mapping.Data = append(sdl_mapping.Data, config.Config_Controller_SDLMap_Control{
				Kind:  control.Kind,
				Index: control.Index,
				Name:  control.Name,
			})
			if control.Kind == sdl_mgr.SDLMgr_Control_Kind_Axis {
				calibration.Data = append(calibration.Data, config.Config_Controller_CalibrationData{
					Id:          control.Name,
					Min:         control.Min,
					Max:         control.Max,
					Idle:        &control.Idle,
					Deadzone:    &control.Deadzone,
					Invert:      &control.Invert,
					EasingCurve: &control.EasingCurve,
				})
			}
		}
	}

	sdl_mapping_filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Select SDL mapping file save location",
		DefaultFilename:  fmt.Sprintf("%s.sdl.json", string_utils.Sluggify(data.Name)),
		DefaultDirectory: "./config/sdl_mappings",
	})
	if err != nil {
		return err
	}

	calibration_filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Select calibration file save location",
		DefaultFilename:  fmt.Sprintf("%s.calibration.json", string_utils.Sluggify(data.Name)),
		DefaultDirectory: "./config/calibration",
	})
	if err != nil {
		return err
	}

	sdl_mapping_file, err := os.OpenFile(sdl_mapping_filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer sdl_mapping_file.Close()

	calibration_file, err := os.OpenFile(calibration_filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer calibration_file.Close()

	encoder_sdl_mapping_file := json.NewEncoder(sdl_mapping_file)
	encoder_sdl_mapping_file.SetIndent("", "  ")
	if err := encoder_sdl_mapping_file.Encode(sdl_mapping); err != nil {
		return err
	}

	encoder_calibration_file := json.NewEncoder(calibration_file)
	encoder_calibration_file.SetIndent("", "  ")
	if err := encoder_calibration_file.Encode(calibration); err != nil {
		return err
	}

	/* register config */
	a.controller_manager.RegisterConfig(sdl_mapping, calibration)

	return nil
}

func (a *App) HasNewerVersion() bool {
	return true
}

func (a *App) UpdateApp() bool {
	return true
}

func (a *App) OpenConfigDirectory() error {
	var cmd *exec.Cmd
	switch go_runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", a.config.GlobalConfigDir)
	case "darwin":
		cmd = exec.Command("open", a.config.GlobalConfigDir)
	default:
		cmd = exec.Command("xdg-open", a.config.GlobalConfigDir)
	}
	fmt.Printf("%#v\n", cmd)
	if err := cmd.Start(); err != nil {
		logger.Logger.Error("[App::OpenConfigDirectory] could not open config directory", "error", err)
		return err
	}
	return nil
}

func (a *App) InstallTrainSimWorldMod() error {
	tsw_exe_path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Train Sim World 5/6 executable (TrainSimWorld.exe)",
	})
	if err != nil {
		return err
	}

	if path.Base(tsw_exe_path) != "TrainSimWorld.exe" {
		return fmt.Errorf("you need to select the TrainSimWorld.exe file to install the mod")
	}

	var manifest ModAssets_Manifest
	manifest_json_bytes, err := mod_assets.ReadFile("mod_assets/manifest.json")
	if err != nil {
		logger.Logger.Error("[App::InstallMod] failed to read manfiest file", "error", err)
		return err
	}

	if err := json.Unmarshal(manifest_json_bytes, &manifest); err != nil {
		return err
	}

	install_path := path.Dir(tsw_exe_path)
	/* go through files to copy */
	for _, file := range manifest.Manifest {
		file_dir := path.Dir(file)
		if err := os.MkdirAll(path.Join(install_path, file_dir), 0755); err != nil {
			logger.Logger.Error("[App::InstallMod] could not create directory", "dir", path.Join(install_path, file_dir))
			return err
		}

		fh, err := mod_assets.Open(path.Join("mod_assets", file))
		if err != nil {
			logger.Logger.Error("[App::InstallMod] could open file", "file", file)
			return fmt.Errorf("could not open file %e", err)
		}
		defer fh.Close()

		out, err := os.Create(path.Join(install_path, file))
		if err != nil {
			logger.Logger.Error("[App::InstallMod] could not create file", "file", path.Join(install_path, file))
			return fmt.Errorf("could not open create %e", err)
		}
		if _, err := io.Copy(out, fh); err != nil {
			logger.Logger.Error("[App::InstallMod] failed to copy file", "file", path.Join(install_path, file))
			return fmt.Errorf("failed to copy file: %w", err)
		}

		defer out.Close()
	}

	/* write version file */
	os.WriteFile(path.Join(install_path, "ue4ss_tsw_controller_mod/Mods/TSWControllerMod/version.txt"), []byte(VERSION), 0755)
	a.program_config.LastInstalledModVersion = VERSION
	a.program_config.Save(path.Join(a.config.GlobalConfigDir, "program.json"))

	return nil
}
