package utils

import (
	"fmt"
	"os"
)

func ReadInputFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %q: %w", filePath, err)
	}

	return data, nil
}

func WriteOutputFile(filePath string, data []byte) error {
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("не удалось записать файл %q: %w", filePath, err)
	}

	return nil
}
