package transfers

import (
	"context"

	"github.com/0xSplits/pulsargocode/pkg/transfers"
)

func (h *Handler) Search(ctx context.Context, req *transfers.SearchI) (*transfers.SearchO, error) {
	return &transfers.SearchO{}, nil
}
