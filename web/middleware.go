package web

import (
	"context"
	"net/http"
)

type MidFunc func(handler HandlerFunc) HandlerFunc

type Chain struct {
	mw []MidFunc
}

func NewChain(mw ...MidFunc) *Chain {
	return &Chain{
		mw: mw,
	}
}

func (c *Chain) Then(handler HandlerFunc) HandlerFunc {
	if handler == nil {
		handler = func(ctx context.Context, r *http.Request) Encoder {
			return nil
		}
	}
	for i := len(c.mw) - 1; i >= 0; i-- {
		handler = c.mw[i](handler)
	}
	return handler
}

func (c *Chain) Append(mw ...MidFunc) *Chain {
	newMW := make([]MidFunc, len(c.mw)+len(mw))
	newMW = append(newMW, c.mw...)
	newMW = append(newMW, mw...)
	return &Chain{
		mw: newMW,
	}
}
