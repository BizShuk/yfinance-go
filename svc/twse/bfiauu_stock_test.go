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

func TestFetchBFIAUUSTOCK_RequiresStockNo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("server should not be called when stockNo is missing")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := FetchBFIAUUSTOCK(context.Background(), c, "20260620", url.Values{})
	if err == nil {
		t.Fatal("expected error when stockNo is missing, got nil")
	}
	if !strings.Contains(err.Error(), "stockNo") {
		t.Fatalf("expected error to mention stockNo, got: %v", err)
	}
}

func TestFetchBFIAUUSTOCK_Decode(t *testing.T) {
	var seenURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"stat":  "OK",
			"title": "單一證券鉅額交易",
			"fields": []string{"證券代號", "成交價", "成交量"},
			"data": [][]string{
				{"2330", "1,000", "500,000", "500,000,000"},
			},
			"date": "20260620",
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	opts := url.Values{}
	opts.Set("stockNo", "2330")
	raw, err := FetchBFIAUUSTOCK(context.Background(), c, "20260620", opts)
	if err != nil {
		t.Fatalf("FetchBFIAUUSTOCK returned error: %v", err)
	}
	resp, ok := raw.(BFIAUResponse)
	if !ok {
		t.Fatalf("expected *BFIAUResponse, got %T", raw)
	}
	if resp.GetStat() != "OK" {
		t.Fatalf("expected stat OK, got %q", resp.GetStat())
	}
	if !strings.Contains(seenURL, "stockNo=2330") {
		t.Errorf("expected URL to contain stockNo=2330, got %q", seenURL)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 data row, got %d", len(resp.Data))
	}
	if resp.Data[0][0] != "2330" {
		t.Errorf("expected stockNo=2330 in row, got %q", resp.Data[0][0])
	}
}

func TestParseBFIAUUSTOCKRow(t *testing.T) {
	row := []string{"2330", "1,000", "500,000", "500,000,000"}
	parsed, err := ParseBFIAUUSTOCKRow(row)
	if err != nil {
		t.Fatalf("ParseBFIAUUSTOCKRow returned error: %v", err)
	}
	if parsed.StockNo != "2330" {
		t.Errorf("StockNo = %q, want %q", parsed.StockNo, "2330")
	}
	if parsed.Price != 1000 {
		t.Errorf("Price = %v, want 1000", parsed.Price)
	}
	if parsed.Volume != 500000 {
		t.Errorf("Volume = %v, want 500000", parsed.Volume)
	}
	if parsed.Amount != 500000000 {
		t.Errorf("Amount = %v, want 500000000", parsed.Amount)
	}
}

func TestParseBFIAUUSTOCKRow_TooShort(t *testing.T) {
	_, err := ParseBFIAUUSTOCKRow([]string{"2330"})
	if err == nil {
		t.Fatal("expected error for short row, got nil")
	}
}
