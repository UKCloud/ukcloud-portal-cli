package main

import "github.com/UKCloud/ukcloud-portal-cli/ukc"

// GitCommit is the git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version is the actual version
const Version = ukc.Version

// VersionPrerelease indicates if dev or not
const VersionPrerelease = ukc.VersionPrerelease
