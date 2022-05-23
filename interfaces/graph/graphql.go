package graph

import (
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	_ "github.com/urfave/cli/v2"
	_ "golang.org/x/tools/go/ast/astutil"
	_ "golang.org/x/tools/go/packages"
	_ "golang.org/x/tools/imports"
	"net/http"
	"sealway-strava/interfaces/graph/generated"
	"sealway-strava/interfaces/rest"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
)

type GraphqlApi struct {
	Resolvers *Resolver
	*rest.DefaultApi
}

func (server *GraphqlApi) RegisterGraphQl() *handler.Server {
	serverName := "graphql"

	schema := generated.NewExecutableSchema(generated.Config{Resolvers: server.Resolvers})
	srv := handler.NewDefaultServer(schema)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	endpoint := server.Prefix(serverName, "/query")
	server.Router.Handle(server.Prefix(serverName, "/"), PlaygroundHandler("GraphQL playground", endpoint))
	server.Router.Handle(endpoint, srv)

	return srv
}
