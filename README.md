# Wikibookgen

`Wikibookgen` is a book generator using Wikipedia as source.

## Build

Copy and edit `.dev.conf.example` to `.dev.conf`, build Go binaries and docker image:
```sh
cp .dev.conf.example .dev.conf
make install
make package
```

## Usage

`Wikibookgen` requires a running `CockroachDB` instance, you can see the setup at [wikipedia-to-cockroachdb](https://github.com/proullon/wikipedia-to-cockroachdb)



