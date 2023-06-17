package transactions

import (
	// transactiondetails "POS-PointofSales/features/transactionDetails"

	"time"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID          string
	ExternalID  string `validate:"required"`
	Amount      int    `validate:"required"`
	QRString    string
	CallbackURL string `validate:"required"`
	Type        string
	Status      string
	Created     string
	Updated     string
	Customer    string
	ItemID      uint
	UserID      uint
	OrderID     string
}

type ItemCore struct {
	Id       uint
	SubTotal int
	Customer string `validate:"required"`
	UserID   uint
	UserName string
	OrderID  string `validate:"required"`
	Status   string
	Details  []DetailCore `validate:"required"`
}

type DetailCore struct {
	Id        uint
	ItemID    uint
	ProductID uint `validate:"required"`
	Quantity  int  `validate:"required,min=0"`
	Total     int
	Price     int
}

type Handler interface {
	AddTransactions() echo.HandlerFunc
	AddPayments() echo.HandlerFunc
	GetHistoryTransactionHandler() echo.HandlerFunc
}

type UseCase interface {
	AddTransactions(userID uint, newDetailTransaction ItemCore) (ItemCore, error)
	AddPayments(userID uint, newTransaction Core) (Core, error)
	GetItemByOrderId(orderID string) (ItemCore, error)
	GetHistoryTransaction(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]ItemCore, int, error)
}

type Repository interface {
	InsertDetailTransactions(userID uint, inputDetail ItemCore) (ItemCore, error)
	InsertPayments(userID uint, newTransaction Core) (Core, error)
	SelectItemByOrderId(orderID string) (ItemCore, error)
	SelectHistoryTransaction(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]ItemCore, int, error)
}
