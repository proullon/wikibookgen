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
	"github.com/proullon/wikibookgen/pkg/cache"
	"github.com/proullon/wikibookgen/pkg/classifier"
	"github.com/proullon/wikibookgen/pkg/clusterer"
	"github.com/proullon/wikibookgen/pkg/generator"
	"github.com/proullon/wikibookgen/pkg/loader"
	"github.com/proullon/wikibookgen/pkg/orderer"

	. "github.com/proullon/wikibookgen/api/model"
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
			Name:   "dumpdb-fr",
			Value:  "wikipedia_fr",
			Usage:  "Database name",
			EnvVar: "DUMP_DB_FR",
		},
		cli.StringFlag{
			Name:   "dumpdb-en",
			Value:  "wikipedia_en",
			Usage:  "Database name",
			EnvVar: "DUMP_DB_EN",
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
	maxconn := c.Int("db-max-conn")

	db, err := openDB(host, dbname, usr, sslRootCert, sslClientCert, sslClientKey, maxconn)
	if err != nil {
		log.Errorf("OpenDB %s: %s", dbname, err)
	}
	dbfr, err := openDB(host, c.String("dumpdb-fr"), usr, sslRootCert, sslClientCert, sslClientKey, maxconn)
	if err != nil {
		log.Errorf("OpenDB %s: %s", c.String("dumpdb-fr"), err)
	}
	dben, err := openDB(host, c.String("dumpdb-en"), usr, sslRootCert, sslClientCert, sslClientKey, maxconn)
	if err != nil {
		log.Errorf("OpenDB %s: %s", c.String("dumpdb-en"), err)
	}

	if c.String("logfile") != "" {
		f, err := os.OpenFile(c.String("logfile"), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return fmt.Errorf("log: %s", err)
		}
		log.SetOutput(f)
	}

	loadermap := make(map[string]Loader)

	dbloader := loader.NewDBLoader(dbfr)
	cacheloader := cache.NewLocalCacheLoader(dbloader)
	loadermap["fr"] = cacheloader

	dbloader = loader.NewDBLoader(dben)
	cacheloader = cache.NewLocalCacheLoader(dbloader)
	loadermap["en"] = cacheloader

	// TODO: classifier should havve loader as LoadGraph argument
	cla, err := classifier.NewV1(loadermap["fr"])
	if err != nil {
		return fmt.Errorf("classifier.NewV1: %s", err)
	}

	gen := generator.NewV1(db, cla, clusterer.NewV1(), orderer.NewV1(), loadermap)
	wg := wikibookgen.New(db, gen)

	ctx := context.WithValue(context.Background(), "wg", wg)
	err = wikibookgen.ListenAndServe(ctx, "8080")
	if err != nil {
		return fmt.Errorf("ListenAndServe: %s", err)
	}
	return nil
}

func openDB(host, dbname, usr, sslRootCert, sslClientCert, sslClientKey string, maxconn int) (*sql.DB, error) {

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
		return nil, fmt.Errorf("sql.Open: %s", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %s", err)
	}
	db.SetMaxOpenConns(maxconn)
	db.SetMaxIdleConns(-1)

	log.Infof("Connected to %s/%s", host, dbname)
	return db, nil
}
