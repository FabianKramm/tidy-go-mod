package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ModuleRegEx retrieves the packages in the go.mod / go.sum file
var ModuleRegEx = regexp.MustCompile("^\\s*([^\\s]+)\\sv.*$")

func main() {
	files := []string{"go.mod", "go.sum"}
	for _, file := range files {
		// Read file
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Skipping file %s, because it couldn't be read: %v", file, err)
			continue
		}

		// Erase lines
		newData := EraseModules(string(data))

		// Write file
		err = ioutil.WriteFile(file, []byte(newData), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully erased all modules that are unused in %s", file)
	}
}

// EraseModules erases all lines that match the go mod regex and do not exist in the vendor folder
func EraseModules(data string) string {
	// Parse lines and find packages that are not in the vendor folder
	newLines := []string{}
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if matches := ModuleRegEx.FindStringSubmatch(line); len(matches) == 2 {
			// Check if module exists
			_, err := os.Stat(filepath.Join("vendor", matches[1]))
			if os.IsNotExist(err) {
				continue
			}
		}

		newLines = append(newLines, line)
	}

	return strings.Join(newLines, "\n")
}
