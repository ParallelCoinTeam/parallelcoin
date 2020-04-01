package btcjson

// VersionResult models objects included in the version response.  In the actual result, these objects are keyed by the program or API name. NOTE: This is a btcsuite extension ported from github.com/decred/dcrd/dcrjson.
type VersionResult struct {
	VersionString string `json:"versionstring"`
	Major         uint32 `json:"major"`
	Minor         uint32 `json:"minor"`
	Patch         uint32 `json:"patch"`
	Prerelease    string `json:"prerelease"`
	BuildMetadata string `json:"buildmetadata"`
}
