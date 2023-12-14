package models

import (
	"database/sql"
	"fmt"
	"github.com/mihailtudos/photosharer/errors"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (service *GalleryService) Create(title string, userID int) (*Gallery, error) {
	// TODO: add validation
	gallery := Gallery{Title: title, UserID: userID}

	row := service.DB.QueryRow(`
		INSERT INTO galleries (user_id, title) 
		VALUES ($1, $2) returning id;
	`, gallery.UserID, gallery.Title)

	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByID(id int) (*Gallery, error) {
	// TODO: add validation
	gallery := Gallery{
		ID: id,
	}

	row := service.DB.QueryRow(`
		SELECT title, user_id
		    FROM galleries
			WHERE id = $1;
	`, gallery.ID)

	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("query gallery: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	// TODO: add validation
	row, err := service.DB.Query(`
		SELECT id, title
		    FROM galleries
			WHERE user_id = $1;
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("query gallieries by user: %w", err)
	}

	var galleries []Gallery

	for row.Next() {
		gallery := Gallery{
			UserID: userID,
		}

		err := row.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query gallery by user id: %w", err)
		}

		galleries = append(galleries, gallery)
	}

	if row.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}

	return galleries, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	// TODO: validation
	_, err := service.DB.Exec(`
		UPDATE galleries
			SET title = $2
		WHERE id = $1;
	`, gallery.ID, gallery.Title)

	if err != nil {
		return fmt.Errorf("update galleries: %w", err)
	}

	return nil
}

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM galleries
			WHERE id = $1;`, id)

	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	return nil
}
