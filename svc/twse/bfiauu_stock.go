package twse

import (
	"context"
	"fmt"
	"net/url"

	"github.com/AmpyFin/yfinance-go/internal/httpx"
)

// FetchBFIAUUSTOCK is a thin wrapper over FetchBFIAU that requires `stockNo`
// to be supplied via opts. It delegates to FetchBFIAU once the parameter is
// validated.
func FetchBFIAUUSTOCK(ctx context.Context, c *httpx.Client, date string, opts url.Values) (any, error) {
	if opts.Get("stockNo") == "" {
		return nil, fmt.Errorf("twse/BFIAUU_STOCK: stockNo is required")
	}
	return FetchBFIAU(ctx, c, date, opts)
}

// ParseBFIAUUSTOCKRow parses a row from the BFIAUU_STOCK endpoint. The data
// shape is identical to BFIAU, so it delegates to ParseBFIAURow.
func ParseBFIAUUSTOCKRow(row []string) (BFIAURow, error) {
	return ParseBFIAURow(row)
}
