package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"production_service/internal/domain/product/model"
	"production_service/pkg/logging"
)

type ProductStorage struct {
	queryBuilder sq.StatementBuilderType
	client       PostgreSQLClient
	logger       *logging.Logger
}

func NewProductStorage(client PostgreSQLClient, logger *logging.Logger) ProductStorage {
	return ProductStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
	}
}

const (
	schema = "public"
	table  = "product"
)

func (s *ProductStorage) All(ctx context.Context) ([]model.Product, error) {
	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		Column("rating").
		Column("category_id").
		Column("specification").
		Column("created_at").
		Column("updated_at").
		From(schema + "." + table) //.ToSql()

	// TODO filtering and sorting

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	fmt.Println(sql, args) // log sql

	return nil, nil
}
