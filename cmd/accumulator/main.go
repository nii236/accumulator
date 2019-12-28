package main

import (
	"accumulator"
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/run"
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
	masterKey := flag.String("master-key", "parliamentary-rutherfordium-goldeneye", "master key for encryption at rest")
	jwtSecret := flag.String("jwt-secret", "contractible-roasted-mollusk", "jwt secret")
	dbseed := flag.Bool("db-seed", false, "Seed fake data")
	stepMinutes := flag.Int("step-minutes", 5, "Step time between scrapes")
	rootPath := flag.String("root-path", "./web/dist", "Path of the webapp")
	serverAddr := flag.String("server-addr", ":8081", "Address to host on")
	loadBalancerAddr := flag.String("loadbalancer-addr", ":8080", "Address to host on")
	flag.Parse()

	conn, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	boil.SetDB(conn)

	if *dbseed {
		fmt.Println("Seeding accumulator system...")
		err = accumulator.Seed()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	fmt.Println("Booting up accumulator system...")
	g := &run.Group{}
	ctx, cancel := context.WithCancel(context.Background())
	g.Add(func() error {
		d, err := accumulator.NewDarer(*masterKey)
		if err != nil {
			return err
		}
		return accumulator.RunServer(ctx, conn, *serverAddr, *jwtSecret, d, accumulator.NewLogToStdOut("server", "0.0.1", false))
	}, func(err error) {
		fmt.Println(err)
		cancel()
	})
	g.Add(func() error {
		return accumulator.RunLoadBalancer(ctx, conn, *loadBalancerAddr, *serverAddr, *rootPath, accumulator.NewLogToStdOut("lb", "0.0.1", false))
	}, func(err error) {
		fmt.Println(err)
		cancel()
	})
	g.Add(func() error {
		return accumulator.RunAttendanceTracker(ctx, *stepMinutes, accumulator.NewLogToStdOut("attendance", "0.0.1", false))
	}, func(err error) {
		fmt.Println(err)
		cancel()
	})
	log.Fatalln(g.Run())
}
