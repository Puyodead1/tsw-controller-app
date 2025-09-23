package main

import (
	"context"
	"os"
	"strings"
	"tsw_controller_app/action_sequencer"
	"tsw_controller_app/config"
	"tsw_controller_app/config_loader"
	"tsw_controller_app/controller_mgr"
	"tsw_controller_app/logger"
	"tsw_controller_app/profile_runner"
	"tsw_controller_app/sdl_mgr"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppEventType = string

const (
	AppEventType_Log                     AppEventType = "log"
	AppEventType_JoyDeviceAddedOrRemoved AppEventType = "joydevice_added_or_removed"
)

type App struct {
	ctx                context.Context
	config_loader      *config_loader.ConfigLoader
	sdl_manager        *sdl_mgr.SDLMgr
	controller_manager *controller_mgr.ControllerManager
	action_sequencer   *action_sequencer.ActionSequencer
	direct_controller  *profile_runner.DirectController
	sync_controller    *profile_runner.SyncController
	profile_runner     *profile_runner.ProfileRunner
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
	/* start logger goroutine */
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-logger.Logger.Channel:
				runtime.EventsEmit(a.ctx, AppEventType_Log, msg)
			}
		}
	}()
}

func (a *App) OnFrontendReady() {
	/* load config from relative config directory */
	sdl_mappings, calibrations, profiles, errors := a.config_loader.FromDirectory("./config")
	for _, err := range errors {
		logger.Logger.LogF("[App] encountered error while reading configuration files (%s)\n", err)
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
			logger.Logger.LogF("[App] registering SDL map and calibration for controller (%s)\n", sdl_mapping.Name)
			a.controller_manager.RegisterConfig(sdl_mapping, *calibration)
		}
	}
	for _, profile := range profiles {
		logger.Logger.LogF("[App] registering profile (%s)\n", profile.Name)
		a.profile_runner.RegisterProfile(profile)
	}

	a.controller_manager.Attach(a.ctx)
	go func() {
		channel, _ := a.controller_manager.SubscribeJoyDeviceOrRemoved()
		for range channel {
			runtime.EventsEmit(a.ctx, AppEventType_JoyDeviceAddedOrRemoved)
		}
	}()
}

func (a *App) SaveLogsAsFile() error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "log.txt",
	})
	if err != nil {
		logger.Logger.LogF("[App::SaveLogsAsFile] invalid save location (%e)", err)
		return err
	}

	if err = os.WriteFile(path, []byte(strings.Join(logger.Logger.Logs, "")), 0644); err != nil {
		logger.Logger.LogF("[App::SaveLogsAsFile] failed to save logs (%e)", err)
		return err
	}

	return nil
}

func (a *App) GetControllers() []Interop_GenericController {
	var controllers []Interop_GenericController
	for _, c := range a.controller_manager.ConfiguredControllers {
		controllers = append(controllers, Interop_GenericController{
			Name:         c.Joystick.Name,
			IsConfigured: true,
		})
	}
	for _, c := range a.controller_manager.UnconfiguredControllers {
		controllers = append(controllers, Interop_GenericController{
			Name:         c.Joystick.Name,
			IsConfigured: false,
		})
	}
	return controllers
}
