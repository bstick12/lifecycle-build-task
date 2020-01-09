package task

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func ConfigPlatformEnvVars(dir string, envVars map[string]string) error {

	if len(envVars) == 0 {
		return nil
	}

	platformEnvDir := filepath.Join(dir, "env")
	err := os.MkdirAll(platformEnvDir, os.ModePerm)
	if err != nil {
		return err
	}

	for k, v := range envVars {
		err = ioutil.WriteFile(filepath.Join(platformEnvDir, k), []byte(v), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil

}
