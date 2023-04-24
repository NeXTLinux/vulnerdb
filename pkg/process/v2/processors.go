package v2

import (
	"github.com/nextlinux/vulnersdb/pkg/data"
	"github.com/nextlinux/vulnersdb/pkg/process/processors"
	"github.com/nextlinux/vulnersdb/pkg/process/v2/transformers/github"
	"github.com/nextlinux/vulnersdb/pkg/process/v2/transformers/nvd"
	"github.com/nextlinux/vulnersdb/pkg/process/v2/transformers/os"
)

func Processors() []data.Processor {
	return []data.Processor{
		processors.NewGitHubProcessor(github.Transform),
		processors.NewNVDProcessor(nvd.Transform),
		processors.NewOSProcessor(os.Transform),
	}
}
