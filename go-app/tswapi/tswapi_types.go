package tswapi

type PropertyName = string

type TSWAPI_ListResponse_Node struct {
	Name string `json:"Name"`
}

type TSWAPI_ListResponse struct {
	Nodes []TSWAPI_ListResponse_Node `json:"Nodes"`
}

type TSWAPI_GetCurrentDrivableActorSubscriptionResponse_Control struct {
	Identifier             string
	PropertyName           string
	CurrentValue           float64
	CurrentNormalizedValue float64
}

type TSWAPI_GetCurrentDrivableActorSubscriptionResponse struct {
	ObjectClass string
	Controls    map[PropertyName]TSWAPI_GetCurrentDrivableActorSubscriptionResponse_Control
}
