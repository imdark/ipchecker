package libs
import (
	"os"
)
func GetConfiguredPort() string {
	configuredPort := os.Getenv("TCP_PORT")
	if configuredPort != "" {
		return configuredPort
	} else {
		return "8081"
	}
}

