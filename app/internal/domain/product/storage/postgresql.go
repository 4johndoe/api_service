package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"production_service/internal/domain/product/model"
	"production_service/pkg/client/postgresql"
	db "production_service/pkg/client/postgresql/model"
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

func (s *ProductStorage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *ProductStorage) All(ctx context.Context) ([]model.Product, error) {
	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		//Column("rating").
		//Column("category_id").
		//Column("specification").
		Column("created_at").
		Column("updated_at").
		From(schema + "." + table) //.ToSql()

	// TODO filtering and sorting

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err := db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	list := make([]model.Product, 0)

	for rows.Next() {
		p := model.Product{}
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.ImageID, &p.Price, &p.CurrencyID, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}
