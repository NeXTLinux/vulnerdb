package nvd

import (
	"github.com/nextlinux/vulnersdb/internal"
	"github.com/nextlinux/vulnersdb/pkg/data"
	"github.com/nextlinux/vulnersdb/pkg/process/v5/transformers"
	"github.com/nextlinux/vulnersdb/pkg/provider/unmarshal"
	"github.com/nextlinux/vulnersdb/pkg/provider/unmarshal/nvd"
	grypeDB "github.com/nextlinux/grype/grype/db/v5"
	"github.com/nextlinux/grype/grype/db/v5/namespace"
)

func Transform(vulnerability unmarshal.NVDVulnerability) ([]data.Entry, error) {
	// TODO: stop capturing record source in the vulnerability metadata record (now that feed groups are not real)
	recordSource := "nvdv2:nvdv2:cves"

	grypeNamespace, err := namespace.FromString("nvd:cpe")
	if err != nil {
		return nil, err
	}

	entryNamespace := grypeNamespace.String()

	uniquePkgs := findUniquePkgs(vulnerability.Configurations...)

	// extract all links
	var links []string
	for _, externalRefs := range vulnerability.References {
		// TODO: should we capture other information here?
		if externalRefs.URL != "" {
			links = append(links, externalRefs.URL)
		}
	}

	// duplicate the vulnerabilities based on the set of unique packages the vulnerability is for
	var allVulns []grypeDB.Vulnerability
	for _, p := range uniquePkgs.All() {
		matches := uniquePkgs.Matches(p)
		cpes := internal.NewStringSet()
		for _, m := range matches {
			cpes.Add(grypeNamespace.Resolver().Normalize(m.Criteria))
		}

		// create vulnerability entry
		allVulns = append(allVulns, grypeDB.Vulnerability{
			ID:                vulnerability.ID,
			VersionConstraint: buildConstraints(uniquePkgs.Matches(p)),
			VersionFormat:     "unknown",
			PackageName:       grypeNamespace.Resolver().Normalize(p.Product),
			Namespace:         entryNamespace,
			CPEs:              cpes.ToSlice(),
			Fix: grypeDB.Fix{
				State: grypeDB.UnknownFixState,
			},
		})
	}

	// create vulnerability metadata entry (a single entry keyed off of the vulnerability ID)
	allCVSS := vulnerability.CVSS()
	metadata := grypeDB.VulnerabilityMetadata{
		ID:           vulnerability.ID,
		DataSource:   "https://nvd.nist.gov/vuln/detail/" + vulnerability.ID,
		Namespace:    entryNamespace,
		RecordSource: recordSource,
		Severity:     nvd.CvssSummaries(allCVSS).Sorted().Severity(),
		URLs:         links,
		Description:  vulnerability.Description(),
		Cvss:         getCvss(allCVSS...),
	}

	return transformers.NewEntries(allVulns, metadata), nil
}

func getCvss(cvss ...nvd.CvssSummary) []grypeDB.Cvss {
	var results []grypeDB.Cvss
	for _, c := range cvss {
		results = append(results, grypeDB.Cvss{
			Version: c.Version,
			Vector:  c.Vector,
			Metrics: grypeDB.CvssMetrics{
				BaseScore:           c.BaseScore,
				ExploitabilityScore: c.ExploitabilityScore,
				ImpactScore:         c.ImpactScore,
			},
		})
	}
	return results
}
