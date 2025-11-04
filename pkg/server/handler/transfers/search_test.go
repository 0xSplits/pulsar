package transfers

import (
	"context"
	"testing"

	"github.com/0xSplits/pulsargocode/pkg/transfers"
	fuzz "github.com/google/gofuzz"
)

func Test_Server_Handler_Transfers_Search_Fuzz(t *testing.T) {
	var han transfers.API
	{
		han = tesHan()
	}

	var fuz *fuzz.Fuzzer
	{
		fuz = fuzz.New()
	}

	for range 1000 {
		var inp *transfers.SearchI
		{
			inp = &transfers.SearchI{}
		}

		{
			fuz.Fuzz(inp)
		}

		{
			_, _ = han.Search(context.Background(), inp)
		}
	}
}
