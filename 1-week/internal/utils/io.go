package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadInputFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл %q: %w", filePath, err)
	}
	defer file.Close()

	var builder strings.Builder
	if info, statErr := file.Stat(); statErr == nil {
		builder.Grow(int(info.Size()))
	}

	reader := bufio.NewReaderSize(file, 64*1024)
	if _, err := io.Copy(&builder, reader); err != nil {
		return "", fmt.Errorf("не удалось прочитать файл %q: %w", filePath, err)
	}

	return builder.String(), nil
}

func WriteOutputFile(filePath string, data string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %q для записи: %w", filePath, err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 64*1024)
	if _, err := writer.WriteString(data); err != nil {
		return fmt.Errorf("не удалось записать файл %q: %w", filePath, err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("не удалось сохранить файл %q: %w", filePath, err)
	}

	return nil
}
