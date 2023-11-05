package sql

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type SQL struct {
	DB *pgx.Conn
}

func NewSQL(host string) (*SQL, error) {
	db, err := pgx.Connect(context.Background(), host)
	if err != nil {
		return nil, err
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()
	//if err = db.Ping(ctx); err != nil {
	//	return nil, err
	//}

	return &SQL{
		DB: db,
	}, nil
}
