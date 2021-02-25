package radiko

import "runtime/debug"

// Version returns a version of this package.
func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	} else {
		return "undefined"
	}
}
