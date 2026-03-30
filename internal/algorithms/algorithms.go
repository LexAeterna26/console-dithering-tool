package algorithms

import (
	"image"
	"image/color"
)

func ThresholdFunc(threshold int) func(*image.NRGBA, *image.NRGBA, int, int) {
	t := uint8(threshold)
	changeColor := func(c uint8) uint8 {
		if c > t {
			return 255
		}
		return 0
	}
	return func(source, newImg *image.NRGBA, x, y int) {
		clr := source.At(x, y).(color.NRGBA)
		grayClr := color.GrayModel.Convert(clr).(color.Gray)
		clr.R = changeColor(grayClr.Y)
		clr.G = changeColor(grayClr.Y)
		clr.B = changeColor(grayClr.Y)
		newImg.SetNRGBA(x, y, clr)
	}
}

func BayerFunc(matrixSize int) func(*image.NRGBA, *image.NRGBA, int, int) {
	return func(source, newImg *image.NRGBA, x, y int) {

	}
}
