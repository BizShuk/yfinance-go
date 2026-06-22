package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
)

type ChartMetadata struct {
	Symbol             string  `json:"symbol"`
	Currency           string  `json:"currency"`
	ExchangeName       string  `json:"exchangeName"`
	InstrumentType     string  `json:"instrumentType"`
	Timezone           string  `json:"timezone"`
	GmtOffset          int     `json:"gmtoffset"`
	FirstTradeDate     int64   `json:"firstTradeDate"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

type metaResult struct {
	Chart struct {
		Result []struct {
			Meta ChartMetadata `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func ExtractMetadata(data []byte) (*ChartMetadata, error) {
	var r metaResult
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if len(r.Chart.Result) == 0 {
		return nil, fmt.Errorf("metadata: empty result")
	}
	m := r.Chart.Result[0].Meta
	return &m, nil
}

func (c *Client) FetchMetadata(ctx context.Context, symbol string) (*ChartMetadata, error) {
	raw, err := c.fetchChartRaw(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return ExtractMetadata(raw)
}
