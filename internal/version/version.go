package version

import "runtime/debug"

// Version is filled at build time with
// -ldflags "-X github.com/wishmatic/genwish/internal/version.Version=v0.4.0".
var Version string

// String reports the build-time version, falling back to the version Go records for a module
// installed with go install module@version, and then to "dev" for a plain source build.
func String() string {
	if Version != "" {
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	return "dev"
}
