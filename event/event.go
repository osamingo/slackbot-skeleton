package event

import (
	"context"

	"github.com/slack-go/slack"
)

type (
	Func   func(ctx context.Context, body []byte) (string, []slack.MsgOption, error)
	Router struct {
		m map[string]Func
	}
	Execution struct {
		eventName string
		channel   string
		fn        Func
	}
)

func NewRouter(es ...Execution) *Router {
	m := make(map[string]Func, len(es))
	for _, e := range es {
		m[e.eventName] = e.fn
	}

	return &Router{
		m: m,
	}
}

func (r *Router) GetFunc(eventName string) Func {
	return r.m[eventName]
}
