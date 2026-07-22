package models

import (
	"database/sql"

	"back-end/database"
)

type Reaction struct {
	ID        int
	UserID    int
	PostID    int
	CommentID *int
	Type      string
}

func InsertReaction(reaction Reaction) (int64, error) {
	query := "INSERT INTO likes_dislikes (user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, reaction.UserID, reaction.PostID, reaction.CommentID, reaction.Type)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func DeleteReaction(userId int, postID int, commentID *int) error {
	query := `DELETE FROM likes_dislikes WHERE user_id = ? AND post_id = ? AND comment_id IS ?`
	_, err := database.DB.Exec(query, userId, postID, commentID)
	return err
}

func CheckReactionByUser(userId int, postID int, commentID *int) (string, error) {
	query := `SELECT type FROM likes_dislikes WHERE likes_dislikes.user_id = ? AND likes_dislikes.post_id = ? AND likes_dislikes.comment_id IS ?`
	row := database.DB.QueryRow(query, userId, postID, commentID)
	var reactionType string
	err := row.Scan(&reactionType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return reactionType, nil
}

func GetReactionByUser(userId, postId int) (string, error) {
	reactionType := ""
	query := "SELECT type FROM likes_dislikes WHERE user_id = ? AND post_id = ? AND comment_id IS NULL LIMIT 1"
	err := database.DB.QueryRow(query, userId, postId).Scan(&reactionType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return reactionType, nil
}

func UpdateReactionType(userId, postId int, reactionType string) error {
	query := "UPDATE likes_dislikes SET type = ? WHERE user_id = ? AND post_id = ? AND comment_id IS NULL"
	_, err := database.DB.Exec(query, reactionType, userId, postId)
	return err
}

func TogglePostReaction(userId, postId int, reactionType string) error {
	currentType, err := GetReactionByUser(userId, postId)
	if err != nil {
		return err
	}

	if currentType == reactionType {
		return DeleteReaction(userId, postId, nil)
	}

	if currentType != "" {
		return UpdateReactionType(userId, postId, reactionType)
	}

	_, err = InsertReaction(Reaction{
		UserID:    userId,
		PostID:    postId,
		CommentID: nil,
		Type:      reactionType,
	})
	return err
}

func CountLikeDislikeByPost(postId int, Type string) (int, error) {
	count := 0
	query := "SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND comment_id IS NULL AND type = ?"
	err := database.DB.QueryRow(query, postId, Type).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CountLikesByComments(commentID int, post_id int, Type string) (int, error) {
	count := 0
	query := `SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND comment_id = ? AND type = ?`
	err := database.DB.QueryRow(query, post_id, commentID, Type).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
