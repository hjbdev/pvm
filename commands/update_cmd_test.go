package commands

import (
	"context"
	"testing"
)

func TestUpdate_NoNewer_DoesNotInstall(t *testing.T) {
	origFetch := fetchLatestRelease
	origInstall := installLatestRunner
	defer func() { fetchLatestRelease = origFetch; installLatestRunner = origInstall }()

	fetchLatestRelease = func(ctx context.Context) (string, error) {
		return "v1.2.3", nil
	}

	installCalled := false
	installLatestRunner = func(auto bool) error {
		installCalled = true
		return nil
	}

	// set current version equal to latest
	prev := version
	version = "v1.2.3"
	defer func() { version = prev }()

	if err := Update([]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if installCalled {
		t.Fatalf("install should not have been called")
	}
}

func TestUpdate_WithYesFlag_Installs(t *testing.T) {
	origFetch := fetchLatestRelease
	origInstall := installLatestRunner
	defer func() { fetchLatestRelease = origFetch; installLatestRunner = origInstall }()

	fetchLatestRelease = func(ctx context.Context) (string, error) {
		return "v1.2.4", nil
	}

	installCalled := false
	installLatestRunner = func(auto bool) error {
		if !auto {
			t.Fatalf("expected auto=true when --yes-update provided")
		}
		installCalled = true
		return nil
	}

	prev := version
	version = "v1.2.3"
	defer func() { version = prev }()

	if err := Update([]string{"--yes-update"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !installCalled {
		t.Fatalf("install should have been called")
	}
}
