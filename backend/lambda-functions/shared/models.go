package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User model
type User struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`
	Email            string    `gorm:"type:varchar(255);not null" json:"email"`
	Age              int       `json:"age"`
	MonthlyIncome    float64   `gorm:"type:numeric(12,2);not null;index:idx_users_income" json:"monthly_income"`
	CreditScore      int       `gorm:"not null;index:idx_users_credit_score" json:"credit_score"`
	EmploymentStatus string    `gorm:"type:varchar(50)" json:"employment_status"`
	CreatedAt        time.Time
}

// LoanProduct model
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

// Match model
type Match struct {
	ID              uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"match_id"`
	UserID          uuid.UUID   `gorm:"type:uuid;not null;index;column:user_id" json:"user_id"`
	User            User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"-"`
	ProductID       uuid.UUID   `gorm:"type:uuid;not null;index;column:product_id" json:"product_id"`
	LoanProduct     LoanProduct `gorm:"constraint:OnDelete:CASCADE;foreignKey:ProductID" json:"-"`
	MatchConfidence bool        `gorm:"default:false" json:"match_confidence"`
	IsNotified      bool        `gorm:"default:false" json:"is_notified"`
	MatchedAt       time.Time   `gorm:"autoCreateTime" json:"matched_at"`
	Reason          string      `gorm:"type:text" json:"reason"`
}

// BatchResult - shared result type for worker pool
type BatchResult struct {
	Inserted  int
	Attempted int
}
