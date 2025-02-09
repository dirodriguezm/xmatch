package allwise_metadata

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	inbox := make(chan indexer.ReaderResult)
	outbox := make(chan indexer.WriterInput[repository.AllwiseMetadata])
	actor := New(inbox, outbox)

	actor.Start()
	rows := make([]repository.InputSchema, 10)
	for i := 0; i < 10; i++ {
		designation := "test"
		w1mpro := 1.0
		w1sigmpro := 1.0
		w2mpro := 2.0
		w2sigmpro := 2.0
		rows[i] = &repository.AllwiseInputSchema{
			Designation: &designation,
			W1mpro:      &w1mpro,
			W1sigmpro:   &w1sigmpro,
			W2mpro:      &w2mpro,
			W2sigmpro:   &w2sigmpro,
		}
	}
	inbox <- indexer.ReaderResult{
		Rows:  rows,
		Error: nil,
	}
	close(inbox)

	for msg := range outbox {
		require.NoError(t, msg.Error)
		require.Len(t, msg.Rows, 10)
		for i := 0; i < 10; i++ {
			require.Equal(t, "test", *msg.Rows[i].Designation)
			require.Equal(t, 1.0, *msg.Rows[i].W1mpro)
			require.Equal(t, 1.0, *msg.Rows[i].W1sigmpro)
			require.Equal(t, 2.0, *msg.Rows[i].W2mpro)
			require.Equal(t, 2.0, *msg.Rows[i].W2sigmpro)
		}
	}
}
