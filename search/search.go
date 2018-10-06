package search

import (
	"context"
	"fmt"
	"net/http"

	search "github.com/fvdveen/mu2-proto/go/proto/search"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type service struct {
	s *youtube.Service
}

// New creates a new seatch service
func New(t string) (search.SearchServiceHandler, error) {
	c := &http.Client{
		Transport: &transport.APIKey{Key: t},
	}

	s, err := youtube.New(c)
	if err != nil {
		return nil, fmt.Errorf("create youtube client: %v", err)
	}

	return &service{
		s: s,
	}, nil
}

func (s *service) Search(ctx context.Context, req *search.SearchRequest, res *search.SearchResponse) error {
	res.Video = &search.Video{}

	resp, err := s.s.Search.List("id,snippet").
		Context(ctx).
		Q(req.Name).
		MaxResults(1).
		Type("video").
		Do()
	if err != nil {
		return fmt.Errorf("call youtube: %v", err)
	}

	if len(resp.Items) == 0 {
		return nil
	}

	res.Video.Id = resp.Items[0].Id.VideoId
	res.Video.Name = resp.Items[0].Snippet.Title
	res.Video.Thumbnail = resp.Items[0].Snippet.Thumbnails.Default.Url
	res.Video.Url = fmt.Sprintf("https://www.youtube.com/watch?v=%s", resp.Items[0].Id.VideoId)

	return nil
}
