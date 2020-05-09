package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/proullon/wikibookgen/api"
	//"github.com/proullon/wikibookgen/pkg/cache"
	"github.com/proullon/wikibookgen/pkg/classifier"
	"github.com/proullon/wikibookgen/pkg/clusterer"
	"github.com/proullon/wikibookgen/pkg/generator"
	"github.com/proullon/wikibookgen/pkg/loader"
	"github.com/proullon/wikibookgen/pkg/orderer"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "crdb.example.com",
			Usage:  "CockroachDB host",
			EnvVar: "CRDB_HOST",
		},
		cli.StringFlag{
			Name:   "dbname",
			Value:  "wikipedia",
			Usage:  "Database name",
			EnvVar: "DB_NAME",
		},
		cli.StringFlag{
			Name:   "user",
			Value:  "wikipedia",
			Usage:  "Database user",
			EnvVar: "DB_USER",
		},
		cli.StringFlag{
			Name:   "ssl-root-cert",
			Value:  "certs/ca.crt",
			Usage:  "Root SSL certificate",
			EnvVar: "SSL_ROOT_CERT",
		},
		cli.StringFlag{
			Name:   "ssl-client-key",
			Value:  "certs/client.wikipedia.key",
			Usage:  "Client SSL key",
			EnvVar: "SSL_CLIENT_KEY",
		},
		cli.StringFlag{
			Name:   "ssl-client-cert",
			Value:  "certs/client.wikipedia.crt",
			Usage:  "Client SSL certificate",
			EnvVar: "SSL_CLIENT_CERT",
		},
		cli.IntFlag{
			Name:   "db-max-conn",
			Value:  300,
			Usage:  "Maximum number of open connection to database",
			EnvVar: "DB_MAX_CONN",
		},
		cli.StringFlag{
			Name:   "logfile",
			Value:  "",
			Usage:  "Log destination",
			EnvVar: "LOGFILE",
		},
	}
	app.Action = start
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}

func start(c *cli.Context) error {
	host := c.String("host")
	dbname := c.String("dbname")
	usr := c.String("user")
	sslRootCert := c.String("ssl-root-cert")
	sslClientKey := c.String("ssl-client-key")
	sslClientCert := c.String("ssl-client-cert")

	dsn := fmt.Sprintf("postgresql://%s@%s:26257/%s?ssl=true&sslmode=require&sslrootcert=%s&sslkey=%s&sslcert=%s",
		usr,
		host,
		dbname,
		sslRootCert,
		sslClientKey,
		sslClientCert,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(c.Int("db-max-conn"))
	db.SetMaxIdleConns(c.Int("db-max-conn"))
	log.Infof("Connected to %s/%s\n", host, dbname)

	if c.String("logfile") != "" {
		f, err := os.OpenFile(c.String("logfile"), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		log.SetOutput(f)
	}

	/*
		dbloader := cache.NewDBLoader(db)
		cacheloader := cache.NewLocalCacheLoader(dbloader)

		cla, err := classifier.NewV1(cacheloader)
	*/
	loader, err := loader.NewFileLoader("/data/mathematiques.json")
	if err != nil {
		return err
	}

	cla, err := classifier.NewV1(loader)
	if err != nil {
		return err
	}

	gen := generator.NewV1(db, cla, clusterer.NewV1(), orderer.NewV1(db))
	wg := wikibookgen.New(db, gen)

	ctx := context.WithValue(context.Background(), "wg", wg)
	err = wikibookgen.ListenAndServe(ctx, "8080")
	if err != nil {
		return err
	}
	return nil
}
