package ukc

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

// Version is main version number that is being run at the moment.
const Version = "0.0.8-dev"

// VersionPrerelease marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
const VersionPrerelease = ""

// SemVersion is an instance of version.Version. This has the secondary
// benefit of verifying during tests and init time that our version is a
// proper semantic version, which should always be the case.
var SemVersion = version.Must(version.NewVersion(Version))

// VersionHeader is the header name used to send the current terraform version
// in http requests.
const VersionHeader = "UKC-Version"

// VersionString is the utlimate output
func VersionString() string {
	if VersionPrerelease != "" {
		return fmt.Sprintf("%s-%s", Version, VersionPrerelease)
	}
	return Version
}
