package config

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"tsw_controller_app/math_utils"

	"github.com/go-playground/validator/v10"
)

type PreferredControlMode = string

const (
	PreferredControlMode_DirectControl PreferredControlMode = "direct_control"
	PreferredControlMode_SyncControl   PreferredControlMode = "Sync_control"
)

type FreeRangeZone struct {
	Start float64
	End   float64
}

type Config_Controller_Profile_Control_Assignment_Action_Keys struct {
	Keys      string   `json:"keys" validate:"required" example:"ctrl+a"`
	PressTime *float64 `json:"press_time,omitempty"`
	WaitTime  *float64 `json:"wait_time,omitempty"`
}

type Config_Controller_Profile_Control_Assignment_Action_DirectControl struct {
	Controls string  `json:"controls" validate:"required"`
	Value    float64 `json:"value"`
	/* sets this value to be a relative adjustment as opposed to an absolute one */
	Relative *bool `json:"relative,omitempty"`
	/* determine whether to hold the value or not; meaning the value will be sent continuously */
	Hold *bool `json:"hold,omitempty"`
	/* whether to apply raw or normalized values */
	UseNormalized *bool `json:"use_normalized,omitempty"`
}

type Config_Controller_Profile_Control_Assignment_Action struct {
	Keys          *Config_Controller_Profile_Control_Assignment_Action_Keys          `json:"-"`
	DirectControl *Config_Controller_Profile_Control_Assignment_Action_DirectControl `json:"-"`
}

type Config_Controller_Profile_Control_Assignment_Condition struct {
	/* this is the other control name to depend on */
	Control  string  `json:"control" validate:"required"`
	Operator string  `json:"operator" validate:"required,oneof=gte,lte,gt,lt"`
	Value    float64 `json:"value"`
}

type Config_Controller_Profile_Control_Assignment_Momentary struct {
	Type      string  `json:"type" validate:"required,eq=momentary"`
	Threshold float64 `json:"threshold"`
	/* which action to perform once the threshold is exceeded */
	ActionActivate Config_Controller_Profile_Control_Assignment_Action `json:"action_activate" validate:"required"`
	/* which action to perform once the threshold is not exceeded anymore; defaults to releasing the activate action if keys */
	ActionDeactivate *Config_Controller_Profile_Control_Assignment_Action `json:"action_deactivate,omitempty"`
}

type Config_Controller_Profile_Control_Assignment_Linear_Threshold struct {
	Value float64 `json:"value"`
	/* ValueEnd and ValueStep can be used to automatically generate a set of thresholds while keeping the same action (ie: throttle) */
	ValueEnd  *float64 `json:"value_end,omitempty"`
	ValueStep *float64 `json:"value_step,omitempty"`
	/* which action to perform once the linear threshold is exceeded */
	ActionActivate   Config_Controller_Profile_Control_Assignment_Action  `json:"action_activate" validate:"required"`
	ActionDeactivate *Config_Controller_Profile_Control_Assignment_Action `json:"action_deactivate,omitempty"`
}

type Config_Controller_Profile_Control_Assignment_Linear struct {
	Type       string                                                          `json:"type" validate:"required,eq=linear"`
	Neutral    *float64                                                        `json:"neutral,omitempty"`
	Thresholds []Config_Controller_Profile_Control_Assignment_Linear_Threshold `json:"thresholds" validate:"required"`
}

type Config_Controller_Profile_Control_Assignment_Toggle struct {
	Type      string  `json:"type" validate:"required,eq=toggle"`
	Threshold float64 `json:"threshold"`
	/* which action to perform once the threshold is exceeded */
	ActionActivate   Config_Controller_Profile_Control_Assignment_Action `json:"action_activate" validate:"required"`
	ActionDeactivate Config_Controller_Profile_Control_Assignment_Action `json:"action_deactivate" validate:"required"`
}

type Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue struct {
	Min  float64  `json:"min"`
	Max  float64  `json:"max"`
	Step *float64 `json:"step,omitempty"`
	/** steps can be combined with null values to create automatic interpolation */
	Steps  *[]*float64 `json:"steps,omitempty"`
	Invert *bool       `json:"invert,omitempty"`
}

type Config_Controller_Profile_Control_Assignment_DirectControl struct {
	Type string `json:"type" validate:"required,eq=direct_control"`
	/* the HID control component as per the UE4SS API */
	Controls string `json:"controls" validate:"required"`
	/* will hold the control in changing */
	Hold *bool `json:"hold"`
	/* whether to apply raw or normalized values */
	UseNormalized *bool                                                                       `json:"use_normalized,omitempty"`
	InputValue    Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue `json:"input_value" validate:"required"`
}

type Config_Controller_Profile_Control_Assignment_SyncControl struct {
	Type string `json:"type" validate:"required,eq=sync_control"`
	/** this is the VHID Identifier Name - differs from the direct control name */
	Identifier     string                                                                      `json:"identifier" validate:"required"`
	InputValue     Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue `json:"input_value" validate:"required"`
	ActionIncrease Config_Controller_Profile_Control_Assignment_Action_Keys                    `json:"action_increase" validate:"required"`
	ActionDecrease Config_Controller_Profile_Control_Assignment_Action_Keys                    `json:"action_decrease" validate:"required"`
}

type Config_Controller_Profile_Control_Assignment struct {
	Momentary     *Config_Controller_Profile_Control_Assignment_Momentary     `json:"-"`
	Linear        *Config_Controller_Profile_Control_Assignment_Linear        `json:"-"`
	Toggle        *Config_Controller_Profile_Control_Assignment_Toggle        `json:"-"`
	DirectControl *Config_Controller_Profile_Control_Assignment_DirectControl `json:"-"`
	SyncControl   *Config_Controller_Profile_Control_Assignment_SyncControl   `json:"-"`
	Conditions    *[]Config_Controller_Profile_Control_Assignment_Condition   `json:"conditions,omitempty"`
}

type Config_Controller_Profile_Control struct {
	Name        string                                          `json:"name"`
	Assignment  *Config_Controller_Profile_Control_Assignment   `json:"assignment,omitempty"`
	Assignments *[]Config_Controller_Profile_Control_Assignment `json:"assignments,omitempty"`
}

type Config_Controller_Profile struct {
	Name string `json:"name" validate:"required"`
	/* if provided ; this profile will automatically be used for this controller */
	UsbID    *string                             `json:"usb_id,omitempty"`
	Controls []Config_Controller_Profile_Control `json:"controls" validate:"required"`
}

func (c *Config_Controller_Profile_Control_Assignment_Action) UnmarshalJSON(data []byte) error {
	var peek struct {
		Controls *string `controls:"keys,omitempty"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}

	v := validator.New()

	/* if controls is defined; try to unmarshal it as a direct control action */
	if peek.Controls != nil {
		var dc_action Config_Controller_Profile_Control_Assignment_Action_DirectControl
		if err := json.Unmarshal(data, &dc_action); err != nil {
			return err
		}
		if err := v.Struct(dc_action); err != nil {
			return err
		}
		c.DirectControl = &dc_action
		return nil
	}

	/* default to a keys action */
	var keys_action Config_Controller_Profile_Control_Assignment_Action_Keys
	if err := json.Unmarshal(data, &keys_action); err != nil {
		return err
	}
	if err := v.Struct(keys_action); err != nil {
		return err
	}
	c.Keys = &keys_action
	return nil
}

func (c Config_Controller_Profile_Control_Assignment_Action) MarshalJSON() ([]byte, error) {
	if c.DirectControl != nil {
		return json.Marshal(c.DirectControl)
	}
	if c.Keys != nil {
		return json.Marshal(c.Keys)
	}
	return nil, fmt.Errorf("unable to marshal control assignment action; has to be one of direct_control or keys but neither was found")
}

func (c *Config_Controller_Profile_Control_Assignment) UnmarshalJSON(data []byte) error {
	v := validator.New()

	var peek struct {
		Type       string                                                    `type:"type"`
		Conditions *[]Config_Controller_Profile_Control_Assignment_Condition `json:"conditions,omitempty"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	if err := v.Struct(peek); err != nil {
		return err
	}

	c.Conditions = peek.Conditions
	switch peek.Type {
	case "momentary":
		var momentary Config_Controller_Profile_Control_Assignment_Momentary
		if err := json.Unmarshal(data, &momentary); err != nil {
			return err
		}
		if err := v.Struct(momentary); err != nil {
			return err
		}
		c.Momentary = &momentary
		return nil
	case "linear":
		var linear Config_Controller_Profile_Control_Assignment_Linear
		if err := json.Unmarshal(data, &linear); err != nil {
			return err
		}
		if err := v.Struct(linear); err != nil {
			return err
		}
		c.Linear = &linear
		return nil
	case "toggle":
		var toggle Config_Controller_Profile_Control_Assignment_Toggle
		if err := json.Unmarshal(data, &toggle); err != nil {
			return err
		}
		if err := v.Struct(toggle); err != nil {
			return err
		}
		c.Toggle = &toggle
		return nil
	case "direct_control":
		var dc Config_Controller_Profile_Control_Assignment_DirectControl
		if err := json.Unmarshal(data, &dc); err != nil {
			return err
		}
		if err := v.Struct(dc); err != nil {
			return err
		}
		c.DirectControl = &dc
		return nil
	case "sync_control":
		var sc Config_Controller_Profile_Control_Assignment_SyncControl
		if err := json.Unmarshal(data, &sc); err != nil {
			return err
		}
		if err := v.Struct(sc); err != nil {
			return err
		}
		c.SyncControl = &sc
		return nil
	}
	return fmt.Errorf("invalid assignment type (%s)", peek.Type)
}

func (c Config_Controller_Profile_Control_Assignment) MarshalJSON() ([]byte, error) {
	if c.Momentary != nil {
		return json.Marshal(c.Momentary)
	}
	if c.Linear != nil {
		return json.Marshal(c.Linear)
	}
	if c.Toggle != nil {
		return json.Marshal(c.Toggle)
	}
	if c.DirectControl != nil {
		return json.Marshal(c.DirectControl)
	}
	if c.SyncControl != nil {
		return json.Marshal(c.SyncControl)
	}
	return nil, fmt.Errorf("unable to marshal control assignment; no valid assignment found")
}

func (c *Config_Controller_Profile_Control_Assignment_Action_DirectControl) ToString() string {
	flags := []string{}
	if c.Hold != nil && *c.Hold {
		flags = append(flags, "hold")
	}
	if c.Relative != nil && *c.Relative {
		flags = append(flags, "relative")
	}
	if c.UseNormalized != nil && *c.UseNormalized {
		flags = append(flags, "normalized")
	}

	return fmt.Sprintf("%s,%f,%s", c.Controls, c.Value, strings.Join(flags, "|"))
}

func (c *Config_Controller_Profile_Control_Assignment_Action) ToString() string {
	if c.Keys != nil {
		return c.Keys.Keys
	}
	if c.DirectControl != nil {
		return c.DirectControl.ToString()
	}
	return ""
}

func (c *Config_Controller_Profile_Control_Assignment_Linear_Threshold) IsExceedingThreshold(value float64) bool {
	if c.Value < 0.0 {
		return value < c.Value
	}
	return value >= c.Value
}

func (c *Config_Controller_Profile_Control_Assignment_Linear) GenerateThresholds() []Config_Controller_Profile_Control_Assignment_Linear_Threshold {
	var thresholds []Config_Controller_Profile_Control_Assignment_Linear_Threshold
	for _, threshold := range c.Thresholds {
		if threshold.ValueEnd == nil || threshold.ValueStep == nil {
			thresholds = append(thresholds, threshold)
		} else {
			current_value := threshold.Value
			for current_value <= *threshold.ValueEnd {
				thresholds = append(thresholds, Config_Controller_Profile_Control_Assignment_Linear_Threshold{
					Value: current_value,
					/* generated thresholds don't need these anymore */
					ValueEnd:         nil,
					ValueStep:        nil,
					ActionActivate:   threshold.ActionActivate,
					ActionDeactivate: threshold.ActionDeactivate,
				})
				current_value = math_utils.RoundToMarginOfError(current_value + *threshold.ValueStep)
			}
		}
	}
	return thresholds
}

/*
Normalizes the input value according to the neutral value
*/
func (c *Config_Controller_Profile_Control_Assignment_Linear) CalculateNeutralizedValue(value float64) float64 {
	if c.Neutral != nil && *c.Neutral > 0 {
		return (value - *c.Neutral) * (1.0 / *c.Neutral)
	}
	return value
}

func (c *Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue) GetFreeRangeZones() []FreeRangeZone {
	var zones []FreeRangeZone
	if c.Steps == nil {
		return zones
	}

	previous_value := c.Min
	is_free_range_zone := false
	for _, step := range *c.Steps {
		if step == nil {
			is_free_range_zone = true
		} else {
			if is_free_range_zone {
				zones = append(zones, FreeRangeZone{
					Start: previous_value,
					End:   *step,
				})
				is_free_range_zone = false
				previous_value = *step
			}
		}
	}

	if is_free_range_zone {
		zones = append(zones, FreeRangeZone{
			Start: previous_value,
			End:   c.Max,
		})
	}

	return zones
}

/*
*
Returns the actual defined steps - excluding free range zones.
Free range zones should be handled separately
*/
func (c *Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue) GetNormalSteps() *[]float64 {
	if c.Steps == nil {
		return nil
	}

	var normal_steps []float64
	for _, step := range *c.Steps {
		if step != nil {
			normal_steps = append(normal_steps, *step)
		}
	}

	return &normal_steps
}

/*
*
The incoming value here can only be [-1, 1]
This calculates the actual value which would be sent to the game
*/
func (c *Config_Controller_Profile_Control_Assignment_DirectOrSyncControl_InputValue) CalculateOutputValue(value float64) float64 {
	input_value := value
	if c.Invert != nil && *c.Invert {
		if value < 0.0 {
			input_value = -1.0 - value
		} else {
			input_value = 1.0 - value
		}
	}

	total_distance := math.Abs(c.Max - c.Min)
	normal := (input_value * total_distance) + c.Min
	normal_steps := c.GetNormalSteps()
	free_zones := c.GetFreeRangeZones()
	fmt.Printf("calc output value from %#v | %#v\n", total_distance, normal)

	if normal_steps == nil && c.Step != nil {
		var auto_steps []float64
		current_value := c.Min
		for {
			auto_steps = append(auto_steps, current_value)
			current_value = math.Min(current_value+*c.Step, c.Max)
			if current_value >= c.Max {
				auto_steps = append(auto_steps, c.Max)
				break
			}
		}
		normal_steps = &auto_steps
	}

	if normal_steps != nil {
		/* check free range first */
		for _, zone := range free_zones {
			if normal >= zone.Start && normal <= zone.End {
				/* is within free range zone; clamp */
				return math_utils.Clamp(normal, c.Min, c.Max)
			}
		}

		/* else find closest step */
		closest_step := (*normal_steps)[0]
		for _, step := range *normal_steps {
			if math.Abs(normal-step) < math.Abs(normal-closest_step) {
				closest_step = step
			}
		}

		return closest_step
	}

	return math_utils.Clamp(normal, c.Min, c.Max)
}

func (c *Config_Controller_Profile) FindControlByName(name string) *Config_Controller_Profile_Control {
	for _, control := range c.Controls {
		if control.Name == name {
			return &control
		}
	}
	return nil
}

func ControllerProfileFromJSON(json_str string) (*Config_Controller_Profile, error) {
	var c Config_Controller_Profile
	if err := json.Unmarshal([]byte(json_str), &c); err != nil {
		return nil, err
	}

	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, err
	}

	return &c, nil
}
