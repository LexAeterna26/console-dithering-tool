package extractor

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LexAeterna26/console-dithering-tool/internal/processor"
)

type Path struct {
	source      string
	destination string
}

var (
	source      string
	destination string
	suffix      string
	funcType    string
	threshold   int
	matrixSize  int
)

func init() {
	const (
		shorthand          = " (shorthand)"
		defaultSource      = "."
		usageSource        = "Image source: file or directory"
		defaultDestination = "."
		usageDestination   = "Destination directory"
		defaultSuffix      = "-dither"
		usageSuffix        = "Processed image file suffix, default: \"-dither\""
		defaultFuncType    = "threshold"
		usageFuncType      = "Dithering function: \"threshold\" (default), \"bayer\""
		defaultThreshold   = 128
		usageThreshold     = "Threshold value for threshold dithering function: [0, 255]"
		defaultMatrixSize  = 2
		usageMatrixSize    = "Matrix size for Bayer Matrix dithering function: 2, 4, 8, 16"
	)

	flag.StringVar(&source, "source", defaultSource, usageSource)
	flag.StringVar(&source, "s", defaultSource, usageSource+shorthand)
	flag.StringVar(&destination, "destination", defaultDestination, usageDestination)
	flag.StringVar(&destination, "d", defaultDestination, usageDestination+shorthand)
	flag.StringVar(&suffix, "suffix", defaultSuffix, usageSuffix)
	flag.StringVar(&funcType, "function", defaultFuncType, usageFuncType)
	flag.StringVar(&funcType, "f", defaultFuncType, usageFuncType+shorthand)
	flag.IntVar(&threshold, "threshold", defaultThreshold, usageThreshold)
	flag.IntVar(&matrixSize, "matrix", defaultMatrixSize, usageMatrixSize)
}

func isImage(p string) bool {
	ext := filepath.Ext(p)
	if ext == ".jpeg" || ext == ".jpg" || ext == ".png" || ext == ".gif" {
		return true
	}
	return false
}

func GetData() ([]Path, processor.ImgConf, error) {
	flag.Parse()

	// Проверка пути к обработанным файлам изображений
	destInfo, err := os.Stat(destination)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("Destination directory is not exist. %s", err)
		} else {
			return nil, nil, fmt.Errorf("Destination path error. %s", err)
		}
	}

	if !destInfo.IsDir() {
		return nil, nil, errors.New("Destination path is not a directory")
	}

	destAbs, err := filepath.Abs(destination)
	if err != nil {
		return nil, nil, fmt.Errorf("Destination absolute path error. %s", err)
	}

	// Получение среза путей к исходным файлам изображений
	sourceInfo, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("Source file or directory is not exist. %s", err)
		} else {
			return nil, nil, fmt.Errorf("Source path error. %s", err)
		}
	}

	var paths []Path
	if sourceInfo.IsDir() {
		sourceFiles, err := os.ReadDir(source)
		if err != nil {
			return nil, nil, fmt.Errorf("Source directory reading error. %s", err)
		}

		for _, file := range sourceFiles {
			if !file.IsDir() && isImage(file.Name()) {
				sourceAbs, err := filepath.Abs(filepath.Join(source, file.Name()))
				if err != nil {
					return nil, nil, fmt.Errorf("Source file %s absolute path error. %s", file.Name(), err)
				}
				paths = append(paths, Path{sourceAbs, destAbs})
			}
		}
		if len(paths) == 0 {
			return nil, nil, fmt.Errorf("Source directory does not have image files. %s", err)
		}
	} else {
		if isImage(source) {
			sourceAbs, err := filepath.Abs(source)
			if err != nil {
				return nil, nil, fmt.Errorf("Source file absolute path error. %s", err)
			}
			paths = append(paths, Path{sourceAbs, destAbs})
		} else {
			return nil, nil, fmt.Errorf("Source file is not an image. %s", err)
		}
	}

	// Создание структуры с конфигурацией для обработки изображения
	var conf processor.ImgConf
	switch funcType {
	case "threshold":
		if threshold > 255 || threshold < 0 {
			return nil, nil, errors.New("Wrong threshold value")
		}
		conf = processor.NewThresholdConf(threshold)
	case "bayer":
		if !(matrixSize == 2 || matrixSize == 4 || matrixSize == 8 || matrixSize == 16) {
			return nil, nil, errors.New("Wrong bayer matrix size")
		}
		conf = processor.NewBayerConf(matrixSize)
	default:
		return nil, nil, errors.New("Unknown dithering function")
	}

	return paths, conf, nil
}
