package cmd

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func CopyFolder(src string, dst string, fs afero.Fs) error {

	// check for images folder
	exists, err := afero.DirExists(fs, dst)
	if err != nil {
		return err
	}
	// delete it and its contents if images folder exists
	if exists {
		if err := fs.RemoveAll(dst); err != nil {
			return err
		}
	}
	// create images folder in dst
	if err := fs.MkdirAll(dst, 0755); err != nil {
		return err
	}
	// read the images folder
	entries, err := afero.ReadDir(fs, src)
	if err != nil {
		return err
	}
	// copy each entry recursively
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyFolder(srcPath, dstPath, fs)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			// open file
			file, err := afero.ReadFile(fs, srcPath)
			if err != nil {
				return err
			}
			// save file
			if err := afero.WriteFile(fs, dstPath, file, 0644); err != nil {
				return err
			}
			if Verbose {
				log.Infof("%s copied.", dstPath)
			}
		}
	}

	return nil
}
