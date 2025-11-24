package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	// 制造一个 随机ID
	id := uuid.New()
	newUUID := strings.ReplaceAll(id.String(), "-", "")
	return strings.ToLower(newUUID)
}
