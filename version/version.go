package version

var (
	URL       string
	GitRef    string
	GitCommit string
	BuildTime string
	Tag       string
	
	Get func() string
)
