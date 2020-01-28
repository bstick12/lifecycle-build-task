package task

import (
	"encoding/base64"
	"fmt"
)

func EncodeRegistryAuth(registry, username, password string) string {
	delimited := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(delimited))
	return fmt.Sprintf(`{"%s":"Basic %s"}`, registry, encoded)

}
