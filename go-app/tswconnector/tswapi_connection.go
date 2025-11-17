package tswconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"time"
	"tsw_controller_app/pubsub_utils"
	tickerutils "tsw_controller_app/ticker_utils"
)

type CurrentFormationClass = string

type TSWAPIConnectionCabState struct {
	Name     string
	Property string
	Value    float64
}

type TSWAPIConnectionConfig struct {
	BaseURL    string `example:"http://localhost:31270"`
	CommAPIKey string
}

type TSWAPIConnection struct {
	context                  context.Context
	transport                *http.Transport
	client                   *http.Client
	ticker                   *tickerutils.PausableTicker
	cabSubscriptionIdCounter int
	cabSubscriptionIds       map[CurrentFormationClass]int /* http://localhost:31270/get/CurrentFormation/0.ObjectClass -> classname:subscriptionID */
	cabStates                map[CurrentFormationClass]map[string]TSWAPIConnectionCabState
	Config                   TSWAPIConnectionConfig
	Subscribers              *pubsub_utils.PubSubSlice[TSWConnector_Message]
}

// http://localhost:31270/get/CurrentFormation/0.ObjectClass
var _ TSWConnector = (*TSWAPIConnection)(nil)

func (c *TSWAPIConnection) parseApiResponse(r io.ReadCloser) (map[string]any, error) {
	var data map[string]any
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}

	if error_code, has_error_code := data["errorCode"]; has_error_code {
		return nil, fmt.Errorf("%s: %s", error_code.(string), data["errorMessage"].(string))
	}

	if _, has_result := data["Result"]; !has_result {
		return nil, fmt.Errorf("invalid_response: Invalid response")
	}

	result := data["Result"].(string)
	if result == "Error" {
		return nil, fmt.Errorf("%s: %s", result, data["Message"].(string))
	}

	return data, nil
}

func (c *TSWAPIConnection) executeTswApiRequest(req *http.Request) (map[string]any, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parseApiResponse(resp.Body)
}

func (c *TSWAPIConnection) currentFormationClass() (CurrentFormationClass, error) {
	req, _ := http.NewRequest("GET", path.Join(c.Config.BaseURL, "get/CurrentFormation/0.ObjectClass"), nil)
	data, err := c.executeTswApiRequest(req)
	if err != nil {
		return "", err
	}
	object_class := data["Values"].(map[string]string)["ObjectClass"]
	return object_class, nil
}

func (c *TSWAPIConnection) getControlInputValue(control string) (float64, error) {
	req, _ := http.NewRequest("GET", path.Join(c.Config.BaseURL, "get/CurrentDrivableActor", fmt.Sprintf("%s.InputValue", control)), nil)
	data, err := c.executeTswApiRequest(req)
	if err != nil {
		return 0, err
	}
	input_value := data["Values"].(map[string]float64)["InputValue"]
	return input_value, nil
}

func (c *TSWAPIConnection) deleteSubscription(subscription_id int) error {
	req, _ := http.NewRequest("DELETE", path.Join(c.Config.BaseURL, fmt.Sprintf("subscription?Subscription=%d", subscription_id)), nil)
	if _, err := c.executeTswApiRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *TSWAPIConnection) createDrivableActorSubscription(subscription_id int) error {
	req, _ := http.NewRequest("GET", path.Join(c.Config.BaseURL, "list/CurrentDrivableActor"), nil)
	data, err := c.executeTswApiRequest(req)
	if err != nil {
		return err
	}

	nodes := data["Nodes"].([]map[string]any)
	for _, node := range nodes {
		name := node["Name"].(string)
		if name != "ModelChildActorComponent0" &&
			name != "RailVehiclePhysicsComponent0" &&
			name != "GameplayTasksComponent0" &&
			name != "GameplayTagStatusComponent0" &&
			name != "Simulation" &&
			name != "DigitalDisplayService" {
			if _, err := c.getControlInputValue(name); err == nil {
				input_value_req, _ := http.NewRequest("POST", path.Join(c.Config.BaseURL, fmt.Sprintf("subscription/CurrentDrivableActor/%s.InputValue?Subscription=%d", name, subscription_id)), nil)
				property_input_identifier_req, _ := http.NewRequest("POST", path.Join(c.Config.BaseURL, fmt.Sprintf("subscription/CurrentDrivableActor/%s.Property.InputIdentifier?Subscription=%d", name, subscription_id)), nil)
				if _, err := c.executeTswApiRequest(input_value_req); err != nil {
					return err
				}
				if _, err := c.executeTswApiRequest(property_input_identifier_req); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *TSWAPIConnection) updateCabState() error {
	formation_class, err := c.currentFormationClass()
	if err != nil {
		return err
	}

	subscription_id, has_subscription_id := c.cabSubscriptionIds[formation_class]
	if !has_subscription_id {
		c.cabSubscriptionIdCounter++
		subscription_id = c.cabSubscriptionIdCounter
		c.cabSubscriptionIds[formation_class] = c.cabSubscriptionIdCounter
		c.createDrivableActorSubscription(c.cabSubscriptionIdCounter)
	}

	req, _ := http.NewRequest("GET", path.Join(c.Config.BaseURL, fmt.Sprintf("subscription?Subscription=%d", subscription_id)), nil)
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	data, err := c.parseApiResponse(response.Body)
	if err != nil {
		return err
	}

	/**
	"Path": "CurrentDrivableActor/AutomaticBrake_F.InputValue",
	"NodeValid": true,
	"Values": {
		"InputValue": 0.80000001192092896
	}
	*/
	incoming_cab_state_entries := data["Entries"].([]map[string]any)
	/* map entries into key values */
	node_path_rx := regexp.MustCompile(`^CurrentDrivableActor\/(.+?)\..+$`)
	/* map[direct_control_name] = { "property": "direct_control_name", "name": "input_identifier", "value": input_value } */
	incoming_cab_states := map[string]TSWAPIConnectionCabState{}
	for _, entry := range incoming_cab_state_entries {
		path_match := node_path_rx.FindStringSubmatch(entry["Path"].(string))
		if path_match == nil || len(path_match) < 2 {
			continue
		}
		direct_control_name := path_match[1]
		values := entry["Values"].(map[string]any)
		control_state, _ := incoming_cab_states[direct_control_name]
		control_state.Property = direct_control_name
		if value, has_identifier := values["identifier"]; has_identifier {
			control_state.Name = value.(string)
		}
		if value, has_input_value := values["InputValue"]; has_input_value {
			control_state.Value = value.(float64)
		}
		incoming_cab_states[direct_control_name] = control_state
	}

	cab_state, _ := c.cabStates[formation_class]
	for control_name, incoming_state := range incoming_cab_states {
		existing_control_state, has_existing_control_state := cab_state[control_name]
		if !has_existing_control_state || existing_control_state.Value != incoming_state.Value {
			cab_state[control_name] = TSWAPIConnectionCabState{
				Name:     incoming_state.Name,
				Property: incoming_state.Property,
				Value:    incoming_state.Value,
			}
			c.Subscribers.EmitTimeout(time.Second, TSWConnector_Message{
				EventName: "sync_control",
				Properties: map[string]any{
					"name":             incoming_state.Name,
					"property":         incoming_state.Property,
					"value":            incoming_state.Value,
					"normalized_value": incoming_state.Value,
				},
			})
		}
	}
	c.cabStates[formation_class] = cab_state

	return nil
}

func (c *TSWAPIConnection) Subscribe() (chan TSWConnector_Message, func()) {
	return make(chan TSWConnector_Message), func() {}
}

func (c *TSWAPIConnection) Send(m TSWConnector_Message) error {
	if m.EventName == "direct_control" {
		controls := m.Properties["controls"]
		input_value := m.Properties["value"]
		flags := m.Properties["flags"] /* hold,relative */
		// http.NewRequest("POST")

	}
	return nil
}

func (c *TSWAPIConnection) Stop() error {
	c.ticker.Pause()
	return nil
}

func (c *TSWAPIConnection) Start() error {
	c.ticker.Start()
	for {
		select {
		case <-c.ticker.C:
			go c.updateCabState()
		case <-c.ticker.Paused():
			return fmt.Errorf("paused")
		}
	}
}

func NewTSWAPIConnection(ctx context.Context, config TSWAPIConnectionConfig) *TSWAPIConnection {
	child_ctx := context.WithoutCancel(ctx)
	transport := &http.Transport{
		IdleConnTimeout: 120 * time.Second,
	}
	conn := TSWAPIConnection{
		context:                  child_ctx,
		transport:                transport,
		client:                   &http.Client{Transport: transport, Timeout: 2 * time.Second},
		ticker:                   tickerutils.NewPausableTicker(child_ctx, 40*time.Millisecond),
		cabSubscriptionIdCounter: 83222112, /* just a random start number for now */
		cabSubscriptionIds:       map[CurrentFormationClass]int{},
		cabStates:                map[CurrentFormationClass]map[string]TSWAPIConnectionCabState{},
		Config:                   config,
		Subscribers:              pubsub_utils.NewPubSubSlice[TSWConnector_Message](),
	}
	return &conn
}

// type TSWAPI_ListResponse_Node struct {
// 	Name string `json:"Name" validate:"required"`
// }

// type TSWAPI_ListResponse struct {
// 	Result   string                     `json:"Result"`
// 	NodePath string                     `json:"NodePath"`
// 	NodeName string                     `json:"NodeName"`
// 	Nodes    []TSWAPI_ListResponse_Node `json:"Nodes"`
// }

// type TSWAPI_GetResponse struct {
// 	Result string         `json:"Result"`
// 	Values map[string]any `json:"Values"`
// }

// type TSWAPI_PatchResponse struct {
// 	Result string `json:"Result"`
// }

// type TSWAPI_PostSubscriptionResponse struct {
// 	SubscriptionID int  `json:"SubscriptionID"`
// 	CurrentlyValid bool `json:"CurrentlyValid"`
// }

// type TSWAPI_GetSubscriptionResponse_Entry struct {
// 	Path      string         `json:"Path"`
// 	NodeValid bool           `json:"NodeValid"`
// 	Values    map[string]any `json:"Values"`
// }

// type TSWAPI_GetSubscriptionResponse struct {
// 	RequestedSubscriptionID int                                    `json:"RequestedSubscriptionID"`
// 	Entries                 []TSWAPI_GetSubscriptionResponse_Entry `json:"Entries"`
// }

// type TSWAPI_TrainController struct {
// }

// func (controller *TSWAPI_TrainController) Subscribe() {
// 	// http://localhost:31270/get/CurrentFormation/0.ObjectClass
// }

// func (controller *TSWAPI_TrainController) SendControlValue(direct_control_name string, value float64) {
// }
