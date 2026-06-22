package yahoo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractMetadata(t *testing.T) {
	raw := []byte(`{"chart":{"result":[{
	  "meta":{"symbol":"AAPL","currency":"USD","exchangeName":"NMS",
	    "instrumentType":"EQUITY","timezone":"EST","gmtoffset":-18000,
	    "firstTradeDate":345479400,"regularMarketPrice":150.0}
	}],"error":null}}`)

	m, err := ExtractMetadata(raw)
	require.NoError(t, err)
	require.Equal(t, "AAPL", m.Symbol)
	require.Equal(t, "USD", m.Currency)
	require.Equal(t, "NMS", m.ExchangeName)
}

func TestExtractMetadata_EmptyResult(t *testing.T) {
	_, err := ExtractMetadata([]byte(`{"chart":{"result":[],"error":null}}`))
	require.Error(t, err)
}
