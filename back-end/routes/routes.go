package routes

import (
	"net/http"

	"back-end/controllers"
	"back-end/middleware"
)

func Route(mux *http.ServeMux) {
	// Authentification
	mux.HandleFunc("POST /api/register", controllers.RegisterHandler)
	mux.HandleFunc("POST /api/login", controllers.LoginHandler)
	mux.HandleFunc("POST /api/logout", middleware.AuthMiddleware(controllers.LogoutHandler))

	// Users
	mux.HandleFunc("GET /api/users", middleware.AuthMiddleware(controllers.GetUsersHandler))
	mux.HandleFunc("GET /api/users/{id}", middleware.AuthMiddleware(controllers.GetUserByIDHandler))

	// Categories
	mux.HandleFunc("GET /api/categories", controllers.GetCategoryHandler)

	// Posts
	mux.HandleFunc("POST /api/posts", middleware.AuthMiddleware(controllers.CreatePostHandler))
	mux.HandleFunc("GET /api/posts", controllers.GetPostsHandler)
	// Post Reactions
	mux.HandleFunc("POST /api/posts/{id}/like", middleware.AuthMiddleware(controllers.LikePostHandler))
	mux.HandleFunc("POST /api/posts/{id}/dislike", middleware.AuthMiddleware(controllers.DislikePostHandler))

	// Comments
	mux.HandleFunc("POST /api/posts/{id}/comment", middleware.AuthMiddleware(controllers.CreateCommentPostHandler))

	// Comment Reactions
	mux.HandleFunc("POST /api/comments/{id}/like", middleware.AuthMiddleware(controllers.LikeCommentHandler))
	mux.HandleFunc("POST /api/comments/{id}/dislike", middleware.AuthMiddleware(controllers.DislikeCommentHandler))

	// Filter
	mux.HandleFunc("GET /api/posts/filter", middleware.AuthMiddleware(controllers.FilterPostsHandler))
}
