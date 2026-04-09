package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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

type TokenKind int

const (
	TokenWord TokenKind = iota
	TokenCommand
	TokenSpace
	TokenPunctuation
)

type Token struct {
	Kind  TokenKind
	Value string
}

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
	n := len(args)
	if n != 2 {
		return "", "", fmt.Errorf("ожидалось 2 аргумента, получено %d: %w", n, errUsage)
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

func tokenize(input string) []Token {
	tokens := make([]Token, 0)

	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)

		switch {
		case unicode.IsSpace(r):
			value := readSpace(input)
			tokens = append(tokens, Token{Kind: TokenSpace, Value: value})
			input = input[len(value):]
		case isPunctuation(r):
			tokens = append(tokens, Token{Kind: TokenPunctuation, Value: input[:size]})
			input = input[size:]
		case r == '(':
			if value, ok := readCommand(input); ok {
				tokens = append(tokens, Token{Kind: TokenCommand, Value: value})
				input = input[len(value):]
				continue
			}

			value := readTextToken(input)
			tokens = append(tokens, Token{Kind: TokenWord, Value: value})
			input = input[len(value):]
		default:
			value := readTextToken(input)
			tokens = append(tokens, Token{Kind: TokenWord, Value: value})
			input = input[len(value):]
		}
	}

	return tokens
}

func readSpace(input string) string {
	end := 0

	for i, r := range input {
		if !unicode.IsSpace(r) {
			break
		}

		end = i + utf8.RuneLen(r)
	}

	return input[:end]
}

func readTextToken(input string) string {
	end := 0

	for i, r := range input {
		if unicode.IsSpace(r) || isPunctuation(r) {
			break
		}

		if r == '(' {
			if _, ok := readCommand(input[i:]); ok {
				break
			}
		}

		end = i + utf8.RuneLen(r)
	}

	if end == 0 {
		_, size := utf8.DecodeRuneInString(input)
		return input[:size]
	}

	return input[:end]
}

func readCommand(input string) (string, bool) {
	end := strings.IndexByte(input, ')')
	if end == -1 {
		return "", false
	}

	raw := input[:end+1]
	if !isValidCommand(raw) {
		return "", false
	}

	return raw, true
}

func isValidCommand(raw string) bool {
	n := len(raw)
	if n < 3 || raw[0] != '(' || raw[n-1] != ')' {
		return false
	}

	inner := raw[1 : n-1]
	if inner == "" || inner != strings.TrimSpace(inner) {
		return false
	}

	switch inner {
	case "up", "low", "cap", "hex", "bin":
		return true
	}

	parts := strings.Split(inner, ",")
	if len(parts) != 2 {
		return false
	}

	command := parts[0]
	if !isCountCommandName(command) {
		return false
	}

	count := strings.TrimSpace(parts[1])
	if count == "" {
		return false
	}

	n, err := strconv.Atoi(count)
	if err != nil || n <= 0 {
		return false
	}

	return true
}

func isCountCommandName(command string) bool {
	switch command {
	case "up", "low", "cap":
		return true
	}

	return false
}

func isPunctuation(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', '\'':
		return true
	default:
		return false
	}
}
