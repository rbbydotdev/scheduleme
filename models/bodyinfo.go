package models

type BodyInfo struct {
	User  *User //TODO do these need to be pointers?
	Event *Event
}

func (bi BodyInfo) ContextKey() string {
	return "BodyInfo"
}

func NewBodyInfo() *BodyInfo {
	return &BodyInfo{
		Event: &Event{},
		User:  &User{},
	}
}
