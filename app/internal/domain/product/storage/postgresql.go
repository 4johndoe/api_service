package storage

import (
	sq "github.com/Masterminds/squirrel"
	"production_service/pkg/logging"
)

type ProductStorage struct {
	QueryBuiler sq.StatementBuilderType
	client      PostgreSQLClient
	logger      *logging.Logger
}
