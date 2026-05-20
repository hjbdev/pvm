package commands

import (
	"context"
	"testing"
)

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.4", "1.2.3", 1},
		{"1.2.3", "1.2.4", -1},
		{"1.10.0", "1.9.9", 1},
		{"1.2", "1.2.0", 0},
		{"1.2.0", "1.2.1", -1},
	}
	for _, tc := range tests {
		if got := compareSemver(tc.a, tc.b); got != tc.want {
			t.Fatalf("compareSemver(%s,%s) = %d; want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestCheckForUpdate_DevIsNewer(t *testing.T) {
	orig := fetchLatestRelease
	fetchLatestRelease = func(ctx context.Context) (string, error) {
		return "v1.2.3", nil
	}
	defer func() { fetchLatestRelease = orig }()

	latest, newer, err := CheckForUpdate("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "v1.2.3" || !newer {
		t.Fatalf("unexpected result latest=%s newer=%v", latest, newer)
	}
}

func TestCheckForUpdate_CurrentEqualsLatest(t *testing.T) {
	orig := fetchLatestRelease
	fetchLatestRelease = func(ctx context.Context) (string, error) {
		return "v1.2.3", nil
	}
	defer func() { fetchLatestRelease = orig }()

	latest, newer, err := CheckForUpdate("v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "v1.2.3" || newer {
		t.Fatalf("unexpected result latest=%s newer=%v", latest, newer)
	}
}

func TestCheckForUpdate_CurrentOlder(t *testing.T) {
	orig := fetchLatestRelease
	fetchLatestRelease = func(ctx context.Context) (string, error) {
		return "v1.2.4", nil
	}
	defer func() { fetchLatestRelease = orig }()

	latest, newer, err := CheckForUpdate("v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "v1.2.4" || !newer {
		t.Fatalf("unexpected result latest=%s newer=%v", latest, newer)
	}
}
