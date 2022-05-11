package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func Install() {
	// perform get request to https://windows.php.net/downloads/releases/archives/
	resp, err := http.Get("https://windows.php.net/downloads/releases/archives/")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	// log.Println(sb)

	// regex match
	re := regexp.MustCompile(`<A HREF="([a-zA-Z0-9./-]+)">([a-zA-Z0-9./-]+)</A>`)
	matches := re.FindAllStringSubmatch(sb, -1)
	for _, match := range matches {
		// path := match[1]
		name := match[2]

		// check if name starts with "php-devel-pack-"
		if name != "" && len(name) > 15 && name[:15] == "php-devel-pack-" {
			continue
		}
		// check if name starts with "php-debug-pack-"
		if name != "" && len(name) > 15 && name[:15] == "php-debug-pack-" {
			continue
		}
		// check if name starts with "php-test-pack-"
		if name != "" && len(name) > 15 && name[:14] == "php-test-pack-" {
			continue
		}

		// check if name contains "src"
		if name != "" && strings.Contains(name, "src") {
			continue
		}

		// check if name does not end in zip
		if name != "" && !strings.HasSuffix(name, ".zip") {
			continue
		}

		fmt.Println(name)

		threadSafe := true

		// check if name contains "nts" or "NTS"
		if name != "" && (strings.Contains(name, "nts") || strings.Contains(name, "NTS")) {
			threadSafe = false
		}

		// regex match name
		re = regexp.MustCompile(`([0-9]{1,3}).([0-9]{1,3}).([0-9]{1,3})`)
		versionMatches := re.FindAllStringSubmatch(name, -1)

		var major string
		var minor string
		var patch string

		for _, versionMatch := range versionMatches {
			major = versionMatch[1]
			minor = versionMatch[2]
			patch = versionMatch[3]
		}

		fmt.Println(major, minor, patch, threadSafe)
	}
}
