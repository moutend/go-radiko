package testutil

import "flag"

// IsVerbose returns true when verbose flag `-v` is specified.
func IsVerbose() bool {
	f := flag.Lookup(`test.v`)

	if f == nil {
		return false
	}

	return f.Value.String() == "true"
}
