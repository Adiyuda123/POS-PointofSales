package helper

import (
	"POS-PointofSales/app/config"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	getQRCodeResponse, err := qrcode.GetQRCode(&qrcode.GetQRCodeParams{
		ExternalID: qrCodeResponse.ExternalID,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("created QR code: %+v\n", qrCodeResponse)
	fmt.Printf("retrieved QR code: %+v\n", getQRCodeResponse)

	return getQRCodeResponse, nil
}

func SendXenditRequest(requestBody interface{}) (*Data, error) {
	url := "https://api.xendit.co/qr_codes"

	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return &Data{}, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return &Data{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	secretAPIKey := config.XenditSecretKey
	authValue := base64.StdEncoding.EncodeToString([]byte(secretAPIKey + ":"))

	req.Header.Set("Authorization", "Basic "+authValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &Data{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Data{}, fmt.Errorf("error reading response body: %w", err)
	}

	var response Data
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &Data{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &response, nil
}

type Data struct {
	ID          string `json:"id"`
	ExternalID  string `json:"external_id"`
	ReferenceID string `json:"reference_id"`
	Currency    string `json:"currency"`
	ChannelCode string `json:"channel_code"`
	Amount      int    `json:"amount"`
	ExpiresAt   string `json:"expires_at"`
	BusinessID  string `json:"business_id"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	QRString    string `json:"qr_string"`
	CallbackURL string `json:"callback_url"`
	Type        string `json:"type"`
	Customer    string `json:"customer"`
	ItemID      int    `json:"item_id"`
	UserID      int    `json:"user_id"`
	Status      string `json:"status"`
}
