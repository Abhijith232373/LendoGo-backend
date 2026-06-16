package services

import (
	"context" 
	"errors"
	"fmt"
	"math"
	"time"
	"lendogo-backend/internal/repositories"
	"lendogo-backend/structures/dto"
	"lendogo-backend/utils"
	"github.com/google/uuid"
)

type WalletService interface {
	GetBalance() (float64, error)
	GenerateRechargeOrder(amount float64) (map[string]interface{}, error)
	ProcessPaymentVerification(orderID, paymentID, signature string, amount float64) error
	DirectFund(amount float64) error
	ProcessDisbursal(req dto.DisburseLoanRequest) error
	GetUserBalance(userID string) (float64, error)
}
type walletServiceImpl struct {
	repo     repositories.WalletRepository
	producer *utils.KafkaProducer
}
func NewWalletService(repo repositories.WalletRepository, producer *utils.KafkaProducer) WalletService {
	return &walletServiceImpl{
		repo:     repo,
		producer: producer,
	}
}
func (s *walletServiceImpl) GetBalance() (float64, error) {
	return s.repo.GetSystemBalance()
}

func (s *walletServiceImpl) GenerateRechargeOrder(amount float64) (map[string]interface{}, error) {
	receipt := fmt.Sprintf("admin_rx_%d", time.Now().Unix())
	return utils.CreateRazorpayOrder(amount, receipt)
}

func (s *walletServiceImpl) ProcessPaymentVerification(orderID, paymentID, signature string, amount float64) error {
	isValid := utils.VerifyRazorpaySignature(orderID, paymentID, signature)
	if !isValid {
		return errors.New("invalid payment signature")
	}
	err := s.repo.CreditSystemWallet(amount)
	if err != nil {
		return err 
	}
	payload := map[string]interface{}{
		"type":           "SYSTEM_RECHARGE",
		"transaction_id": paymentID,
		"order_id":       orderID,
		"amount":         amount,
		"status":         "COMPLETED",
	}
	kafkaErr := s.producer.PublishEvent(context.Background(), "telemetry.payments", "SYSTEM_RECHARGE_SUCCESS", payload)
	if kafkaErr != nil {
		fmt.Println("KAFKA ERROR: Failed to publish system recharge event:", kafkaErr)
	} else {
		fmt.Println("KAFKA EVENT PUBLISHED: SYSTEM_RECHARGE_SUCCESS")
	}
	return nil
}
func (s *walletServiceImpl) DirectFund(amount float64) error {
	return s.repo.CreditSystemWallet(amount)
}

//  OUTBOUND CAPITAL (Loan Disbursements)
func (s *walletServiceImpl) ProcessDisbursal(req dto.DisburseLoanRequest) error {
	expectedNet := req.SanctionedAmt - req.ProcessingFee
	if math.Abs(req.NetPayout-expectedNet) > 0.01 {
		return errors.New("security alert: payout amount mismatch detected")
	}
	loanUUID, err := uuid.Parse(req.LoanID)
	if err != nil {
		return errors.New("invalid loan ID format")
	}
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	err = s.repo.ExecuteDisbursal(loanUUID, userUUID, req.NetPayout)
	if err != nil {
		return err 
	}

//  KAFKA TRIGGER: BLAST THE EVENT AFTER SUCCESSFUL DISBURSAL!
	payload := map[string]interface{}{
		"loan_id":   req.LoanID,
		"user_id":   req.UserID,
		"amount":    req.NetPayout,
		"status":    "DISBURSED",
		"timestamp": time.Now().Unix(),
	}
	kafkaErr := s.producer.PublishEvent(context.Background(), "telemetry.loans", "LOAN_DISBURSED", payload)
	if kafkaErr != nil {
		fmt.Println("KAFKA ERROR: Failed to publish loan disbursal event:", kafkaErr)
	} else {
		fmt.Println("KAFKA EVENT PUBLISHED: LOAN_DISBURSED")
	}
	return nil 
}

//  USER WALLET LOGIC
func (s *walletServiceImpl) GetUserBalance(userID string) (float64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	return s.repo.GetUserBalance(userUUID)
}