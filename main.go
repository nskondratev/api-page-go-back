package main

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/nskondratev/api-page-go-back/conf"
	"github.com/nskondratev/api-page-go-back/db"
	eventStore "github.com/nskondratev/api-page-go-back/events/store"
	"github.com/nskondratev/api-page-go-back/gql"
	"github.com/nskondratev/api-page-go-back/handler"
	"github.com/nskondratev/api-page-go-back/logger"
	"github.com/nskondratev/api-page-go-back/pages"
	pageStore "github.com/nskondratev/api-page-go-back/pages/store"
	"github.com/nskondratev/api-page-go-back/router"
	"github.com/nskondratev/api-page-go-back/ws"
)

func main() {
	r := router.New()
	c, err := conf.GetAppConfig()

	if err != nil {
		r.Logger.Fatal(err)
	}

	baseGroup := r.Group(c.BaseUrl)
	apiGroup := baseGroup.Group("/api")

	l := logger.New(r.Logger)

	d, err := db.NewGorm(&db.MysqlDBConfig{
		ConnectionString: c.DBConnectionString,
	})

	if err != nil {
		r.Logger.Fatal(err)
	}

	ps := pageStore.NewGorm(&pageStore.GormConfig{
		DB:     d,
		Logger: l,
	})

	es := eventStore.NewGorm(&eventStore.GormConfig{
		DB:     d,
		Logger: l,
	})

	wsHub := ws.NewHub()
	go wsHub.Run()

	gqlHub := gql.NewGraphQLHub()

	gqlHub.AddType(pages.GraphQLType)

	if err := pages.RegisterGraphQLQueries(ps, gqlHub); err != nil {
		r.Logger.Fatalf("Error while registering graphql queries from pages: %s", err.Error())
	}

	if err := gqlHub.Compile(); err != nil {
		r.Logger.Fatalf("Error while compiling graphql schema: %s", err.Error())
	}

	h := handler.New(&handler.Config{
		Logger:     l,
		PageStore:  ps,
		EventStore: es,
		WsHub:      wsHub,
		GraphQLHub: gqlHub,
	})
	h.Register(apiGroup, baseGroup)

	r.Server.Addr = c.Addr
	r.Logger.Fatal(gracehttp.Serve(r.Server))
}
