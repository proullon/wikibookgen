module github.com/proullon/wikibookgen

go 1.13

require (
	github.com/gorilla/mux v1.7.4
	github.com/lib/pq v1.4.0
	github.com/proullon/graph v0.0.0-00010101000000-000000000000
	github.com/proullon/workerpool v0.0.0-20200429190315-8cc98e318cde
	github.com/sirupsen/logrus v1.5.0
	github.com/urfave/cli v1.22.4
)

replace github.com/proullon/graph => ../graph
