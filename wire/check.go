package wire

import (
	"github.com/walnuts1018/crd-preview/backend/infra/postgres"
	"github.com/walnuts1018/crd-preview/backend/usecase"
)

var _ usecase.ManifestSouceRepo = &postgres.PostgresClient{}
var _ usecase.CRDRepo = &postgres.PostgresClient{}
var _ usecase.SchemaRepo = &postgres.PostgresClient{}
