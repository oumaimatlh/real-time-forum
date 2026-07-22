package models

import "back-end/database"

type Category struct {
	IdCat int
	Name  string
}

func GetCategoryByID(id int) (Category, error) {
	var category Category
	query := `SELECT id, name FROM category WHERE id = ?`

	err := database.DB.QueryRow(query, id).Scan(&category.IdCat, &category.Name,)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func InsertPostCategory(postID int64, categoryID int) error {
	_, err := database.DB.Exec("INSERT INTO post_category (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	return err
}


func GetAllCategories() ([]Category, error) {
	categories := []Category{}
	query := "SELECT id, name FROM category"
	lignes, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	for lignes.Next() {
		category := Category{}
		err := lignes.Scan(&category.IdCat, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := lignes.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

