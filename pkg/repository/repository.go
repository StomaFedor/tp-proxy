package repository

import (
	"context"
	"tp-proxy/pkg/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Request interface {
	Save(ctx context.Context, req []byte) (int, error)
	GetAll(ctx context.Context) ([]models.Request, error)
	GetById(ctx context.Context, id int) (models.Request, error)
}

type Responce interface {
	Save(ctx context.Context, req []byte) (int, error)
}

type Repository struct {
	Request  Request
	Responce Responce
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Request:  NewRequestPostgres(db),
		Responce: NewResponcePostgres(db),
	}
}
