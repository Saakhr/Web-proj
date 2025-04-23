package models

import (
	"fmt"
	"time"

	"github.com/Saakhr/Web-proj/pkg/database"
)

type Announcement struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Display  string    `json:"display"` // general, computer_science, physics, chemistry, math
	DateTime time.Time `json:"datetime"`
}

func CreateAnnouncement(a *Announcement) error {
	_, err := database.DB.Exec(`
		INSERT INTO announcements (title, content, display, datetime)
		VALUES (?, ?, ?, ?)`,
		a.Title, a.Content, a.Display, a.DateTime)
	return err
}
func DeleteAnnouncement(id int) error {
    // Prepare the DELETE query
    query := `
        DELETE FROM announcements
        WHERE id = $1;`

    // Execute the query
    result, err := database.DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete student: %v", err)
    }

    // Check if any rows were affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check affected rows: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no student found with ID %d", id)
    }

    return nil
}

func GetAllAnnouncements() ([]Announcement, error) {
	rows, err := database.DB.Query(`
		SELECT id, title, content, display, datetime 
		FROM announcements 
		ORDER BY datetime DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []Announcement
	for rows.Next() {
		var a Announcement
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.Display,
			&a.DateTime)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func GetAnnouncementsByDepartment(dept string) ([]Announcement, error) {
	rows, err := database.DB.Query(`
		SELECT id, title, content, display, datetime 
		FROM announcements 
		WHERE display = ? 
		ORDER BY datetime DESC`, dept)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []Announcement
	for rows.Next() {
		var a Announcement
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.Display,
			&a.DateTime)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, nil
}
