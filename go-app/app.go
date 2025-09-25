package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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

type AppEventType = string

const (
	AppEventType_JoyDevicesUpdated AppEventType = "joydevices_updated"
	AppEvent_RawEvent              AppEventType = "rawevent"
)

type AppRawSubscriber struct {
	Channel   chan controller_mgr.ControllerManager_RawEvent
	Cancel    func()
	LastEvent *controller_mgr.ControllerManager_RawEvent
}

type App struct {
	ctx                context.Context
	config_loader      *config_loader.ConfigLoader
	sdl_manager        *sdl_mgr.SDLMgr
	controller_manager *controller_mgr.ControllerManager
	action_sequencer   *action_sequencer.ActionSequencer
	direct_controller  *profile_runner.DirectController
	sync_controller    *profile_runner.SyncController
	profile_runner     *profile_runner.ProfileRunner

	raw_subscriber *AppRawSubscriber
}

func NewApp() *App {
	sdl_manager := sdl_mgr.New()
	sdl_manager.PanicInit()

	controller_manager := controller_mgr.New(sdl_manager)
	action_sequencer := action_sequencer.New()
	direct_controller := profile_runner.NewDirectController()
	sync_controller := profile_runner.NewSyncController()

	return &App{
		config_loader:      config_loader.New(),
		sdl_manager:        sdl_manager,
		controller_manager: controller_manager,
		action_sequencer:   action_sequencer,
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
}

func (a *App) domReady(ctx context.Context) {
}

func (a *App) OnFrontendReady() {
	/* load config from relative config directory */
	sdl_mappings, calibrations, profiles, errors := a.config_loader.FromDirectory("./config")
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
			logger.Logger.Info("[App] registering SDL map and calibration for controller", "name", sdl_mapping.Name)
			a.controller_manager.RegisterConfig(sdl_mapping, *calibration)
		}
	}
	for _, profile := range profiles {
		logger.Logger.Info("[App] registering profile", "profile", profile.Name)
		a.profile_runner.RegisterProfile(profile)
	}

	a.controller_manager.Attach(a.ctx)
	go func() {
		channel, _ := a.controller_manager.SubscribeJoyDevicesUpdated()
		for range channel {
			runtime.EventsEmit(a.ctx, AppEventType_JoyDevicesUpdated)
		}
	}()
}

func (a *App) GetControllers() []Interop_GenericController {
	var controllers []Interop_GenericController
	for _, c := range a.controller_manager.ConfiguredControllers {
		controllers = append(controllers, Interop_GenericController{
			UsbID:        c.Joystick.ToString(),
			Name:         c.Joystick.Name,
			IsConfigured: true,
		})
	}
	for _, c := range a.controller_manager.UnconfiguredControllers {
		controllers = append(controllers, Interop_GenericController{
			GUID:         c.Joystick.GUID,
			UsbID:        c.Joystick.ToString(),
			Name:         c.Joystick.Name,
			IsConfigured: false,
		})
	}
	return controllers
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
	if j, has_unconfigured_joystick := a.controller_manager.UnconfiguredControllers[guid]; has_unconfigured_joystick {
		joystick = &j.Joystick
	} else if j, has_configured_joystick := a.controller_manager.ConfiguredControllers[guid]; has_configured_joystick {
		joystick = &j.Joystick
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
					runtime.EventsEmit(a.ctx, AppEvent_RawEvent, event)
				}
			}
		}
	}()
	a.raw_subscriber = &raw_subscriber

	return nil
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
					EasingCurve: &[]float64{0.0, 0.0, 1.0, 1.0},
				})
			}
		}
	}

	sdl_mapping_filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Select SDL mapping file save location",
		DefaultFilename: fmt.Sprintf("%s.sdl.json", string_utils.Sluggify(data.Name)),
	})
	if err != nil {
		return err
	}

	calibration_filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Select calibration file save location",
		DefaultFilename: fmt.Sprintf("%s.calibration.json", string_utils.Sluggify(data.Name)),
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
	if err := encoder_sdl_mapping_file.Encode(sdl_mapping); err != nil {
		return err
	}

	encoder_calibration_file := json.NewEncoder(calibration_file)
	if err := encoder_calibration_file.Encode(calibration); err != nil {
		return err
	}

	/* register config */
	a.controller_manager.RegisterConfig(sdl_mapping, calibration)

	return nil
}
