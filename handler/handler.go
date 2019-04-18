package handler

import (
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/gql"
	"github.com/nskondratev/api-page-go-back/logger"
	"github.com/nskondratev/api-page-go-back/pages"
	"github.com/nskondratev/api-page-go-back/ws"
)

type Handler struct {
	logger     logger.Logger
	pageStore  pages.Store
	eventStore events.Store
	wsHub      ws.IHub
	gqlHub     *gql.GraphQLHub
}

type Config struct {
	Logger     logger.Logger
	PageStore  pages.Store
	EventStore events.Store
	WsHub      ws.IHub
	GraphQLHub *gql.GraphQLHub
}

func New(hc *Config) *Handler {
	return &Handler{
		logger:     hc.Logger,
		pageStore:  hc.PageStore,
		eventStore: hc.EventStore,
		wsHub:      hc.WsHub,
		gqlHub:     hc.GraphQLHub,
	}
}
