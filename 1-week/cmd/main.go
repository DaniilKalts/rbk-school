package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/DaniilKalts/rbk-school/1-week/internal/app"
	"github.com/DaniilKalts/rbk-school/1-week/internal/utils"
)

const (
	msgErrorPrefix = "Ошибка:"
	msgUsage       = "Использование: go run ./cmd <входной-файл> <выходной-файл>"
)

func main() {
	args := os.Args[1:]

	paths, err := utils.ParseArgs(args)
	if err != nil {
		if errors.Is(err, utils.ErrUsage) {
			fmt.Fprintf(os.Stderr, "%s ожидалось 2 аргумента, получено %d\n", msgErrorPrefix, len(args))
			fmt.Fprintln(os.Stderr, msgUsage)
			os.Exit(2)
		}

		fail(err, 1)
	}

	inputPath := paths.InputPath
	outputPath := paths.OutputPath

	if err := app.Run(inputPath, outputPath); err != nil {
		fail(err, 1)
	}

	fmt.Println("Входной файл:", inputPath)
	fmt.Println("Выходной файл:", outputPath)
}

func fail(err error, code int) {
	fmt.Fprintln(os.Stderr, msgErrorPrefix, err)
	os.Exit(code)
}
