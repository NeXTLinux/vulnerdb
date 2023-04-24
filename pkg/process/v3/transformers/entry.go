package transformers

import (
	"github.com/nextlinux/vulnersdb/pkg/data"
	grypeDB "github.com/nextlinux/grype/grype/db/v3"
)

func NewEntries(vs []grypeDB.Vulnerability, metadata grypeDB.VulnerabilityMetadata) []data.Entry {
	entries := []data.Entry{
		{
			DBSchemaVersion: grypeDB.SchemaVersion,
			Data:            metadata,
		},
	}
	for _, vuln := range vs {
		entries = append(entries, data.Entry{
			DBSchemaVersion: grypeDB.SchemaVersion,
			Data:            vuln,
		})
	}
	return entries
}
