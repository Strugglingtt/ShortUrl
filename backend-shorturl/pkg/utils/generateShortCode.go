package utils

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/teris-io/shortid"
)

// GenerateShortCode 生成短码
func GenerateShortCode() (string, error) {
	// 1. 生成短码
	shortCode, err := shortid.Generate()
	if err != nil {
		log.Errorf("generate short code error: %v", err)
		return "", err
	}
	return shortCode, nil
}
