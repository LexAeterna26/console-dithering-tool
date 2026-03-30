package main

import (
	"os"

	"github.com/LexAeterna26/console-dithering-tool/internal/extractor"
	"github.com/LexAeterna26/console-dithering-tool/internal/logger"
	"github.com/LexAeterna26/console-dithering-tool/internal/processor"
)

func main() {
	logger := logger.New()

	// Получение путей к файлам и конфигурации
	imgPaths, imgConfig, err := extractor.GetData()
	if err != nil {
		logger.Error("Failed to retrieve data:", "error", err)
		os.Exit(1)
	}

	// Выбор функции для дизеринга
	f, err := processor.GetFunction(imgConfig)
	if err != nil {
		logger.Error("Failed to get dithering function:", "error", err)
		os.Exit(1)
	}

	// Обработка изображений
	for _, imgPath := range imgPaths {
		err := processor.ProcessImage(imgPath, f)
		if err != nil {
			logger.Error("Failed to process image.", "error", err)
		} else {
			logger.Info("Image processed successfully.", "source", imgPath.Source)
		}
	}
}
