package database

import "log"

func TableCreation() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
  			nickName VARCHAR UNIQUE NOT NULL,
			firstName VARCHAR  NOT NULL,
  			lastName VARCHAR  NOT NULL,
  			email VARCHAR UNIQUE NOT NULL,
			age  INTEGER NOT NULL, 
			gender VARCHAR NOT NULL CHECK(gender IN ('male', 'female')),
  			password VARCHAR UNIQUE NOT NULL,
  			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,

		`CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY,
			title VARCHAR,
			content TEXT,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS category(
			id INTEGER PRIMARY KEY,
			name VARCHAR NOT NULL
	)`,

		`CREATE TABLE IF NOT EXISTS post_category(
			post_id INTEGER ,
			category_id INTEGER,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS comments(
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			content VARCHAR
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS likes_dislikes(		
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			post_id INTEGER,
			comment_id INTEGER,
			type 	TEXT NOT NULL CHECK(type IN ('like','dislike')),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS sessions(
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			token VARCHAR UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY,
			sender_id INTEGER NOT NULL,
			receiver_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
	)`,
	}

	for _, query := range queries {
		prep, err := DB.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}

		_, err = prep.Exec()
		if err != nil {
			log.Fatal(err)
		}
	}
}
