package utils

import "os"

func FindRootModulePath(maxDepth int) (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dirs, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, dir := range dirs {
		if dir.Name() == "go.mod" {
			return currDir, nil
		}
	}

	os.Chdir("..")
	return FindRootModulePath(maxDepth - 1)
}
