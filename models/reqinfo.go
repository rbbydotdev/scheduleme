package models

type AcceptType string

const (
	HTMLKey AcceptType = "text/html"
	JSONKey AcceptType = "application/json"
)

type RequestInfo struct {
	Accept AcceptType
}

func (ri RequestInfo) ContextKey() string {
	return "RequestInfo"
}

func (ri *RequestInfo) IsJSON() bool {
	return ri.Accept == JSONKey
}
func (ri *RequestInfo) IsHTML() bool {
	return ri.Accept == HTMLKey
}

// // TODO figure out if JSON or HTML default is best to use
// func (ri *RequestInfo) DetermineAndSetAccept(acceptStr string) {
// 	if acceptStr == "application/json" {
// 		ri.Accept = JSONKey
// 	} else {
// 		ri.Accept = HTMLKey
// 	}
// }

func NewRequestInfo() *RequestInfo {
	return &RequestInfo{
		Accept: HTMLKey,
	}
}
