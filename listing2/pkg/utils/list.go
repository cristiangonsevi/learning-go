package utils

import (
	"fmt"
	"listing/pkg/logger"
	"os"
)

func ListCurrentDirectory() error {
	files, err := os.ReadDir(".")

	if err != nil {

		return err
	}

	for _, file := range files {
		var icon rune

		if file.IsDir() {
			icon = '📂'
		} else {
			icon = '🧾'
		}

		msg := fmt.Sprintf("%c %s", icon, file.Name())

		logger.Info(msg)
	}

	return nil
}
