package util

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepository(url, path string) error {
	options := git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	_, err := git.PlainClone(path, false, &options)
	return err
}
