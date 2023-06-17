package helper

import (
	"POS-PointofSales/app/config"
	"fmt"
	"log"

	// "github.com/labstack/gommon/log"

	// "github.com/skip2/go-qrcode"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"github.com/xendit/xendit-go/qrcode"
)

func CreateInvoice(amount int, externalID string) (string, error) {
	xendit.Opt.SecretKey = config.XenditSecretKey
	amountFloat := float64(amount)

	createInvoiceData := invoice.CreateParams{
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

// CreateQRCodeHelper creates a QR code and returns the QR code response.
func CreateQRCodeHelper2(externalID string, callbackURL string, amount int) (*xendit.QRCode, error) {
	xendit.Opt.SecretKey = config.XenditSecretKey
	amountFloat := float64(amount)
	createQRCodeData := qrcode.CreateQRCodeParams{
		ExternalID:  externalID,
		CallbackURL: callbackURL,
		Type:        xendit.DynamicQRCode,
		Amount:      amountFloat,
	}

	qrCodeResponse, err := qrcode.CreateQRCode(&createQRCodeData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("created QR code: %+v\n", qrCodeResponse)
	return qrCodeResponse, nil
}
