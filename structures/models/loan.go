package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type LoanApplication struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    ReferenceNumber string    `gorm:"type:varchar(20);uniqueIndex" json:"reference_number"`

    // Applicant Details
    FullName     string `gorm:"type:varchar(100)" json:"full_name"`
    DOB          string `gorm:"type:varchar(20)" json:"dob"`
    Email        string `gorm:"type:varchar(100)" json:"email"`
    MobileNumber string `gorm:"type:varchar(20)" json:"mobile_number"`
    Address      string `json:"address"`
    City         string `json:"city"`
    State        string `json:"state"`
    Pincode      string `gorm:"type:varchar(20)" json:"pincode"`

    // Loan Configuration
    LoanTrack       string  `gorm:"type:varchar(50)" json:"loan_track"`
    ProductCategory string  `gorm:"type:varchar(50)" json:"product_category"`
    PrincipalAmount float64 `gorm:"not null" json:"principal_amount"`
    TenureMonths    int     `gorm:"not null" json:"tenure_months"`
    InterestRate    float64 `json:"interest_rate"`
    EstimatedEMI    float64 `json:"estimated_emi"`

    ApplicationStatus string `gorm:"type:varchar(20);default:'UNDER_REVIEW'" json:"status"`

    KYC           KYCDocuments     `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE;" json:"kyc_documents"`
    FinancialDocs FinancialDetails `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE;" json:"financial_details"`

    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` 
}

// 2. KYC Table (Strictly Identity & Security)
type KYCDocuments struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	LoanID           uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"loan_id"` // uniqueIndex forces 1:1 Match
	LiveSelfiePath   string    `json:"live_selfie_path"`
	AadhaarFrontPath string    `json:"aadhaar_front_path"`
	AadhaarBackPath  string    `json:"aadhaar_back_path"`
	PanCardPath      string    `json:"pan_card_path"`
}

// 3. Financial Table (Strictly Income & Bank files)
type FinancialDetails struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	LoanID                uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"loan_id"` // uniqueIndex forces 1:1 Match
	EmploymentStatus      string    `gorm:"type:varchar(50)" json:"employment_status"`
	MonthlyIncome         float64   `json:"monthly_income"`
	BankStatementPath     string    `json:"bank_statement_path"`
	PropertyAgreemntPath  string    `json:"property_agreemnt_path"`
	IncomeProofPath       string    `json:"income_proof_path"`
}
func (l *LoanApplication) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return
}
func (k *KYCDocuments) BeforeCreate(tx *gorm.DB) (err error) {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return
}
func (f *FinancialDetails) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}