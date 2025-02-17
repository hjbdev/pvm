package commands

import (
	"fmt"
	"hjbdev/pvm/common"
	"hjbdev/pvm/theme"
)

func Remove(args []string) {
	if len(args) < 1 {
		theme.Error("You must specify a path of external php.")
		return
	}

	removePath := args[0]

	// Add to versions.json
	common.RemoveFromVersionJson(removePath)

	theme.Success(fmt.Sprintf("Finished remove PHP %s", removePath))
}
