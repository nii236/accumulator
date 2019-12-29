package main

import (
	"accumulator"
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kelseyhightower/envconfig"
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

type Config struct {
	MasterKey        string `default:"9A1F3DE2BB279CB966CC1167BC6C538FDE97268E3EE5F581D918309409520AE3"`
	JWTSecret        string `default:"contractible-roasted-mollusk"`
	StepMinutes      int    `default:"5"`
	RootPath         string `default:"./web/dist"`
	ServerAddr       string `default:":8081"`
	LoadBalancerAddr string `default:":8080"`
}

func main() {
	dbseed := flag.Bool("db-seed", false, "Seed fake data")
	showConfig := flag.Bool("config", false, "Show config variables")

	c := &Config{}
	err := envconfig.Process("ACCUMULATOR", c)
	if err != nil {
		log.Fatal(err.Error())
	}
	flag.Parse()
	conn, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	boil.SetDB(conn)
	if *showConfig {
		envconfig.Usage("ACCUMULATOR", c)
		return
	}
	if *dbseed {
		fmt.Println("Seeding accumulator system...")
		err = accumulator.Seed(c.MasterKey)
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
		d, err := accumulator.NewDarer(c.MasterKey)
		if err != nil {
			return err
		}
		return accumulator.RunServer(ctx, conn, c.ServerAddr, c.JWTSecret, d, accumulator.NewLogToStdOut("server", "0.0.1", false))
	}, func(err error) {
		fmt.Println(err)
		cancel()
	})
	g.Add(func() error {
		return accumulator.RunLoadBalancer(ctx, conn, c.LoadBalancerAddr, c.ServerAddr, c.RootPath, accumulator.NewLogToStdOut("lb", "0.0.1", false))
	}, func(err error) {
		fmt.Println(err)
		cancel()
	})
	g.Add(func() error {
		d, err := accumulator.NewDarer(c.MasterKey)
		if err != nil {
			return err
		}
		return accumulator.RunAttendanceTracker(ctx, d, c.StepMinutes, accumulator.NewLogToStdOut("attendance", "0.0.1", false))
	}, func(err error) {
		fmt.Println(err)
		cancel()
	})
	log.Fatalln(g.Run())
}
