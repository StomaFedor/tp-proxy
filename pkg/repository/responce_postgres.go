package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResponcePostgres struct {
	db *pgxpool.Pool
}

func NewResponcePostgres(db *pgxpool.Pool) *ResponcePostgres {
	return &ResponcePostgres{db: db}
}

func (r *ResponcePostgres) Save(ctx context.Context, req []byte) (int, error) {
	var id int

	query, args, err := psql.Insert(responceTable).
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
