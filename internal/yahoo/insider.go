package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
)

type InsiderDTO struct {
	Transactions     []InsiderTransaction
	PurchaseActivity NetSharePurchaseActivity
	Roster           []InsiderHolder
}

type InsiderTransaction struct {
	FilerName       string `json:"filerName"`
	TransactionText string `json:"transactionText"`
	Shares          RawInt `json:"shares"`
	Value           RawInt `json:"value"`
	StartDate       RawInt `json:"startDate"`
	OwnershipType   string `json:"ownership"`
}

type NetSharePurchaseActivity struct {
	Period         string `json:"period"`
	BuyInfoShares  RawInt `json:"buyInfoShares"`
	SellInfoShares RawInt `json:"sellInfoShares"`
	NetInfoShares  RawInt `json:"netInfoShares"`
}

type InsiderHolder struct {
	Name            string `json:"name"`
	Relation        string `json:"relation"`
	PositionDirect  RawInt `json:"positionDirect"`
	LatestTransDate RawInt `json:"latestTransDate"`
}

type insiderResult struct {
	QuoteSummary struct {
		Result []struct {
			InsiderTransactions struct {
				Transactions []InsiderTransaction `json:"transactions"`
			} `json:"insiderTransactions"`
			NetSharePurchaseActivity NetSharePurchaseActivity `json:"netSharePurchaseActivity"`
			InsiderHolders           struct {
				Holders []InsiderHolder `json:"holders"`
			} `json:"insiderHolders"`
		} `json:"result"`
	} `json:"quoteSummary"`
}

func DecodeInsider(data []byte) (*InsiderDTO, error) {
	var r insiderResult
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if len(r.QuoteSummary.Result) == 0 {
		return nil, fmt.Errorf("insider: empty result")
	}
	res := r.QuoteSummary.Result[0]
	return &InsiderDTO{
		Transactions:     res.InsiderTransactions.Transactions,
		PurchaseActivity: res.NetSharePurchaseActivity,
		Roster:           res.InsiderHolders.Holders,
	}, nil
}

func (c *Client) FetchInsider(ctx context.Context, symbol string) (*InsiderDTO, error) {
	raw, err := c.FetchQuoteSummary(ctx, symbol,
		[]string{"insiderTransactions", "netSharePurchaseActivity", "insiderHolders"})
	if err != nil {
		return nil, err
	}
	return DecodeInsider(raw)
}
