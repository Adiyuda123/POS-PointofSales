package helper

import (
	"POS-PointofSales/app/config"

	"github.com/google/uuid"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(amount int, email string) (string, error) {
	xendit.Opt.SecretKey = config.XenditSecretKey
	amountFloat := float64(amount)

	externalID := uuid.New().String()

	createInvoiceData := invoice.CreateParams{
		ExternalID:  externalID,
		Amount:      amountFloat,
		PayerEmail:  email,
		Description: "Invoice for POS",
	}

	inv, err := invoice.Create(&createInvoiceData)
	if err != nil {
		return "", err
	}

	return inv.InvoiceURL, nil
}
