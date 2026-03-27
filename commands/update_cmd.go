package commands

import (
	"fmt"
	"os"

	"hjbdev/pvm/theme"
)

// Update is the CLI entry for `pvm update`.
func Update(args []string) error {
	// detect explicit yes flag from args
	yes := false
	for _, a := range args {
		if a == "--yes-update" || a == "-y" || a == "--yes" {
			yes = true
		}
	}

	latest, newer, err := CheckForUpdate(version)
	if err != nil {
		return fmt.Errorf("could not check for updates: %w", err)
	}

	theme.Info(fmt.Sprintf("Current %s", version))
	theme.Info(fmt.Sprintf("Latest %s", latest))

	if !newer {
		theme.Info("pvm is up to date.")
		return nil
	}

	// newer available
	theme.Info("A newer version is available.")

	// Always run installer non-interactively when a newer release is found.
	theme.Info("Running installer...")
	// If user explicitly provided --yes-update/-y, prefer the safe download runner
	// (useful for automated CI and tests). Otherwise, if `PVM_INSTALL_SCRIPT` is
	// set and user didn't explicitly request --yes-update, run that script.
	if yes {
		if err := installLatestRunner(true); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
	} else if os.Getenv("PVM_INSTALL_SCRIPT") != "" {
		if err := InstallLatest(true); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
	} else {
		if err := installLatestRunner(true); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
	}
	theme.Info("Install completed.")
	return nil
}
