package main

import (
	"context"
	"fmt"

	searchpb "github.com/fvdveen/mu2-proto/go/proto/search"
	"github.com/micro/cli"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/sirupsen/logrus"
)

func main() {
	var loc, term string

	s := grpc.NewService(
		micro.Flags(
			cli.StringFlag{
				Destination: &loc,
				Name:        "location",
				Usage:       "Search service location",
				Value:       "mu2.service.search",
			},
			cli.StringFlag{
				Destination: &term,
				Name:        "search-term",
				Usage:       "The term to search on",
				Value:       "Never gonna give you up",
			},
		),
	)

	s.Init()

	cl := searchpb.NewSearchService(loc, s.Client())

	res, err := cl.Search(context.Background(), &searchpb.SearchRequest{
		Name: term,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	fmt.Printf("%+v", *res.Video)
}
