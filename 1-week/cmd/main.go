package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	msgErrorPrefix = "Ошибка:"
	msgUsage       = "Использование: go run ./cmd <входной-файл> <выходной-файл>"
)

var (
	errUsage           = errors.New("неверное использование команды")
	errInputPathEmpty  = errors.New("путь к входному файлу пустой")
	errOutputPathEmpty = errors.New("путь к выходному файлу пустой")
	errSamePaths       = errors.New("входной и выходной файлы должны быть разными")
)

func main() {
	inputPath, outputPath, err := parseArgs(os.Args[1:])
	if err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprintf(os.Stderr, "%s ожидалось 2 аргумента, получено %d\n", msgErrorPrefix, len(os.Args[1:]))
			fmt.Fprintln(os.Stderr, msgUsage)
			os.Exit(2)
		}

		fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
		os.Exit(1)
	}

	inputData, err := readInputFile(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
		os.Exit(1)
	}

	if err := writeOutputFile(inputData, outputPath); err != nil {
		fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
		os.Exit(1)
	}

	fmt.Println("Входной файл:", inputPath)
	fmt.Println("Выходной файл:", outputPath)
}

func parseArgs(args []string) (string, string, error) {
	if len(args) != 2 {
		return "", "", fmt.Errorf("ожидалось 2 аргумента, получено %d: %w", len(args), errUsage)
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

func readInputFile(inputPath string) ([]byte, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать входной файл: %w", err)
	}

	return data, nil
}

func writeOutputFile(inputData []byte, outputPath string) error {
	if err := os.WriteFile(outputPath, inputData, 0o644); err != nil {
		return fmt.Errorf("не удалось записать выходной файл: %w", err)
	}

	return nil
}
