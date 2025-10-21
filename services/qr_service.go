package services

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code and returns it as a base64 encoded string
func GenerateQRCode(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
