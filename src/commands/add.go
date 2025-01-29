package commands

import (
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Add(args []string) {
	if len(args) < 1 {
		theme.Error("You must specify a path of external php.")
		return
	}

	addPath := args[0]

	// Verify that the specified path exists and contains a php executable
	phpPath := filepath.Join(addPath, "php.exe")
	if _, err := os.Stat(phpPath); os.IsNotExist(err) {
		theme.Error(fmt.Sprintf("The file php.exe was not found in the specified path: %s", addPath))
		return
	}

	// Run the "php -v" command to get the PHP version
	cmd := exec.Command(phpPath, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		theme.Error(fmt.Sprintf("failed to execute php -v: %v", err))
	}

	// Parse the version from the output
	// The output usually looks like: PHP 7.4.3 (cli) (built: Mar  4 2020 22:44:12) (ZTS)
	versionPattern := `PHP\s+([0-9]+\.[0-9]+\.[0-9]+)`
	re := regexp.MustCompile(versionPattern)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		theme.Error("failed to parse PHP version from output")
	}

	// Determine if the version is thread-safe (TS) or non-thread-safe (NTS)
	threadSafe := ""
	if strings.Contains(strings.ToLower(phpPath), "nts") {
		threadSafe = " nts"
	}

	// Add to versions.json
	common.AddToVersionsJson(phpPath, matches[1]+threadSafe, "external")

	theme.Success(fmt.Sprintf("Finished add PHP %s", phpPath))
}
