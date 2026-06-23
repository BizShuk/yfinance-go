package twse

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AmpyFin/yfinance-go/internal/httpx"
	"github.com/stretchr/testify/require"
)

// newTestHttpxClient returns an *httpx.Client whose Config.BaseURL points
// at the httptest server, so the client's Call method resolves the path
// correctly. It also overrides BaseURL so FetchJSON's internal URL path
// matches the test server.
func newTestHttpxClient(t *testing.T, srv *httptest.Server) *httpx.Client {
	t.Helper()
	cfg := httpx.DefaultConfig()
	cfg.MaxAttempts = 1
	cfg.BaseURL = srv.URL
	old := BaseURL
	BaseURL = srv.URL
	t.Cleanup(func() { BaseURL = old })
	return httpx.NewClient(cfg)
}

func TestFetchJSON_Decodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"stat":"OK","title":"MI_INDEX","fields":["a","b"],"data":[["1","2"],["3","4"]]}`))
	}))
	defer srv.Close()

	c := newTestHttpxClient(t, srv)
	got, err := FetchJSON[TestResponse](context.Background(), c, "/test/endpoint", nil)
	require.NoError(t, err)
	require.Equal(t, "OK", got.Stat)
	require.Equal(t, "MI_INDEX", got.Title)
	require.Len(t, got.Data, 2)
}

func TestFetchJSON_NoDataReturnsErrNoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"stat":"沒有符合條件的資料","fields":[],"data":[]}`))
	}))
	defer srv.Close()

	c := newTestHttpxClient(t, srv)
	_, err := FetchJSON[TestResponse](context.Background(), c, "/test/endpoint", nil)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoData))
}

func TestFetchJSON_StatAtTopLevel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"stat":"OK","data":[]}`))
	}))
	defer srv.Close()

	c := newTestHttpxClient(t, srv)
	got, err := FetchJSON[TestResponse](context.Background(), c, "/test/endpoint", nil)
	require.NoError(t, err)
	require.Equal(t, "OK", got.Stat)
}

type EmbeddedResponse struct {
	Response
	Date string `json:"date"`
}

func (r *EmbeddedResponse) GetStat() string { return r.Response.Stat }

func TestFetchJSON_EmbeddedStructReportsStat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"stat":"OK","date":"20221230","fields":["a"],"data":[["x"]]}`))
	}))
	defer srv.Close()

	c := newTestHttpxClient(t, srv)
	got, err := FetchJSON[EmbeddedResponse](context.Background(), c, "/test/endpoint", nil)
	require.NoError(t, err)
	require.Equal(t, "OK", got.GetStat())
	require.Equal(t, "20221230", got.Date)
}

// TestResponse is a sample struct matching the TWSE JSON envelope.
type TestResponse struct {
	Stat   string     `json:"stat"`
	Title  string     `json:"title"`
	Fields []string   `json:"fields"`
	Data   [][]string `json:"data"`
}

func (r *TestResponse) GetStat() string { return r.Stat }
