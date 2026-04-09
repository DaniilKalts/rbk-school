package app

import (
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc"
	"github.com/DaniilKalts/rbk-school/1-week/internal/utils"
)

func Run(inputPath, outputPath string) error {
	input, err := utils.ReadInputFile(inputPath)
	if err != nil {
		return err
	}

	result := textproc.Process(input)

	return utils.WriteOutputFile(outputPath, result)
}
