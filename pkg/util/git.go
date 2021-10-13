package util

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepository(url, dest string) error {
	options := git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	_, err := git.PlainClone(dest, false, &options)
	return err
}
