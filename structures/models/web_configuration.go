package models

import "time"

type WebConfiguration struct {
	ID uint `gorm:"primaryKey;autoIncrement:false;default:1" json:"-"`

	ApplyLoanEnabled     bool `gorm:"default:true" json:"apply_loan_enabled"`
	LoginEnabled         bool `gorm:"default:true" json:"login_enabled"`
	RegisterEnabled      bool `gorm:"default:true" json:"register_enabled"`
	ApplyJobEnabled      bool `gorm:"default:true" json:"apply_job_enabled"`
	
	ProfileUpdateEnabled bool `gorm:"default:true" json:"profile_update_enabled"`
	FeedbackEnabled      bool `gorm:"default:true" json:"feedback_enabled"`
	LoanHistoryEnabled   bool `gorm:"default:true" json:"loan_history_enabled"`
	RepayEnabled         bool `gorm:"default:true" json:"repay_enabled"`

	AutoPayEnabled       bool `gorm:"default:false" json:"auto_pay_enabled"`
	InternalScoreEnabled bool `gorm:"default:false" json:"internal_score_enabled"`
	CibilScoreEnabled    bool `gorm:"default:false" json:"cibil_score_enabled"`

	BlogEnabled             bool `gorm:"default:true" json:"blog_enabled"`
	ChatSupportEnabled      bool `gorm:"default:true" json:"chat_support_enabled"`
	FreeConsultationEnabled bool `gorm:"default:true" json:"free_consultation_enabled"`

	MinCreditScore   int     `gorm:"default:650" json:"min_credit_score"`
	BaseInterestRate float64 `gorm:"default:14.0" json:"base_interest_rate"`

	UpdatedAt time.Time `json:"updated_at"`
}