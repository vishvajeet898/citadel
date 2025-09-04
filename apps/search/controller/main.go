package controller

import (
	"github.com/Orange-Health/citadel/apps/search/service"
)

type Search struct {
	SearchService service.SearchServiceInterface
}

func InitSearchController() *Search {
	return &Search{
		SearchService: service.InitializeSearchService(),
	}
}
