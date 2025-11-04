package transfers

import (
	"github.com/0xSplits/pulsargocode/pkg/transfers"
	"github.com/xh3b4sd/logger"
)

func tesHan() transfers.API {
	return New(Config{
		Log: logger.Fake(),
	})
}
