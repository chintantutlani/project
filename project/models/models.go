package models

type User struct {
	ID          uint `gorm:"primaryKey"`
	FirstName   string
	LastName    string
	CompanyName string
	Address     string
	City        string
	County      string
	Postal      string
	Phone       string
	Email       string
	Web         string
}
