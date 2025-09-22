package sdl_mgr

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

/* the SDL control kind like Button, Hat, Axis */
type SDLMgr_Control_Kind = string

const (
	SDLMgr_Control_Kind_Button SDLMgr_Control_Kind = "button"
	SDLMgr_Control_Kind_Hat    SDLMgr_Control_Kind = "hat"
	SDLMgr_Control_Kind_Axis   SDLMgr_Control_Kind = "axis"
)

type SDLMgr_Joystick struct {
	GUID      sdl.JoystickGUID
	Name      string
	VendorID  int
	ProductID int
	Index     int

	IsOpen           bool
	InternalJoystick *sdl.Joystick
}

type SDLMgr struct {
	Initialized bool
}

func New() *SDLMgr {
	return &SDLMgr{}
}

/*
Initializes the SDL library for the app
sdl.Init is guarded to only be ran once per app
*/
func (mgr *SDLMgr) PanicInit() bool {
	if !mgr.Initialized {
		/* try to initialize if not already initialized */
		if err := sdl.Init(sdl.INIT_GAMECONTROLLER | sdl.INIT_JOYSTICK); err != nil {
			panic(err)
		}
	}

	mgr.Initialized = true
	return true
}

/* Just a passthrough for the sdl quit method */
func (mgr *SDLMgr) Quit() {
	sdl.Quit()
}

func (mgr *SDLMgr) GetJoystickByIndex(index int) (*SDLMgr_Joystick, error) {
	if index >= sdl.NumJoysticks() {
		return nil, fmt.Errorf("index is out of range for number of registered SDL joysticks")
	}

	name := sdl.JoystickNameForIndex(index)
	guid := sdl.JoystickGetDeviceGUID(index)
	usb_vendor := sdl.JoystickGetDeviceVendor(index)
	usb_product := sdl.JoystickGetDeviceProduct(index)

	return &SDLMgr_Joystick{
		GUID:      guid,
		Name:      name,
		VendorID:  usb_vendor,
		ProductID: usb_product,
		Index:     index,
		IsOpen:    false,
	}, nil
}

/*
Starts polling for events within a go-routine every 60ms
Can be cancelled using the context
Returns a channel to listen to events
*/
func (mgr *SDLMgr) StartPolling(ctx context.Context) (chan sdl.Event, context.CancelFunc) {
	ctx_with_cancel, cancel := context.WithCancel(ctx)
	event_channel := make(chan sdl.Event)
	go func() {
		/* check for done */
		for {
			if event := sdl.PollEvent(); event != nil {
				/* pass event to channel */
				event_channel <- event
			}

			timer := time.NewTimer(60 * time.Millisecond)
			select {
			case <-ctx_with_cancel.Done():
				/* polling stopped during wait */
				timer.Stop()
				return
			case <-timer.C:
				/* continue to next iteration */
			}
		}
	}()
	return event_channel, cancel
}

func (joystick *SDLMgr_Joystick) ToString() string {
	return fmt.Sprintf("%d:%d", joystick.VendorID, joystick.ProductID)
}

func (joystick *SDLMgr_Joystick) Open() error {
	if joystick.IsOpen {
		return nil
	}

	joystick.InternalJoystick = sdl.JoystickOpen(joystick.Index)
	if joystick.InternalJoystick == nil {
		return fmt.Errorf("could not open joystick for use")
	}
	joystick.IsOpen = true
	return nil
}

func (joystick *SDLMgr_Joystick) Close() error {
	if !joystick.IsOpen {
		return fmt.Errorf("joystick is not open")
	}

	if joystick.InternalJoystick == nil {
		return fmt.Errorf("internal joystick not assigned")
	}

	joystick.InternalJoystick.Close()
	joystick.IsOpen = false
	joystick.InternalJoystick = nil
	return nil
}
