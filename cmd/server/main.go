package main

import (
	_ "github.com/ilyes-rhdi/buildit-Gql/docs"
	"github.com/ilyes-rhdi/buildit-Gql/internal/database"
	"github.com/ilyes-rhdi/buildit-Gql/internal/server"
)

func main() {
	s := server.NewServer(":8080")
	database.InitDB()
	s.Run()
}
