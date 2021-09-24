package graph

import (
	_ "github.com/urfave/cli/v2"
	_ "golang.org/x/tools/go/ast/astutil"
	_ "golang.org/x/tools/go/packages"
	_ "golang.org/x/tools/imports"

	"sealway-strava/api"
	"sealway-strava/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type GraphqlApi struct {
	Resolvers *Resolver
	*api.DefaultApi
}

func (server *GraphqlApi) RegisterGraphQl() *handler.Server {
	serverName := "graphql"

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: server.Resolvers}))

	server.Router.Handle(server.Prefix(serverName, "/"), playground.Handler("GraphQL playground", "/query"))
	server.Router.Handle(server.Prefix(serverName, "/query"), srv)

	return srv
}
