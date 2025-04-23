package models

import (
	"fmt"

	"github.com/Saakhr/Web-proj/pkg/database"
)

type Projects struct{
  ID int `json:"id"`
  Title string `json:"title"`
  Description string `json:"description"`
}
func CreateProject(a *Projects) error{
  _,err:=database.DB.Exec(`
    INSERT INTO projects (title,description)
    VALUES (?, ?);
    `,
    a.Title, a.Description)
  return err
}

func DeleteProject(id int) error {
    // Prepare the DELETE query
    query := `
        DELETE FROM projects
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

func GetStudentProjectsRec(studentID int) ([]Projects,error){
  rows,err:=database.DB.Query(`
SELECT p.id, p.title, p.description 
FROM projects p
WHERE p.id NOT IN (
    SELECT wp.project_id 
    FROM student_project_wishlist wp
    WHERE wp.student_id = $1  
    )`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
  var projects []Projects
  for rows.Next(){
    var p Projects
    err:=rows.Scan(
      &p.ID,
      &p.Title,
      &p.Description)
		if err != nil {
			return nil, err
		}
    projects=append(projects, p)
  }
  return projects,nil

}

func GetProjects() ([]Projects,error){
  rows,err:=database.DB.Query(`
    SELECT id, title, description
		FROM projects;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
  var projects []Projects
  for rows.Next(){
    var p Projects
    err:=rows.Scan(
      &p.ID,
      &p.Title,
      &p.Description)
		if err != nil {
			return nil, err
		}
    projects=append(projects, p)
  }
  return projects,nil

}
