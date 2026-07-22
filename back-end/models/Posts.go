package models

import (
	"database/sql"
	"fmt"
	"time"

	"back-end/database"
)

type Post struct {
	IdPost    int
	Title     string
	Content   string
	UserId    int
	NickName  string
	Comments  []Comments
	Likes     int
	Dislikes  int
	CreatedAt time.Time
}

func InsertPost(post Post) (int64, error) {
	fmt.Printf("UserId = %d\n", post.UserId)
	fmt.Printf("%+v\n", post)
	query := "INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, post.Title, post.Content, post.UserId)
	if err != nil {
		return 0, err
	}
	lastPostId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastPostId, nil
}

func GetAllPosts() ([]Post, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.user_id, posts.created_at, users.nickName
		FROM posts
		INNER JOIN users ON posts.user_id = users.id 
		ORDER BY posts.created_at DESC`

	return getPostsFromQuery(query)
}

func GetPostsByCategory(idcat int) ([]Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, p.created_at, u.nickName
		FROM posts p
		INNER JOIN post_category pc ON p.id = pc.post_id
		INNER JOIN users u ON p.user_id = u.id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC`

	return getPostsFromQuery(query, idcat)
}

func GetFilteredPosts(userID int, selectedCats []string, filterLikes, filterMyPosts bool) ([]Post, error) {
	seen := make(map[int]bool)
	allPosts := []Post{}

	addPosts := func(posts []Post) {
		for _, p := range posts {
			if !seen[p.IdPost] {
				seen[p.IdPost] = true
				allPosts = append(allPosts, p)
			}
		}
	}
	for _, cat := range selectedCats {
		if cat == "all" {
			continue
		}
		query := `
			SELECT DISTINCT p.id, p.title, p.content, p.user_id, p.created_at, u.nickName
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			INNER JOIN post_category pc ON p.id = pc.post_id
			WHERE pc.category_id = ?
			ORDER BY p.created_at DESC`

		posts, err := getPostsFromQuery(query, cat)
		if err != nil {
			return nil, err
		}
		addPosts(posts)
	}

	if filterLikes {
		query := `
			SELECT DISTINCT p.id, p.title, p.content, p.user_id, p.created_at, u.nickName
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			INNER JOIN likes_dislikes ld ON p.id = ld.post_id
			WHERE ld.user_id = ? AND ld.type = 'like'
			ORDER BY p.created_at DESC`

		posts, err := getPostsFromQuery(query, userID)
		if err != nil {
			return nil, err
		}
		addPosts(posts)
	}
	if filterMyPosts {
		query := `
			SELECT p.id, p.title, p.content, p.user_id, p.created_at, u.nickName
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			WHERE p.user_id = ?
			ORDER BY p.created_at DESC`

		posts, err := getPostsFromQuery(query, userID)
		if err != nil {
			return nil, err
		}
		addPosts(posts)
	}

	return allPosts, nil
}

//----------//

func scanPost(row *sql.Row) (Post, error) {
	var p Post
	err := row.Scan(&p.IdPost, &p.Title, &p.Content, &p.UserId, &p.CreatedAt, &p.NickName)
	return p, err
}

func GetPostByID(id int) (Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, p.created_at, u.nickName
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`
	return scanPost(database.DB.QueryRow(query, id))
}

func scanPostsRows(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		post := Post{}
		if err := rows.Scan(&post.IdPost, &post.Title, &post.Content,
			&post.UserId, &post.CreatedAt, &post.NickName); err != nil {
			return nil, err
		}

		comments, err := GetCommentsByPost(post.IdPost)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		likes, err := CountLikeDislikeByPost(post.IdPost, "like")
		if err != nil {
			return nil, err
		}
		post.Likes = likes

		dislikes, err := CountLikeDislikeByPost(post.IdPost, "dislike")
		if err != nil {
			return nil, err
		}
		post.Dislikes = dislikes

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func getPostsFromQuery(query string, args ...interface{}) ([]Post, error) {
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPostsRows(rows)
}
