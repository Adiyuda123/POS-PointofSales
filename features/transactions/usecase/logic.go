package usecase

import (
	"POS-PointofSales/features/transactions"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
)

type transactionLogic struct {
	t transactions.Repository
}

func New(r transactions.Repository) transactions.UseCase {
	return &transactionLogic{
		t: r,
	}
}

// CreateTransactionsIfNotExists implements transactions.UseCase.
func (tl *transactionLogic) CreateTransactionsIfNotExists(userId uint, id uint, customer string) (transactions.Core, error) {
	existingTransaction, err := tl.t.SelectTransactionById(id)
	if err != nil {
		return transactions.Core{}, err
	}

	if existingTransaction.Id != 0 {
		return existingTransaction, nil
	}

	newTransaction := transactions.Core{
		Id:       id,
		Status:   "Pending",
		Customer: customer,
		UserID:   userId,
	}

	createdTransaction, err := tl.t.InsertTransactions(userId, newTransaction)
	if err != nil {
		return transactions.Core{}, err
	}

	if createdTransaction.Id != id {
		return transactions.Core{}, errors.New("failed to create transaction with the provided ID")
	}

	return createdTransaction, nil
}

// GetTransactionById implements transactions.UseCase.
func (tl *transactionLogic) GetTransactionById(id uint) (transactions.Core, error) {
	result, err := tl.t.SelectTransactionById(id)
	if err != nil {
		log.Error("failed to find transaction", err.Error())
		return transactions.Core{}, errors.New("internal server error")
	}

	return result, nil
}

// GetTotalAmount implements transactions.UseCase.
func (tl *transactionLogic) GetTotalAmount(externalID string, customer string) (int, error) {
	amount, err := tl.t.GetTotalAmount(externalID, customer)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

// AddTransactions implements transactions.UseCase.
func (tl *transactionLogic) AddTransactions(userID uint, newDetailTransaction transactions.DetailCore) (transactions.DetailCore, error) {
	res, err := tl.t.InsertDetailTransactions(userID, newDetailTransaction)
	if err != nil {
		log.Error("failed on calling add product query")
		if strings.Contains(err.Error(), "open") {
			log.Error("errors occurs on opening picture file")
			return transactions.DetailCore{}, errors.New("product photo are not allowed")
		} else if strings.Contains(err.Error(), "upload file in path") {
			log.Error("upload file in path are error")
			return transactions.DetailCore{}, errors.New("cannot upload file in path")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on add product")
			return transactions.DetailCore{}, errors.New("data is up to date")
		}
		return transactions.DetailCore{}, err
	}
	return res, nil
}

// AddTransactions implements transactions.UseCase.
func (tl *transactionLogic) CreateTransactions(userId uint, newTransaction transactions.Core) (transactions.Core, error) {
	res, err := tl.t.InsertTransactions(userId, newTransaction)
	if err != nil {
		return transactions.Core{}, err
	}
	return res, nil
}
