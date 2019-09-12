package lifecycle

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Auths map[string]AuthEntry `json:"auths,omitempty"`
}

type AuthEntry struct {
	Auth string `json:"auth"`
}

var ErrNoConfigDir = errors.New("could not determine docker config dir")

func ConfigDir() (string, error) {
	if dc := os.Getenv("DOCKER_CONFIG"); dc != "" {
		return dc, nil
	}
	if h := dockerUserHome(); h != "" {
		return filepath.Join(dockerUserHome(), ".docker"), nil
	}
	return "", ErrNoConfigDir
}

func dockerUserHome() string {
	return os.Getenv("HOME")
}

func (c Config) Write() (string, error) {

	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(configDir, "config.json")
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(configFile, data, 0600)
	if err != nil {
		return "", err
	}

	return configFile, nil
}

func WriteConfig(registry, username, password string) (string, error) {

	delimited := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(delimited))

	config := Config{
		Auths: map[string]AuthEntry{
			registry: AuthEntry{
				Auth: encoded,
			},
		},
	}

	return config.Write()

}

func ResetEnv(origEnv []string) {
	for _, pair := range origEnv {
		i := strings.Index(pair[1:], "=") + 1
		if err := os.Setenv(pair[:i], pair[i+1:]); err != nil {
			panic(fmt.Sprintf("Setenv(%q, %q) failed during reset: %v", pair[:i], pair[i+2:], err))
		}
	}
}
