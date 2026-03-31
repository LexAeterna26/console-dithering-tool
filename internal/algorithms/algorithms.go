package algorithms

import (
	"image"
	"image/color"
)

func ThresholdFunc(threshold int) func(*image.NRGBA, *image.NRGBA, int, int) {
	t := uint8(threshold)
	return func(source, newImg *image.NRGBA, x, y int) {
		clr := source.At(x, y).(color.NRGBA)
		grayClr := color.GrayModel.Convert(clr).(color.Gray)
		var newGrayClr uint8
		if grayClr.Y > t {
			newGrayClr = 255
		} else {
			newGrayClr = 0
		}
		clr.R = newGrayClr
		clr.G = newGrayClr
		clr.B = newGrayClr
		newImg.SetNRGBA(x, y, clr)
	}
}

func BayerMatrix(n int) [][]int {
	if n == 1 {
		return [][]int{{0}}
	}

	halfN := n / 2
	matrixQuarter := BayerMatrix(halfN)
	for i := range matrixQuarter {
		for j := range matrixQuarter[i] {
			matrixQuarter[i][j] *= 4
		}
	}

	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}
	points := [][]int{
		{0, 0},
		{halfN, halfN},
		{0, halfN},
		{halfN, 0},
	}

	val := 0
	for _, point := range points {
		row, col := point[0], point[1]
		for i := row; i < row+halfN; i++ {
			for j := col; j < col+halfN; j++ {
				matrix[i][j] = matrixQuarter[i%halfN][j%halfN] + val
			}
		}
		val++
	}

	return matrix
}

func BayerFunc(matrixSize int) func(*image.NRGBA, *image.NRGBA, int, int) {
	matrix := BayerMatrix(matrixSize)
	normMatrix := make([][]uint8, len(matrix))
	for i := range matrixSize {
		normMatrix[i] = make([]uint8, len(matrix))
		for j := range matrixSize {
			normMatrix[i][j] = uint8(255 * matrix[i][j] / (matrixSize * matrixSize))
		}
	}

	return func(source, newImg *image.NRGBA, x, y int) {
		clr := source.At(x, y).(color.NRGBA)
		grayClr := color.GrayModel.Convert(clr).(color.Gray)
		xMatrix := x % matrixSize
		yMatrix := y % matrixSize
		valMatrix := normMatrix[xMatrix][yMatrix]
		var newGrayClr uint8
		if grayClr.Y > valMatrix {
			newGrayClr = 255
		} else {
			newGrayClr = 0
		}
		clr.R = newGrayClr
		clr.G = newGrayClr
		clr.B = newGrayClr
		newImg.SetNRGBA(x, y, clr)
	}
}
