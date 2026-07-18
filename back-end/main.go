package main

import (
	"back-end/database"
	"back-end/routes"
	"fmt"
	"log"
	"net/http"
)
func main(){
	database.DBinit()
	mux := http.NewServeMux()
	routes.Route(mux)
	
	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}