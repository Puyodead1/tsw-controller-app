package main

import (
	"tsw_controller_app/config"
	"tsw_controller_app/sdl_mgr"
)

type Interop_GenericController struct {
	GUID         string
	UsbID        string
	Name         string
	IsConfigured bool
}

type Interop_Profile struct {
	Name  string
	UsbID string
}

type Interop_RawEvent struct {
	GUID      string
	UsbID     string
	Kind      sdl_mgr.SDLMgr_Control_Kind
	Index     int
	Value     float64
	Timestamp int
}

type Interop_ControllerCalibration_Control struct {
	Kind        sdl_mgr.SDLMgr_Control_Kind
	Index       int
	Name        string
	Min         float64
	Max         float64
	Idle        float64
	Deadzone    float64
	EasingCurve []float64
	Invert      bool
}

type Interop_ControllerCalibration struct {
	Name     string
	UsbId    string
	Controls []Interop_ControllerCalibration_Control
}

type Interop_ControllerConfiguration struct {
	Calibration Interop_ControllerCalibration
	SDLMapping  config.Config_Controller_SDLMap
}

type Interop_SyncController_ControlState struct {
	Identifier             string
	PropertyName           string
	CurrentValue           float64
	CurrentNormalizedValue float64
	TargetValue            float64
	/** [-1,0,1] -> decreasing, idle, increasing */
	Moving int
}

type Interop_SharedProfile struct {
	Name  string
	UsbID string
	Url   string
}
