package search

import (
	"context"
	"sync"

	search "github.com/fvdveen/mu2-proto/go/proto/search"
)

// Wrapper is a wrapper around a search service
type Wrapper interface {
	search.SearchServiceHandler
	SetService(search.SearchServiceHandler) error
}

// NewWrapper creates a new wrapper
func NewWrapper(s search.SearchServiceHandler) Wrapper {
	return &wrapper{
		s: s,
	}
}

type wrapper struct {
	s  search.SearchServiceHandler
	mu sync.RWMutex
}

func (w *wrapper) SetService(s search.SearchServiceHandler) error {
	w.mu.Lock()
	w.s = s
	w.mu.Unlock()
	return nil
}

func (w *wrapper) Search(ctx context.Context, req *search.SearchRequest, res *search.SearchResponse) error {
	return w.s.Search(ctx, req, res)
}
