package main

import (
	"database/sql"
	"net"

	"github.com/icrxz/fullcycle-grpc/internal/database"
	"github.com/icrxz/fullcycle-grpc/internal/pb"
	"github.com/icrxz/fullcycle-grpc/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDB := database.NewCategory(db)
	service := service.NewCategoryService(*categoryDB)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(list); err != nil {
		panic(err)
	}
}
