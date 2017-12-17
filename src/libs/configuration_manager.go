package libs
import (
	"os"
	"strconv"
)

func GetConfiguredPort() string {
	configuredPort := os.Getenv("TCP_PORT")
	if configuredPort != "" {
		return configuredPort
	} else {
		return "8081"
	}
}

func ReadAgentsUrls() []string {
	ips_num_string := os.Getenv("agent_ips_num")
	if ips_num_string == "" {
		ips_num_string = "0"
	}

	ips_num, _ := strconv.Atoi(ips_num_string)

	agentUrls := make([]string, ips_num)
	for i := 0; i < ips_num; i++ {
		agentUrls[i] = os.Getenv("agent_ip" + strconv.Itoa(i + 1))
	}

	return agentUrls

}
