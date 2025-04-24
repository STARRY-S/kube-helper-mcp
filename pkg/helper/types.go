package helper

import (
	"encoding/json"
)

type Result interface {
	String() string
}

type listResult struct {
	Results []Result `json:"results"`
}

func (l *listResult) Add(r Result) {
	l.Results = append(l.Results, r)
}

func (l *listResult) String() string {
	if l == nil {
		return ""
	}
	if len(l.Results) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(l.Results)
	return string(b)
}
