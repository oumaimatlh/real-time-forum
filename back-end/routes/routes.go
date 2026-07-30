package routes

import (
	"net/http"

	"back-end/controllers"
	"back-end/middleware"
)

func Route(mux *http.ServeMux) {
	// Authentification
	mux.Handle("POST /api/register", middleware.RateLimiter(controllers.RegisterHandler))
	mux.Handle("POST /api/login", middleware.RateLimiter(controllers.LoginHandler))
	mux.Handle("POST /api/logout", middleware.RateLimiter(middleware.AuthMiddleware(controllers.LogoutHandler)))

	// Users
	mux.Handle("GET /api/users", middleware.RateLimiter(middleware.AuthMiddleware(controllers.GetUsersHandler)))
	mux.Handle("GET /api/users/{id}", middleware.RateLimiter(middleware.AuthMiddleware(controllers.GetUserByIDHandler)))

	// Categories
	mux.Handle("GET /api/categories", middleware.RateLimiter(controllers.GetCategoryHandler))

	// Posts
	mux.Handle("POST /api/posts", middleware.RateLimiter(middleware.AuthMiddleware(controllers.CreatePostHandler)))
	mux.Handle("GET /api/posts", middleware.RateLimiter(controllers.GetPostsHandler))
	// Post Reactions
	mux.Handle("POST /api/posts/{id}/like", middleware.RateLimiter(middleware.AuthMiddleware(controllers.LikePostHandler)))
	mux.Handle("POST /api/posts/{id}/dislike", middleware.RateLimiter(middleware.AuthMiddleware(controllers.DislikePostHandler)))

	// Comments
	mux.Handle("POST /api/posts/{id}/comment", middleware.RateLimiter(middleware.AuthMiddleware(controllers.CreateCommentPostHandler)))

	// Comment Reactions
	mux.Handle("POST /api/comments/{id}/like", middleware.RateLimiter(middleware.AuthMiddleware(controllers.LikeCommentHandler)))
	mux.Handle("POST /api/comments/{id}/dislike", middleware.RateLimiter(middleware.AuthMiddleware(controllers.DislikeCommentHandler)))

	// Filter
	mux.Handle("GET /api/posts/filter", middleware.RateLimiter(middleware.AuthMiddleware(controllers.FilterPostsHandler)))

	// Conversation
	mux.Handle("GET /api/conversations", middleware.RateLimiter(middleware.AuthMiddleware(controllers.GetConversationsHandler)))

	// CHAT
	mux.Handle("GET /api/chat", middleware.RateLimiter(middleware.AuthMiddleware(controllers.ChatHandler)))

	// Historique d Message entre 2 user
	mux.Handle("GET /api/history", middleware.RateLimiter(middleware.AuthMiddleware(controllers.GetConversationMessagesHandler)))
}
