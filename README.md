# go-my-mutex

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

Mutex lock based on MySQL user-level locks for Go. Can be used to acquire a lock on a resource exclusively across several running instances of the application.

## Install

```bash
$ go get -u -v github.com/vgarvardt/go-my-mutex
```

## PostgreSQL drivers

The store accepts an adapter interface that interacts with the DB. Adapter and implementations are extracted to separate package [`github.com/vgarvardt/go-pg-adapter`](https://github.com/vgarvardt/go-pg-adapter) for easier maintenance.

## Usage example

```go
package main

import (
	"context"
    "fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vgarvardt/go-my-mutex"
	"github.com/vgarvardt/go-pg-adapter/sqladapter"
)

func main() {
	db, _ := sql.Open("mysql", os.Getenv("DB_URI"))
    defer db.Close()
    conn, _ := db.Conn(context.Background())
    defer conn.Close()

	m, _ := mymutex.New(sqladapter.NewConn(conn))

	useExclusiveResource(m)
	
	if !useExclusiveResourceOrNoOp(m) {
		fmt.Println("Resource is busy, doing nothing")
	}
}

func useExclusiveResource(m *pgmutex.PgMutex) {
	lockName := "my-lock"
	_ := m.Lock(lockName)
	defer m.Unlock(lockName)

	// do something with resource exclusively across several instances
}

func useExclusiveResourceOrNoOp(m *pgmutex.PgMutex) bool {
	lockName := "my-try-lock"
	if success, _ := m.TryLock(lockName); !success {
		return false
	}
	defer m.Unlock(lockName)

	// do something with resource exclusively across several instances
}
```

## How to run tests

You will need running MySQL instance. E.g. the one running in docker and exposing a port to a host system

```bash
docker run --rm -p 3306:3306 -it -e MYSQL_ROOT_PASSWORD=mutex -e MYSQL_DATABASE=mutex mysql:5.7
```

Now you can run tests using the running PostgreSQL instance using `MY_URI` environment variable

```bash
MY_URI="root:mutex@tcp(localhost:3306)/mutex" go test -cover ./...
```

## MIT License

```
Copyright (c) 2019 Vladimir Garvardt
```

[Build-Status-Url]: https://travis-ci.org/vgarvardt/go-my-mutex
[Build-Status-Image]: https://travis-ci.org/vgarvardt/go-my-mutex.svg?branch=master
[codecov-url]: https://codecov.io/gh/vgarvardt/go-my-mutex
[codecov-image]: https://codecov.io/gh/vgarvardt/go-my-mutex/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/vgarvardt/go-my-mutex
[reportcard-image]: https://goreportcard.com/badge/github.com/vgarvardt/go-my-mutex
[godoc-url]: https://godoc.org/github.com/vgarvardt/go-my-mutex
[godoc-image]: https://godoc.org/github.com/vgarvardt/go-my-mutex?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
