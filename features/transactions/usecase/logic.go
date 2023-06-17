package usecase

import (
	"POS-PointofSales/features/transactions"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type transactionLogic struct {
	t         transactions.Repository
	validator *validator.Validate
}

func New(r transactions.Repository) transactions.UseCase {
	return &transactionLogic{
		t:         r,
		validator: validator.New(),
	}
}

// GetHistoryTransaction implements transactions.UseCase.
func (tl *transactionLogic) GetHistoryTransaction(userID uint, limit int, offset int, search string, fromDate time.Time, toDate time.Time) ([]transactions.ItemCore, int, error) {
	result, totaldata, err := tl.t.SelectHistoryTransaction(userID, limit, offset, search, fromDate, toDate)
	if err != nil {
		log.Error("failed to find all restock", err.Error())
		return []transactions.ItemCore{}, totaldata, errors.New("internal server error")
	}

	return result, totaldata, nil
}

// GetItemById implements transactions.UseCase.
func (tl *transactionLogic) GetItemByOrderId(orderID string) (transactions.ItemCore, error) {
	result, err := tl.t.SelectItemByOrderId(orderID)
	if err != nil {
		log.Error("failed to find item", err.Error())
		return transactions.ItemCore{}, errors.New("internal server error")
	}

	return result, nil
}

// AddPayments implements transactions.UseCase.
func (tl *transactionLogic) AddPayments(userID uint, newTransaction transactions.Core) (transactions.Core, error) {
	err := tl.validator.Struct(newTransaction)
	if err != nil {
		log.Error("validation error:", err.Error())
		return transactions.Core{}, err
	}

	res, err := tl.t.InsertPayments(userID, newTransaction)
	if err != nil {
		log.Error("failed on calling add product query")
		if strings.Contains(err.Error(), "open") {
			log.Error("errors occurs on opening picture file")
			return transactions.Core{}, errors.New("product photo are not allowed")
		} else if strings.Contains(err.Error(), "upload file in path") {
			log.Error("upload file in path are error")
			return transactions.Core{}, errors.New("cannot upload file in path")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on add product")
			return transactions.Core{}, errors.New("data is up to date")
		}
		return transactions.Core{}, err
	}
	return res, nil
}

// AddTransactions implements transactions.UseCase.
func (tl *transactionLogic) AddTransactions(userID uint, newDetailTransaction transactions.ItemCore) (transactions.ItemCore, error) {
	err := tl.validator.Struct(newDetailTransaction)
	if err != nil {
		log.Error("validation error:", err.Error())
		return transactions.ItemCore{}, err
	}

	res, err := tl.t.InsertDetailTransactions(userID, newDetailTransaction)
	if err != nil {
		log.Error("failed on calling add product query")
		if strings.Contains(err.Error(), "open") {
			log.Error("errors occurs on opening picture file")
			return transactions.ItemCore{}, errors.New("product photo are not allowed")
		} else if strings.Contains(err.Error(), "upload file in path") {
			log.Error("upload file in path are error")
			return transactions.ItemCore{}, errors.New("cannot upload file in path")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on add product")
			return transactions.ItemCore{}, errors.New("data is up to date")
		}
		return transactions.ItemCore{}, err
	}
	return res, nil
}
