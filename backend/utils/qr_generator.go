package utils

import (
    "encoding/base64"
    "encoding/json"
    "github.com/skip2/go-qrcode"
)

type QRData struct {
    DocumentCode string `json:"document_code"`
    Type         string `json:"type"`
    UserID       uint   `json:"user_id"`
    TransactionID *uint  `json:"transaction_id"`
    Timestamp    int64  `json:"timestamp"`
}

func GenerateQRCode(data string) (string, error) {
    var png []byte
    png, err := qrcode.Encode(data, qrcode.Medium, 256)
    if err != nil {
        return "", err
    }
    
    base64Str := base64.StdEncoding.EncodeToString(png)
    return "data:image/png;base64," + base64Str, nil
}

func GenerateQRCodeFromMap(data map[string]interface{}) (string, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    
    return GenerateQRCode(string(jsonData))
}

func ParseQRCode(qrData string) (*QRData, error) {
    var data QRData
    err := json.Unmarshal([]byte(qrData), &data)
    if err != nil {
        return nil, err
    }
    return &data, nil
}