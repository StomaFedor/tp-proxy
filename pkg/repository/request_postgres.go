package repository

import (
	"context"
	"encoding/json"
	"strings"

	"tp-proxy/pkg/models"

	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5"
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

func (r *RequestPostgres) GetAll(ctx context.Context) ([]models.Request, error) {
	query, args, err := psql.
		Select("id, data").
		From(requestTable).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		var data string
		err = rows.Scan(&request.Id, &data)
		if err != nil {
			return nil, err
		}
		data = strings.ReplaceAll(data, "\\", "")
		err = json.Unmarshal([]byte(data), &request)
		requests = append(requests, request)
	}
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	return requests, nil
}

func (r *RequestPostgres) GetById(ctx context.Context, id int) (models.Request, error) {
	query, args, err := psql.
		Select("id, data").
		From(requestTable).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return models.Request{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	if err != nil {
		return models.Request{}, err
	}

	request, err := scanRequest(row)

	return request, err
}

func scanRequest(row pgx.Row) (models.Request, error) {
	var request models.Request
	var data string
	if err := row.Scan(&request.Id, &data); err != nil {
		return models.Request{}, err
	}
	data = strings.ReplaceAll(data, "\\", "")
	json.Unmarshal([]byte(data), &request)

	return request, nil
}
