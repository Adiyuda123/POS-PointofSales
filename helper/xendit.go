package helper

import (
	"POS-PointofSales/app/config"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(amount int, externalID string) (string, error) {
	// return "", nil
	xendit.Opt.SecretKey = config.XenditSecretKey
	amountFloat := float64(amount)

	createInvoiceData := invoice.CreateParams{
		// ExternalID:  strconv.Itoa(externalID),
		ExternalID:  externalID,
		Amount:      amountFloat,
		Description: "Invoice for POS",
	}

	inv, err := invoice.Create(&createInvoiceData)
	if err != nil {
		return "", err
	}

	return inv.InvoiceURL, nil
}
