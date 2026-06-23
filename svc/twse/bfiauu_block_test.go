package twse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFetchBlockBFIAUU_Decode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/block/BFIAUU") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("date") != "20260620" {
			t.Errorf("missing/wrong date param: %s", r.URL.RawQuery)
		}
		// Verify that stockNo=2330 from opts propagates to URL.
		if r.URL.Query().Get("stockNo") != "2330" {
			t.Errorf("expected stockNo=2330 in URL, got: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"stat":  "OK",
			"title": "鉅額交易",
			"fields": []string{"序號", "證券代號", "證券名稱", "買進證券商", "賣出證券商", "成交數量", "成交金額", "成交價格", "成交時間", "買進成交價"},
			"data": [][]string{
				{"1", "2330", "台積電", "元大", "凱基", "100,000", "60,000,000", "600.00", "13:30:00", "599.50"},
			},
			"date":    "20260620",
			"stockNo": "2330",
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	opts := url.Values{}
	opts.Set("stockNo", "2330")
	raw, err := FetchBlockBFIAUU(context.Background(), c, "20260620", opts)
	if err != nil {
		t.Fatalf("FetchBlockBFIAUU returned error: %v", err)
	}
	resp, ok := raw.(BlockBFIAUUResponse)
	if !ok {
		t.Fatalf("expected *BlockBFIAUUResponse, got %T", raw)
	}
	if resp.GetStat() != "OK" {
		t.Fatalf("expected stat OK, got %q", resp.GetStat())
	}
	if resp.Date != "20260620" {
		t.Errorf("Date = %q, want 20260620", resp.Date)
	}
	if resp.StockNo != "2330" {
		t.Errorf("StockNo = %q, want 2330", resp.StockNo)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 data row, got %d", len(resp.Data))
	}
	row, err := ParseBlockBFIAUURow(resp.Data[0])
	if err != nil {
		t.Fatalf("ParseBlockBFIAUURow returned error: %v", err)
	}
	if row.StockCode != "2330" {
		t.Errorf("StockCode = %q", row.StockCode)
	}
	if row.StockName != "台積電" {
		t.Errorf("StockName = %q", row.StockName)
	}
	if row.BuyBroker != "元大" {
		t.Errorf("BuyBroker = %q", row.BuyBroker)
	}
	if row.SellBroker != "凱基" {
		t.Errorf("SellBroker = %q", row.SellBroker)
	}
	if row.TradeVolume != 100000 {
		t.Errorf("TradeVolume = %v, want 100000", row.TradeVolume)
	}
	if row.TradeAmount != 60000000 {
		t.Errorf("TradeAmount = %v, want 60000000", row.TradeAmount)
	}
	if row.TradePrice != 600.00 {
		t.Errorf("TradePrice = %v, want 600.00", row.TradePrice)
	}
	if row.TradeTime != "13:30:00" {
		t.Errorf("TradeTime = %q", row.TradeTime)
	}
	if row.BuyTradePrice != 599.50 {
		t.Errorf("BuyTradePrice = %v, want 599.50", row.BuyTradePrice)
	}
}

func TestFetchBlockBFIAUU_NoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"stat":"很抱歉，沒有符合條件的資料!","fields":[],"data":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := FetchBlockBFIAUU(context.Background(), c, "20260620", url.Values{})
	if err == nil {
		t.Fatal("expected error for no-data response, got nil")
	}
	if !strings.Contains(err.Error(), "no data") {
		t.Fatalf("expected ErrNoData, got: %v", err)
	}
}