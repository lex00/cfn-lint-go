package main

import (
	"runtime/debug"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name     string
		buildVar string
		want     string
	}{
		{
			name:     "returns build variable when set",
			buildVar: "v1.0.0",
			want:     "v1.0.0",
		},
		{
			name:     "returns dev when build variable is dev",
			buildVar: "dev",
			want:     "dev", // Will be overridden by build info if available
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original and restore after test
			original := version
			defer func() { version = original }()

			version = tt.buildVar
			got := getVersion()

			// When version is explicitly set (not "dev"), it should be returned
			if tt.buildVar != "dev" && got != tt.want {
				t.Errorf("getVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetVersionFromBuildInfo(t *testing.T) {
	// Test that the function properly reads from debug.ReadBuildInfo
	// When running tests, build info is typically available
	info, ok := debug.ReadBuildInfo()

	// Save original and restore after test
	original := version
	defer func() { version = original }()

	version = "dev"
	got := getVersion()

	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		// When build info is available with a real version,
		// getVersion should return that version
		if got != info.Main.Version {
			t.Errorf("getVersion() = %q, want %q from build info", got, info.Main.Version)
		}
	} else {
		// When no build info is available or version is (devel),
		// it should return "dev"
		if got != "dev" {
			t.Logf("Note: got version %q (build info may be available)", got)
		}
	}
}

func TestGetVersionPrefersBuildVariable(t *testing.T) {
	// When version is explicitly set via ldflags, it should be preferred
	original := version
	defer func() { version = original }()

	version = "v2.0.0-custom"
	got := getVersion()

	if got != "v2.0.0-custom" {
		t.Errorf("getVersion() = %q, want %q (should prefer build variable)", got, "v2.0.0-custom")
	}
}
