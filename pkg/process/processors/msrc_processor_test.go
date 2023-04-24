package processors

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nextlinux/vulnersdb/pkg/data"
	testUtils "github.com/nextlinux/vulnersdb/pkg/process/tests"
	"github.com/nextlinux/vulnersdb/pkg/provider/unmarshal"
)

func mockMSRCProcessorTransform(vulnerability unmarshal.MSRCVulnerability) ([]data.Entry, error) {
	return []data.Entry{
		{
			DBSchemaVersion: 0,
			Data:            vulnerability,
		},
	}, nil
}

func TestMSRCProcessor_Process(t *testing.T) {
	f, err := os.Open("test-fixtures/msrc.json")
	require.NoError(t, err)
	defer testUtils.CloseFile(f)

	processor := NewMSRCProcessor(mockMSRCProcessorTransform)
	entries, err := processor.Process(f)

	require.NoError(t, err)
	assert.Len(t, entries, 2)
}
