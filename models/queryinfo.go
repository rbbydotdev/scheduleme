package models

import (
	"net/url"
	"time"
)

type QueryInfo struct {
	AvailQuery *AvailQuery
}

type AvailQuery struct {
	StartTime time.Time
	EndTime   time.Time
}

func (aq *AvailQuery) UpdatesQueryInfo(qi *QueryInfo, a *AvailQuery) {
	qi.AvailQuery = a
}
func (aq *AvailQuery) New() *AvailQuery {
	return &AvailQuery{}
}
func (aq *AvailQuery) Parse(v url.Values) (err error) {
	if (v.Get("start") == "") && (v.Get("end") == "") {
		aq.StartTime = time.Now()
		aq.EndTime = time.Now().AddDate(0, 1, 0)
		return
	}
	s, err := time.Parse(time.RFC3339, v.Get("start"))
	if err != nil {
		return
	}
	e, err := time.Parse(time.RFC3339, v.Get("end"))
	if err != nil {
		return
	}
	aq.StartTime = s
	aq.EndTime = e
	return
}

func (qi QueryInfo) ContextKey() string {
	return "QueryInfo"
}
