package transactions

import (
	// transactiondetails "POS-PointofSales/features/transactionDetails"

	"github.com/labstack/echo/v4"
)

type Core struct {
	Id         uint
	ExternalID string
	Status     string
	InvoiceURL string
	Amount     int
	Customer   string
	UserID     uint
	Details    []DetailCore
}

type DetailCore struct {
	Id            uint
	ExternalID    string
	ProductID     uint
	Quantity      int
	Total         int
	Customer      string
	UserID        uint
	TransactionID uint
}

type Handler interface {
	CreateTransactions() echo.HandlerFunc
	AddTransactions() echo.HandlerFunc
	// GetTotalAmount() echo.HandlerFunc
}

type UseCase interface {
	CreateTransactions(userId uint, newTransaction Core) (Core, error)
	AddTransactions(userID uint, newDetailTransaction DetailCore) (DetailCore, error)
	GetTotalAmount(externalID string, customer string) (int, error)
	GetTransactionById(id uint) (Core, error)
	CreateTransactionsIfNotExists(userId uint, id uint, customer string) (Core, error)
}

type Repository interface {
	InsertTransactions(userId uint, input Core) (Core, error)
	InsertDetailTransactions(userID uint, inputDetail DetailCore) (DetailCore, error)
	GetTotalAmount(externalID string, customer string) (int, error)
	SelectTransactionById(id uint) (Core, error)
	// GetTransactionByExternalID(externalID string) (Core, error)
}
