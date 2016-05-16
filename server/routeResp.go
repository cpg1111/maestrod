package server

// RouteResp is a struct representing the response to a get on a route that displays subroutes
type RouteResp struct {
	Routes []string `json:"routes"`
}
