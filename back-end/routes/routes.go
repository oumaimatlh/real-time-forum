package routes

import (
	"net/http"

	"back-end/controllers"
)

func Route(mux *http.ServeMux) {
	// Authentification
	mux.HandleFunc("POST /register", controllers.RegisterHandler)
	mux.HandleFunc("POST /login", controllers.LoginHandler)
	mux.HandleFunc("POST /logout", controllers.LogoutHandler)
	

	// //Users
	// mux.HandleFunc("GET /users"))
	// mux.HandleFunc("GET /users/{id}")

	// //Posts
	// mux.HandleFunc("GET /posts")
	// mux.HandleFunc("GET /posts/{id}")

	// //Categories
	// mux.HandleFunc("GET /categories")

	// //Comments
	// mux.HandleFunc("/posts/{id}/comments")
}
