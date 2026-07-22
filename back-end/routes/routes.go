package routes

import (
	"net/http"

	"back-end/controllers"
	"back-end/middleware"
)

func Route(mux *http.ServeMux) {
	// Authentification
	mux.HandleFunc("POST /register", middleware.AuthMiddleware(controllers.RegisterHandler))
	mux.HandleFunc("POST /login", middleware.AuthMiddleware(controllers.LoginHandler))
	mux.HandleFunc("POST /logout", middleware.AuthMiddleware(controllers.LogoutHandler))

	// Users
	mux.HandleFunc("GET /users", middleware.AuthMiddleware(controllers.GetUsersHandler))
	mux.HandleFunc("GET /users/{id}", middleware.AuthMiddleware(controllers.GetUserByIDHandler))

	// Posts
	mux.HandleFunc("GET /posts", controllers.GetPostsHandler)
	mux.HandleFunc("POST /posts", middleware.AuthMiddleware(controllers.CreatePostHandler))
	//-------//
	mux.HandleFunc("GET /posts/category", controllers.GetPostsByCategoryHandler)

	//Categories
	mux.HandleFunc("GET /categories", middleware.AuthMiddleware(controllers.GetCategoryHandler))

	/*
		add comment  POST /api/posts/:id/comments
		POST /api/posts/:id/like
		POST /api/posts/:id/dislike
		GET /api/posts?mine=true


	*/
}
