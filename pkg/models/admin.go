package models

import (
	"github.com/Saakhr/Web-proj/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

func (a *Admin) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashed)
	return nil
}

func (a *Admin) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
}

func GetAdminByEmail(email string) (*Admin, error) {
	admin := &Admin{}
	err := database.DB.QueryRow(`
		SELECT id, email, password, full_name 
		FROM admins WHERE email = ?`, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Password,
		&admin.FullName)
	return admin, err
}
