package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/v-starostin/goph-keeper/internal/handler"
	"github.com/v-starostin/goph-keeper/internal/service"
	"github.com/v-starostin/goph-keeper/internal/storage"
	"github.com/v-starostin/goph-keeper/pkg/pb"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migration", "postgres", instance)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println(err)
		return
	}

	authStorage := storage.New(db)
	authService := service.NewAuth(authStorage, []byte("secret"))
	authHandler := handler.New(authService)

	server := grpc.NewServer()
	pb.RegisterAuthServer(server, authHandler)

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = server.Serve(l); err != nil {
		fmt.Println(err)
		return
	}
}
