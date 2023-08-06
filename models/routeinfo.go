package models

import "scheduleme/values"

type RouteInfo struct {
	Availability values.DateSlots
	Event        Event
	Events       Events
	User         User
	Auth         Auth
	APIKey       APIKey
	Offset       int
	Page         int
	Filter       string
}

func (ri RouteInfo) ContextKey() string {
	return "RouteInfo"
}

func NewRouteInfo() *RouteInfo {
	return &RouteInfo{}
}
