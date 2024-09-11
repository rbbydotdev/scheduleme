package body_info

import (
	"scheduleme/domain/event"
)

type BodyInfo struct {
	Event *event.Event
}

func (bi BodyInfo) ContextKey() string {
	return "BodyInfo"
}

func NewBodyInfo() *BodyInfo {
	return &BodyInfo{
		Event: &event.Event{},
	}
}
