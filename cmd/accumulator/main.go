package main

import (
	"accumulator"
	"context"
	"flag"
	"fmt"
	"github.com/nii236/vrchat-go/client"
	vrc "github.com/nii236/vrchat-go/client"
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
	vrcClient, err := vrc.NewClient(client.ReleaseAPIURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Booting up accumulator system...")
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
	g := &run.Group{}
	ctx, cancel := context.WithCancel(context.Background())

	g.Add(func() error {
		return accumulator.RunServer(ctx, conn, *serverAddr, vrcClient, accumulator.NewLogToStdOut("server", "0.0.1", false))
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
	log.Fatalln(g.Run())
}
