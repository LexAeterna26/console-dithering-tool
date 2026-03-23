package main

import (
	"fmt"
	"os"

	"github.com/LexAeterna26/console-dithering-tool/internal/extractor"
	"github.com/LexAeterna26/console-dithering-tool/internal/logger"
)

func main() {
	logger := logger.New()

	imgPaths, imgConfig, err := extractor.GetData()
	if err != nil {
		logger.Error("Failed to retrieve data:", "error", err)
		os.Exit(1)
	}

	fmt.Println(imgPaths)
	fmt.Println(imgConfig)
}
