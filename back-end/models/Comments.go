package models

import (
	"time"

	"back-end/database"
)

type Comments struct {
	IdComment int
	UserId    int
	PostId    int
	Content   string
	CreatedAt time.Time
	NickName  string
	Likes     int
	Dislikes  int
}

func InsertComment(c Comments) error {
	_, err := database.DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", c.UserId, c.PostId, c.Content)
	return err
}

func GetCommentsByPost(postID int) ([]Comments, error) {
	rows, err := database.DB.Query(`
                SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, u.nickName
                FROM comments c
                JOIN users u ON u.id = c.user_id
                WHERE c.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comments

	for rows.Next() {
		var c Comments
		err := rows.Scan(
			&c.IdComment,
			&c.UserId,
			&c.PostId,
			&c.Content,
			&c.CreatedAt,
			&c.NickName,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i, c := range comments {
		likees, _ := CountLikesByComments(c.IdComment, c.PostId, "like")
		dislikes, _ := CountLikesByComments(c.IdComment, c.PostId, "dislike")
		comments[i].Likes = likees
		comments[i].Dislikes = dislikes
	}
	return comments, nil
}

func GetCommentByID(commentID int) (Comments, error) {
	var comment Comments

	query := `
		SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, u.nickName
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.id = ?
	`

	err := database.DB.QueryRow(query, commentID).Scan(
		&comment.IdComment,
		&comment.UserId,
		&comment.PostId,
		&comment.Content,
		&comment.CreatedAt,
		&comment.NickName,
	)
	if err != nil {
		return Comments{}, err
	}

	likes, _ := CountLikesByComments(comment.IdComment, comment.PostId, "like")
	dislikes, _ := CountLikesByComments(comment.IdComment, comment.PostId, "dislike")

	comment.Likes = likes
	comment.Dislikes = dislikes

	return comment, nil
}

func ToggleCommentReaction(userId, postId, commentId int, reactionType string) error {
	currentType, err := CheckReactionByUser(userId, postId, &commentId)
	if err != nil {
		return err
	}

	if currentType == reactionType {
		return DeleteReaction(userId, postId, &commentId)
	}

	if currentType != "" {
		query := `
            UPDATE likes_dislikes
            SET type = ?
            WHERE user_id = ?
            AND post_id = ?
            AND comment_id = ?
        `

		_, err := database.DB.Exec(query,
			reactionType,
			userId,
			postId,
			commentId,
		)

		return err
	}

	_, err = InsertReaction(Reaction{
		UserID:    userId,
		PostID:    postId,
		CommentID: &commentId,
		Type:      reactionType,
	})

	return err
}
