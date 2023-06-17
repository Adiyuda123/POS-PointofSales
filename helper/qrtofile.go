package helper

import (
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQRCode(url string, customer string) (string, error) {
	fullURL := "https://api.xendit.co/qr_codes/" + url
	fileName := fmt.Sprintf("QR_%s.png", customer)

	err := qrcode.WriteFile(fullURL, qrcode.Medium, 256, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return fileName, nil
}
