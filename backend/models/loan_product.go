package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type LoanProduct struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"product_id"`
	BankName         string         `gorm:"type:varchar(100);not null" json:"bank_name"`
	ProductName      string         `gorm:"type:varchar(255);not null" json:"product_name"`
	InterestRate     string         `gorm:"type:varchar(50)" json:"interest_rate"`
	MinCreditScore   int            `gorm:"default:0" json:"min_credit_score"`
	MinMonthlyIncome float64        `gorm:"type:numeric(12,2);default:0" json:"min_monthly_income"`
	RawCriteria      datatypes.JSON `gorm:"type:jsonb" json:"raw_criteria"`
	ProductURL       string         `gorm:"type:text;uniqueIndex" json:"product_url"`
	Age              int            `gorm:"constraint:check=age>=18;column:age;" json:"age"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at"`
}
