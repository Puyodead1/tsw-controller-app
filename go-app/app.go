package main

import (
	"context"
	"tsw_controller_app/action_sequencer"
	"tsw_controller_app/config_loader"
	"tsw_controller_app/controller_mgr"
	"tsw_controller_app/profile_runner"
	"tsw_controller_app/sdl_mgr"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppEventType = string

const (
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
	a.controller_manager.Attach(ctx)
	go func() {
		channel, _ := a.controller_manager.SubscribeJoyDeviceOrRemoved()
		for range channel {
			runtime.EventsEmit(ctx, AppEventType_JoyDeviceAddedOrRemoved)
		}
	}()
	a.ctx = ctx
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
