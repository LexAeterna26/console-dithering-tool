package processor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/LexAeterna26/console-dithering-tool/internal/algorithms"
	"github.com/LexAeterna26/console-dithering-tool/internal/validator"
)

type Path struct {
	Source      string
	Destination string
}

func GetFunction(config ImgConf) (func(*image.NRGBA, *image.NRGBA, int, int), error) {
	funcType := config.FuncType()
	var f func(*image.NRGBA, *image.NRGBA, int, int)
	switch funcType {
	case "threshold":
		thresholdConfig := config.(*ThresholdConf)
		f = algorithms.ThresholdFunc(thresholdConfig.threshold)
	case "bayer":
		BayerConf := config.(*BayerConf)
		f = algorithms.BayerFunc(BayerConf.matrixSize)
	default:
		return nil, errors.New("Wrong dithering function configuration.")
	}
	return f, nil
}

func MakeImage(sourceImg *image.NRGBA, f func(*image.NRGBA, *image.NRGBA, int, int)) *image.NRGBA {
	bounds := sourceImg.Bounds()
	newImg := image.NewNRGBA(bounds)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			f(sourceImg, newImg, x, y)
		}
	}
	return newImg
}

func EncodeImage(imageFormat string, newFile *os.File, newImage *image.NRGBA) error {
	switch imageFormat {
	case "png":
		if err := png.Encode(newFile, newImage); err != nil {
			return err
		}
	case "gif":
		if err := gif.Encode(newFile, newImage, nil); err != nil {
			return err
		}
	case "jpg", "jpeg":
		if err := jpeg.Encode(newFile, newImage, nil); err != nil {
			return err
		}
	default:
		return errors.New("Wrong image format.")
	}

	return nil
}

func ProcessImage(imgPath Path, f func(*image.NRGBA, *image.NRGBA, int, int)) error {
	// Открытие исходного файла
	sourceFile, err := os.ReadFile(imgPath.Source)
	if err != nil {
		return fmt.Errorf("Failed to open source image. Path: %s. %s", imgPath.Source, err)
	}

	// Получение конфигурации изображения и проверка размеров
	sourceConfReader := bytes.NewReader(sourceFile)
	sourceConf, _, err := image.DecodeConfig(sourceConfReader)
	if err != nil {
		return fmt.Errorf("Failed to read source image configuration. Path: %s. %s", imgPath.Source, err)
	}

	if !validator.ValidateSize(sourceConf.Width, sourceConf.Height) {
		return fmt.Errorf("Image size is too big. Path: %s. %s", imgPath.Source, err)
	}

	// Получение исходного изображения
	sourceImgReader := bytes.NewReader(sourceFile)
	sourceImageRaw, imageFormat, err := image.Decode(sourceImgReader)
	if err != nil {
		return fmt.Errorf("Failed to decode source image. Path: %s. %s", imgPath.Source, err)
	}

	sourceImage, ok := sourceImageRaw.(*image.NRGBA)
	if !ok {
		bounds := sourceImageRaw.Bounds()
		sourceImage = image.NewNRGBA(bounds)
		draw.Draw(sourceImage, bounds, sourceImageRaw, bounds.Min, draw.Src)
	}

	// Обработка изображения и создание файла
	newImage := MakeImage(sourceImage, f)

	newFile, err := os.Create(imgPath.Destination)
	if err != nil {
		return fmt.Errorf("Failed to create new image. Path: %s. %s", imgPath.Destination, err)
	}
	defer newFile.Close()

	if err := EncodeImage(imageFormat, newFile, newImage); err != nil {
		return fmt.Errorf("Failed to encode new image. Path: %s. %s", imgPath.Destination, err)
	}

	return nil
}
