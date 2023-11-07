package sql

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type SQL struct {
	DB *pgx.Conn
}

func NewSQL(host string) (*SQL, error) {
	db, err := pgx.Connect(context.Background(), host)
	if err != nil {
		return nil, err
	}

	// TODO: CONNECTION POOL
	// https://github.com/jackc/pgx/wiki/Getting-started-with-pgx#using-a-connection-pool

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	s := &SQL{
		DB: db,
	}

	err = s.MigrateDown(host)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	err = s.MigrateUp(host)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	return s, nil
}

func (q SQL) MigrateUp(dataSourceName string) error {
	fmt.Println("Migration Up Started")
	m, err := migrate.New(
		"file://internal/server/config/db/migrations",
		dataSourceName)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}

	fmt.Println("Migration Up Ended")
	return nil
}

func (q SQL) MigrateDown(dataSourceName string) error {
	fmt.Println("Migration Down Started")
	m, err := migrate.New(
		"file://internal/server/config/db/migrations",
		dataSourceName)
	if err != nil {
		return err
	}
	err = m.Down()
	if err != nil {
		return err
	}

	fmt.Println("Migration Down Ended")
	return nil
}
