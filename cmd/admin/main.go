package main

import (
	"flag"
	"fmt"

	"accumulator/bindata"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	migrate_bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
)

func connect() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("sqlite3", "./accumulator.db")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func main() {
	dbversion := flag.Bool("db-version", false, "Get the DB version")
	dbmigrate := flag.Bool("db-migrate", false, "Migrate DB")
	dbdrop := flag.Bool("db-drop", false, "Drop DB")
	flag.Parse()

	conn, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	boil.SetDB(conn)
	if *dbversion {
		fmt.Println("Getting DB version...")
		v, d, err := Version(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Version: %d, Dirty: %v\n", v, d)
		return
	}
	if *dbmigrate {
		fmt.Println("Migrating accumulator system...")
		err = Migrate(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	if *dbdrop {
		fmt.Println("Dropping accumulator system...")
		err = Drop(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

}

func newMigrateInstance(conn *sqlx.DB) (*migrate.Migrate, error) {
	s := migrate_bindata.Resource(bindata.AssetNames(),
		func(name string) ([]byte, error) {
			return bindata.Asset(name)
		})
	d, err := migrate_bindata.WithInstance(s)
	if err != nil {
		return nil, fmt.Errorf("bindata instance: %w", err)
	}
	dbDriver, err := sqlite3.WithInstance(conn.DB, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("db instance: %w", err)
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("migrate instance: %w", err)
	}
	return m, nil
}
func Migrate(conn *sqlx.DB) error {
	m, err := newMigrateInstance(conn)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = m.Up()
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
func Drop(conn *sqlx.DB) error {
	m, err := newMigrateInstance(conn)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = m.Drop()
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
func Version(conn *sqlx.DB) (uint, bool, error) {
	m, err := newMigrateInstance(conn)
	if err != nil {
		return 0, false, fmt.Errorf("migrate: %w", err)
	}
	v, d, err := m.Version()
	if err != nil {
		return 0, false, fmt.Errorf("migrate: %w", err)
	}
	return v, d, nil
}
