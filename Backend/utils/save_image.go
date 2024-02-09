package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

func SaveBase64Image(imageData string) (string, error) {
	parts := strings.Split(imageData, ";base64,")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 data")
	}
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	imagePath := filepath.Join("./media/avatars", filename)

	if err := ioutil.WriteFile(imagePath, decoded, 0644); err != nil {
		return "", err
	}

	return imagePath, nil
}
