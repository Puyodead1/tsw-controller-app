package main

import (
	"context"
	"fmt"
	"tsw_controller_app/action_sequencer"
	"tsw_controller_app/controller_mgr"
	"tsw_controller_app/profile_runner"
	"tsw_controller_app/sdl_mgr"
)

type App struct {
	ctx                context.Context
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

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
