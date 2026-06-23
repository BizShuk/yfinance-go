package twse

import (
	"context"
	"fmt"
	"net/url"

	"github.com/AmpyFin/yfinance-go/internal/httpx"
)

// BFIAUResponse embeds the common Response envelope and adds the
// `date` field that TWSE returns on this endpoint.
type BFIAUResponse struct {
	Response
	Date string `json:"date"`
}

// GetStat returns the embedded stat field.
func (r *BFIAUResponse) GetStat() string { return r.Response.Stat }

// BFIAURow is a typed representation of one BFIAU/BFIAUU_STOCK data row.
// Fields: 證券代號, 成交價, 成交量, 成交金額.
type BFIAURow struct {
	StockNo string  // 證券代號
	Price   float64 // 成交價
	Volume  int64   // 成交量
	Amount  int64   // 成交金額
}

// FetchBFIAU retrieves the daily block-trade data for `date`.
// `opts` may include a `stockNo=...` parameter to filter to a single security.
func FetchBFIAU(ctx context.Context, c *httpx.Client, date string, opts url.Values) (any, error) {
	if date == "" {
		return nil, fmt.Errorf("twse/BFIAU: date is required")
	}
	q := url.Values{}
	q.Set("date", date)
	for k, vs := range opts {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	return FetchJSON[BFIAUResponse](ctx, c, "/block/BFIAUU", q)
}

// ParseBFIAURow converts one raw `data` row into a typed BFIAURow.
func ParseBFIAURow(row []string) (BFIAURow, error) {
	if len(row) < 4 {
		return BFIAURow{}, fmt.Errorf("BFIAUU: row too short: %d cols", len(row))
	}
	return BFIAURow{
		StockNo: row[0],
		Price:   ParseFloat(row[1]),
		Volume:  ParseInt(row[2]),
		Amount:  ParseInt(row[3]),
	}, nil
}
