package config

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestConfigProfile_SyncControlValidation(t *testing.T) {
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
	assert.NoError(t, v.Struct(value))
}

func TestConfigProfile_LinearThreshold_IsExceedingThreshold(t *testing.T) {
	threshold := Config_Controller_Profile_Control_Assignment_Linear_Threshold{
		Value: 0.5,
	}

	// 0.51 should exceed 0.5
	assert.True(t, threshold.IsExceedingThreshold(0.51))

	// 0.49 should not exceed 0.49
	assert.False(t, threshold.IsExceedingThreshold(0.49))
}

func TestConfigProfile_Linear_GenerateThresholds(t *testing.T) {
	var auto_step_threshold_value_end float64 = 0.4
	var auto_step_threshold_value_step float64 = 0.03

	linear := Config_Controller_Profile_Control_Assignment_Linear{
		Type: "linear",
		Thresholds: []Config_Controller_Profile_Control_Assignment_Linear_Threshold{
			// simple thresholds
			{Value: 0.0},
			{Value: 0.1},
			{Value: 0.2},
			// auto step threshold - value end is exclusive
			{Value: 0.3, ValueEnd: &auto_step_threshold_value_end, ValueStep: &auto_step_threshold_value_step},
			{Value: 0.6},
		},
	}

	thresholds := linear.GenerateThresholds()

	// should have generated 9 thresholds
	assert.Equal(t, thresholds[0].Value, 0.0)
	assert.Equal(t, thresholds[1].Value, 0.1)
	assert.Equal(t, thresholds[2].Value, 0.2)
	assert.Equal(t, thresholds[3].Value, 0.3)
	assert.Equal(t, thresholds[4].Value, 0.33)
	assert.Equal(t, thresholds[5].Value, 0.36)
	assert.Equal(t, thresholds[6].Value, 0.39)
	assert.Equal(t, thresholds[7].Value, 0.6)
}

func TestConfigProfile_Linear_CalculateNeutralizedValue(t *testing.T) {
	var neutral_value float64 = 0.5
	linear := Config_Controller_Profile_Control_Assignment_Linear{
		Type:       "linear",
		Neutral:    &neutral_value,
		Thresholds: []Config_Controller_Profile_Control_Assignment_Linear_Threshold{},
	}

	assert.Equal(t, 0.0, linear.CalculateNeutralizedValue(0.5))
	assert.Equal(t, -1.0, linear.CalculateNeutralizedValue(0))
	assert.Equal(t, 1.0, linear.CalculateNeutralizedValue(1))
}

func TestConfigProfile_DirectOrSyncControl_InputValue_GetFreeRangeZones(t *testing.T) {
	steps_0_0 := 0.0
	steps_0_5 := 0.5
	steps_1_0 := 1.0
	input_value := Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue{
		Min: 0.0,
		Max: 1.0,
		Steps: &[]*float64{
			&steps_0_0, &steps_0_5, nil, &steps_1_0,
		},
	}
	free_range_zones := input_value.GetFreeRangeZones()
	assert.Len(t, free_range_zones, 1)
	assert.Equal(t, FreeRangeZone{Start: 0.5, End: 1.0}, free_range_zones[0])
}

func TestConfigProfile_DirectOrSyncControl_InputValue_GetNormalSteps(t *testing.T) {
	steps_0_0 := 0.0
	steps_0_5 := 0.5
	steps_1_0 := 1.0
	input_value := Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue{
		Min: 0.0,
		Max: 1.0,
		Steps: &[]*float64{
			&steps_0_0, &steps_0_5, nil, &steps_1_0,
		},
	}
	normal_steps := input_value.GetNormalSteps()
	assert.Len(t, *normal_steps, 3)
	assert.ElementsMatch(t, *normal_steps, []float64{0.0, 0.5, 1.0})
}

func TestConfigProfile_DirectOrSyncControl_InputValue_CalculateOutputValue(t *testing.T) {
	steps_0_0 := 0.0
	steps_0_5 := 0.5
	steps_1_0 := 1.0
	input_value := Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue{
		Min: 0.0,
		Max: 1.0,
		Steps: &[]*float64{
			&steps_0_0, &steps_0_5, nil, &steps_1_0,
		},
	}
	assert.Equal(t, 0.0, input_value.CalculateOutputValue(0.0))
	assert.Equal(t, 0.0, input_value.CalculateOutputValue(0.1))
	assert.Equal(t, 0.0, input_value.CalculateOutputValue(0.2))
	assert.Equal(t, 0.5, input_value.CalculateOutputValue(0.3))
	assert.Equal(t, 0.5, input_value.CalculateOutputValue(0.4))
	assert.Equal(t, 0.5, input_value.CalculateOutputValue(0.5))
	assert.Equal(t, 0.6, input_value.CalculateOutputValue(0.6))
	assert.Equal(t, 0.7, input_value.CalculateOutputValue(0.7))
	assert.Equal(t, 0.8, input_value.CalculateOutputValue(0.8))
	assert.Equal(t, 0.9, input_value.CalculateOutputValue(0.9))
	assert.Equal(t, 1.0, input_value.CalculateOutputValue(1.0))
}
