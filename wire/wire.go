//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/crd-preview/backend/config"
	"github.com/walnuts1018/crd-preview/backend/graph"
	"github.com/walnuts1018/crd-preview/backend/infra/postgres"
	"github.com/walnuts1018/crd-preview/backend/router"
	"github.com/walnuts1018/crd-preview/backend/usecase"
)

func CreateRouter(
	ctx context.Context,
	cfg config.Config,
) (*gin.Engine, error) {
	wire.Build(
		ConfigSet,
		postgresSet,
		usecase.NewUsecase,
		resolverSet,
		router.NewRouter,
	)

	return &gin.Engine{}, nil
}

var ConfigSet = wire.FieldsOf(new(config.Config),
	"PSQLDSN",
)

var resolverSet = wire.NewSet(
	graph.NewResolver,
	wire.Bind(new(graph.ResolverRoot), new(*graph.Resolver)),
)

var postgresSet = wire.NewSet(
	postgres.NewPostgres,
	wire.Bind(new(usecase.ManifestSouceRepo), new(*postgres.PostgresClient)),
	wire.Bind(new(usecase.CRDRepo), new(*postgres.PostgresClient)),
	wire.Bind(new(usecase.SchemaRepo), new(*postgres.PostgresClient)),
)
