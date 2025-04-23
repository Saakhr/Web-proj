package models

import (
	"fmt"

	"github.com/Saakhr/Web-proj/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type Student struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
}

func (a *Student) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashed)
	return nil
}

func (a *Student) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
}

func (a *Student) CreateStudent()error{
  if err := a.HashPassword(); err != nil {
    return err
  }

	_, err := database.DB.Exec(`
		INSERT INTO students (email, first_name, last_name, password)
		VALUES (?, ?, ?, ?)`,
		a.Email, a.FirstName, a.LastName, a.Password)
	return err
}

func DeleteStudent(studentID int) error {
    // Prepare the DELETE query
    query := `
        DELETE FROM students
        WHERE id = $1;`

    // Execute the query
    result, err := database.DB.Exec(query, studentID)
    if err != nil {
        return fmt.Errorf("failed to delete student: %v", err)
    }

    // Check if any rows were affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check affected rows: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no student found with ID %d", studentID)
    }

    return nil
}

func GetStudents() ([]Student,error){
  rows,err:=database.DB.Query(`
    SELECT id, email, first_name, last_name
		FROM students;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
  var projects []Student
  for rows.Next(){
    var p Student
    err:=rows.Scan(
      &p.ID,
      &p.Email,
      &p.FirstName,
      &p.LastName)
		if err != nil {
			return nil, err
		}
    projects=append(projects, p)
  }
  return projects,nil

}

func GetStudentByEmail(email string) (*Student, error) {
	student := &Student{}
	err := database.DB.QueryRow(`
		SELECT id, email, password, first_name, last_name 
		FROM students WHERE email = ?`, email).Scan(
		&student.ID,
		&student.Email,
		&student.Password,
		&student.FirstName,
    &student.LastName)
	return student, err
}
