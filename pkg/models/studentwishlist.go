package models

import (
	"fmt"

	"github.com/Saakhr/Web-proj/pkg/database"
)

type StudentWishlist struct {
	ID          []int        `json:"id"`
	StudentName Student    `json:"studentname"`
	Projects    []Projects `json:"projects"`
}

func CreateWish(studentID, ProjectID int) error {
	_, err := database.DB.Exec(`
		INSERT INTO student_project_wishlist (student_id, project_id)
		VALUES (?, ?)`,
		studentID, ProjectID)
	return err
}

func DeleteWishlistItem(studentID int) error {
	// Prepare the DELETE query
	query := `
        DELETE FROM student_project_wishlist
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
		return fmt.Errorf("no wishlist found with ID %d", studentID)
	}

	return nil
}

func DeleteStudentWishlistItem(id,studentID int) error {
	// Prepare the DELETE query
	query := `
        DELETE FROM student_project_wishlist
        WHERE id = $1 AND student_id = $2;`

	// Execute the query
	result, err := database.DB.Exec(query, id, studentID)
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

func GetAllWishlists() ([]StudentWishlist, error) {
	query := `SELECT 
    s.id,
    ts.id,
    ts.first_name,
    ts.last_name,
    p.id,
    p.Title
FROM 
    ` + "`student_project_wishlist`" + `s
JOIN 
    students ts ON s.student_id = ts.id
JOIN 
  projects  p ON s.project_id = p.id;
  `
	rows, err := database.DB.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Map to store Group data
	groupsMap := make(map[int]StudentWishlist)

	for rows.Next() {
		var wishlistID int
		var studentID int
		var studentName string
		var studentName2 string
		var projectID int
		var projectName string

		err := rows.Scan(&wishlistID, &studentID, &studentName, &studentName2, &projectID, &projectName)
		if err != nil {
			panic(err.Error())
		}

		// Check if Group exists in map, if not create a new Group
		group, ok := groupsMap[studentID]
		if !ok {
			group = StudentWishlist{StudentName: Student{
				ID:        studentID,
				FirstName: studentName,
				LastName:  studentName2,
			}}
		}

		// Append the course to the Group's Modules
		group.Projects = append(group.Projects, Projects{Title: projectName, ID: projectID})
		group.ID = append(group.ID, wishlistID)

		// Update the Group in the map
		groupsMap[studentID] = group
	}

	// Convert map to slice and return
	groups := make([]StudentWishlist, 0, len(groupsMap))
	for _, group := range groupsMap {
		groups = append(groups, group)
	}
	return groups, nil
}

func GetStudentWishlists(id int) (StudentWishlist, error) {
	query := `SELECT 
    s.id,
    ts.id,
    ts.first_name,
    ts.last_name,
    p.id,
    p.Title
FROM 
    ` + "`student_project_wishlist`" + `s
JOIN 
    students ts ON s.student_id = ts.id
JOIN 
  projects  p ON s.project_id = p.id
  WHERE ts.id = $1;
  `
	rows, err := database.DB.Query(query, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Map to store Group data
	groupsMap := make(map[int]StudentWishlist)
  var wishlist StudentWishlist

	for rows.Next() {
		var wishlistID int
		var studentID int
		var studentName string
		var studentName2 string
		var projectID int
		var projectName string

		err := rows.Scan(&wishlistID, &studentID, &studentName, &studentName2, &projectID, &projectName)
		if err != nil {
			panic(err.Error())
		}

		// Check if Group exists in map, if not create a new Group
		group, ok := groupsMap[studentID]
		if !ok {
			group = StudentWishlist{StudentName: Student{
				ID:        studentID,
				FirstName: studentName,
				LastName:  studentName2,
			}}
		}

		// Append the course to the Group's Modules
		group.Projects = append(group.Projects, Projects{Title: projectName, ID: projectID})
		group.ID = append(group.ID, wishlistID)

		// Update the Group in the map
		groupsMap[studentID] = group
    wishlist=group
	}

	// Convert map to slice and return
	groups := make([]StudentWishlist, 0, len(groupsMap))
	for _, group := range groupsMap {
		groups = append(groups, group)
	}
  
	return wishlist, nil
}
