package config

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestSyncControlValidation(t *testing.T) {
	value := Config_Controller_Profile_Control_Assignment_SyncControl{
		Type:       "sync_control",
		Identifier: "Throttle",
		InputValue: Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue{
			Min: 0.0,
			Max: 1.0,
		},
		ActionIncrease: Config_Controller_Profile_Control_Assignment_Action_Keys{
			Keys: "a",
		},
		ActionDecrease: Config_Controller_Profile_Control_Assignment_Action_Keys{
			Keys: "d",
		},
	}

	v := validator.New()
	if err := v.Struct(value); err != nil {
		t.Fatalf("sync control validation failed %e", err)
	}
}
