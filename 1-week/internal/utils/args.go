package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrUsage = errors.New("неверное использование команды")

	errInputPathEmpty  = errors.New("путь к входному файлу пустой")
	errOutputPathEmpty = errors.New("путь к выходному файлу пустой")
	errSamePaths       = errors.New("входной и выходной файлы должны быть разными")
)

func ParseArgs(args []string) (string, string, error) {
	n := len(args)
	if n != 2 {
		return "", "", fmt.Errorf("ожидалось 2 аргумента, получено %d: %w", n, ErrUsage)
	}

	inputPath := args[0]
	outputPath := args[1]

	if inputPath == "" {
		return "", "", errInputPathEmpty
	}

	if outputPath == "" {
		return "", "", errOutputPathEmpty
	}

	if inputPath == outputPath {
		return "", "", errSamePaths
	}

	if err := validateInputPath(inputPath); err != nil {
		return "", "", err
	}

	if err := validateOutputPath(outputPath); err != nil {
		return "", "", err
	}

	return inputPath, outputPath, nil
}

func validateInputPath(inputPath string) error {
	inputFileInfo, err := os.Stat(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("входной файл не найден: %s", inputPath)
		}

		return fmt.Errorf("не удалось проверить входной файл: %w", err)
	}

	if inputFileInfo.IsDir() {
		return fmt.Errorf("входной путь указывает на директорию, а не на файл: %s", inputPath)
	}

	return nil
}

func validateOutputPath(outputPath string) error {
	outputFileInfo, err := os.Stat(outputPath)
	if err == nil {
		if outputFileInfo.IsDir() {
			return fmt.Errorf("выходной путь указывает на директорию, а не на файл: %s", outputPath)
		}

		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("не удалось проверить выходной файл: %w", err)
	}

	outputDir := filepath.Dir(outputPath)
	outputDirInfo, err := os.Stat(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("директория выходного файла не найдена: %s", outputDir)
		}

		return fmt.Errorf("не удалось проверить директорию выходного файла: %w", err)
	}

	if !outputDirInfo.IsDir() {
		return fmt.Errorf("путь к директории выходного файла не является директорией: %s", outputDir)
	}

	return nil
}
