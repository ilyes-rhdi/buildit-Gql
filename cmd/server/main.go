package main

import (
	"github.com/ilyes-rhdi/buildit-Gql/internal/database"
	"github.com/ilyes-rhdi/buildit-Gql/internal/server"
	"github.com/ilyes-rhdi/buildit-Gql/pkg/redis"
)

func main() {
	s := server.NewServer(":8080")
	database.InitDB()
	redis.Connect()
	s.Run()
}
