package main

/*
#include <stdint.h>
#include <stdlib.h>
typedef void (*MessageCallback)(const char*);

static inline void CallMessageCallback(MessageCallback cb, const char* msg) {
	cb(msg);
}
*/
import "C"
import (
	"context"
	"unsafe"
)

type DLLState struct {
	Context    context.Context
	Cancel     func()
	Callback   *C.MessageCallback
	Connection *Connection
}

var state DLLState = DLLState{
	Context:    context.TODO(),
	Cancel:     func() {},
	Callback:   nil,
	Connection: nil,
}

//export tsw_controller_mod_start
func tsw_controller_mod_start() {
	ctx, cancel := context.WithCancel(context.Background())
	state.Connection = NewConnection(ctx)
	state.Cancel = cancel

	callback_channel, _ := state.Connection.Subscribe()
	go func() {
		for msg := range callback_channel {
			if state.Callback != nil {
				cmsg := C.CString(msg)
				C.CallMessageCallback(*state.Callback, cmsg)
				C.free(unsafe.Pointer(cmsg))
			}
		}
	}()

	go func() {
		state.Connection.Listen()
	}()
}

//export tsw_controller_mod_stop
func tsw_controller_mod_stop() {
	state.Cancel()
}

//export tsw_controller_mod_set_receive_message_callback
func tsw_controller_mod_set_receive_message_callback(callback C.MessageCallback) {
	state.Callback = &callback
}

//export tsw_controller_mod_send_message
func tsw_controller_mod_send_message(message *C.char) {
	if state.Connection != nil {
		state.Connection.Send(C.GoString(message))
	}
}

//export DllMain
func DllMain(_hinstDLL unsafe.Pointer, _fdwReason C.uint32_t, _lpReserved unsafe.Pointer) C.int {
	return 1
}

func main() {}
