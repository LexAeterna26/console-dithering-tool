package processor

type ImgConf interface {
	FuncType() string
}

// Конфигурация для порогового дизеринга
type ThresholdConf struct {
	funcType  string
	threshold int
}

func (conf *ThresholdConf) FuncType() string {
	return conf.funcType
}

func (conf *ThresholdConf) Threshold() int {
	return conf.threshold
}

func NewThresholdConf(th int) *ThresholdConf {
	conf := &ThresholdConf{
		funcType:  "threshold",
		threshold: th,
	}
	return conf
}

// Конфигурация для дизеринга матрицей Байера
type BayerConf struct {
	funcType   string
	matrixSize int
}

func (conf *BayerConf) FuncType() string {
	return conf.funcType
}

func (conf *BayerConf) MatrixSize() int {
	return conf.matrixSize
}

func NewBayerConf(ms int) *BayerConf {
	conf := &BayerConf{
		funcType:   "bayer",
		matrixSize: ms,
	}
	return conf
}
