package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/DaniilKalts/rbk-school/1-week/internal/utils"
)

const (
	msgErrorPrefix = "Ошибка:"
	msgUsage       = "Использование: go run ./cmd <входной-файл> <выходной-файл>"
)

func main() {
	args := os.Args[1:]

	inputPath, outputPath, err := utils.ParseArgs(args)
	if err != nil {
		if errors.Is(err, utils.ErrUsage) {
			fmt.Fprintf(os.Stderr, "%s ожидалось 2 аргумента, получено %d\n", msgErrorPrefix, len(args))
			fmt.Fprintln(os.Stderr, msgUsage)
			os.Exit(2)
		}

		fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
		os.Exit(1)
	}

	inputData, err := utils.ReadInputFile(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
		os.Exit(1)
	}

	if err := utils.WriteOutputFile(outputPath, inputData); err != nil {
		fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
		os.Exit(1)
	}

	fmt.Println("Входной файл:", inputPath)
	fmt.Println("Выходной файл:", outputPath)
}
