package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./school.db")
	if err != nil {
		return err
	}

	// Enable foreign key constraints
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	return createTables()
}

func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS admins (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			full_name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS announcements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			display TEXT NOT NULL CHECK(display IN ('general', 'computer_science', 'physics', 'chemistry', 'math')),
			datetime DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS student_project_wishlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id INTEGER NOT NULL,
			project_id INTEGER NOT NULL,
			FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
