package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestPostgres struct {
	db *pgxpool.Pool
}

func NewRequestPostgres(db *pgxpool.Pool) *RequestPostgres {
	return &RequestPostgres{db: db}
}

func (r *RequestPostgres) Save(ctx context.Context, req []byte) (int, error) {
	var id int

	query, args, err := psql.Insert(requestTable).
		Columns("data").
		Values(req).
		ToSql()

	if err != nil {
		return 0, err
	}

	query += " RETURNING id"
	row := r.db.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
