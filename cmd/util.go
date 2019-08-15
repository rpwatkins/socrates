package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func runValidation(fs afero.Fs) []string {

	missing := []string{}
	incs := checkMaster(fs, master)
	for _, i := range incs {
		if !i.Found {
			missing = append(missing, i.Path)
		}
		missing = append(missing, processValidation(fs, i.Includes)...)

	}
	return missing
}

func processValidation(fs afero.Fs, includes []include) []string {
	missing := []string{}

	for _, i := range includes {
		if !i.Found {
			missing = append(missing, i.Path)
		}
		missing = append(missing, processValidation(fs, i.Includes)...)
	}
	return missing
}

func getAttributes(fs afero.Fs, path string) (map[string]string, error) {
	atts := make(map[string]string)
	// check path
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return nil, err
	}
	if exists {
		content, err := afero.ReadFile(fs, path)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" && line[:1] == ":" {
				parts := strings.Split(line, ":")
				k := strings.TrimSpace(parts[1])
				v := strings.TrimSpace(parts[2])
				atts[k] = v
			}
		}
	}
	return atts, nil
}

func getImagePath(fs afero.Fs, path string) (string, error) {
	m, err := getAttributes(fs, path)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	for k, v := range m {
		if k == "imagesdir" {
			return v, nil
		}
	}
	return "", errors.New("no imagesdir attribute found in master.adoc.")
}

func plural(num int, item string) string {
	if num == 1 {
		return item
	} else {
		return item + "s"
	}
}
