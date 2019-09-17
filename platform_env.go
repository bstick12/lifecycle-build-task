package lifecycle

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func ConfigPlatformEnvVars(dir string, envVars []EnvVariable) error {

	platformEnvDir := filepath.Join(dir, "env")
	err := os.MkdirAll(platformEnvDir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, envVar := range envVars {
		err = ioutil.WriteFile(filepath.Join(platformEnvDir, envVar.Name), []byte(envVar.Value), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil

}
