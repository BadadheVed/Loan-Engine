package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name  string    `gorm:"type:varchar(100);not null" json:"name"`
	Email string    `gorm:"type:varchar(255);not null" json:"email"`
	Age   int       `json:"age"`

	MonthlyIncome float64 `gorm:"type:numeric(12,2);not null;index:idx_users_income" json:"monthly_income"`
	CreditScore   int     `gorm:"not null;index:idx_users_credit_score" json:"credit_score"`

	EmploymentStatus string `gorm:"type:varchar(50)" json:"employment_status"`
	CreatedAt        time.Time
}
