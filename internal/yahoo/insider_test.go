package yahoo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeInsider_ParsesTransactions(t *testing.T) {
	raw := []byte(`{"quoteSummary":{"result":[{
	  "insiderTransactions":{"transactions":[
	    {"filerName":"DOE JOHN","transactionText":"Sale at price 150.00",
	     "shares":{"raw":1000},"value":{"raw":150000},
	     "startDate":{"raw":1700000000}}]},
	  "netSharePurchaseActivity":{"period":"6m","buyInfoShares":{"raw":5000},
	     "sellInfoShares":{"raw":2000},"netInfoShares":{"raw":3000}},
	  "insiderHolders":{"holders":[
	    {"name":"DOE JANE","relation":"Director","positionDirect":{"raw":20000}}]}
	}],"error":null}}`)

	d, err := DecodeInsider(raw)
	require.NoError(t, err)
	require.Len(t, d.Transactions, 1)
	require.Equal(t, "DOE JOHN", d.Transactions[0].FilerName)
	require.NotNil(t, d.PurchaseActivity.NetInfoShares.Raw)
	require.Equal(t, int64(3000), *d.PurchaseActivity.NetInfoShares.Raw)
	require.Len(t, d.Roster, 1)
	require.Equal(t, "DOE JANE", d.Roster[0].Name)
}

func TestDecodeInsider_EmptyResult(t *testing.T) {
	_, err := DecodeInsider([]byte(`{"quoteSummary":{"result":[],"error":null}}`))
	require.Error(t, err)
}
