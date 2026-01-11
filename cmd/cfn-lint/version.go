package main

import "runtime/debug"

// getVersion returns the application version.
// Priority:
// 1. If version is set via ldflags (not "dev"), use that
// 2. Try to get version from Go module info (for go install)
// 3. Fall back to "dev"
func getVersion() string {
	// If version was set via ldflags, use it
	if version != "dev" {
		return version
	}

	// Try to get version from build info (works with go install pkg@version)
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}
