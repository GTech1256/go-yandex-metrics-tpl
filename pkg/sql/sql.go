package sql

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"time"
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

	//err = s.MigrateUp("postgres", host)
	//if err != nil {
	//	fmt.Println("ERROR:", err)
	//}

	return s, nil
}

func (q SQL) MigrateUp(driverName, dataSourceName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)

	// or m.Step(2) if you want to explicitly set the number of migrations to run
	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
