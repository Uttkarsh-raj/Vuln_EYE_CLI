package main

import (
	"fmt"
	"log"

	"github.com/Uttkarsh-raj/PS-1708/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Server...")
	server := gin.New() // New server
	server.Use(gin.Logger())
	routes.RegisterRoutes(server) // Register the Different routes to the server
	log.Fatal(server.Run(":3000"))
}
