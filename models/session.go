package models

import (
	"context"
	"scheduleme/frame"
	"scheduleme/util"
	"scheduleme/values"
)

type SessionInfo struct {
	IsAdmin     bool
	UserID      values.ID
	Flash       string
	RedirectURL string
	CSRFToken   string
	State       string
}

type SessionInterface interface {
	IsLoggedIn() bool
	RotateCSRFToken() string
}

func (ses SessionInfo) ContextKey() string {
	return "SessionInfo"
}

func (ses *SessionInfo) IsLoggedIn() bool {
	return ses.UserID != 0
}

func generateCSRFToken() string {
	return util.RandomStr(32)
}

func (ses *SessionInfo) RotateCSRFToken() string {
	oldToken := ses.CSRFToken
	ses.CSRFToken = generateCSRFToken()
	return oldToken
}

// State
func StateFromContext(ctx context.Context) frame.CtxState {
	return frame.CtxState(frame.FromContext[SessionInfo](ctx).State)
}

func NewSession() *SessionInfo {
	s := &SessionInfo{
		IsAdmin: false,
	}
	s.RotateCSRFToken()
	return s
}
