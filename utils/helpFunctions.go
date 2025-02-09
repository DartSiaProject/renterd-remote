package utils

import (
	"encoding/json"
	"net"
)

func HttpHeaderMapToString(header map[string][]string) string {
	mergedData := "{"
	for names, values := range header {
		// Loop over all values for the name.
		subString := "["
		for i := 0; i < len(values); i++ {
			subString = subString + "\"" + values[i] + "\"]"
		}
		//justString := strings.Join(values, "\", ")
		mergedData = mergedData + "\"" + names + "\":" + subString + ","
	}
	mergedData = mergedData[:len(mergedData)-1] + "}"
	return mergedData
}

func StringToHttpHeaderMap(header string) map[string][]string {
	var jsonMap map[string][]string
	json.Unmarshal([]byte(header), &jsonMap)
	return jsonMap
}

type Result struct {
	ContentType string `json:"content-type"`
}

func StringToJSON(header string) Result {

	var jsonMap Result
	json.Unmarshal([]byte(header), &jsonMap)
	return jsonMap
}

// GetLocalIP returns all no loopback local IPs of the host
func GetLocalIP() []string {
	ipAdress := []string{"All Interfaces"}
	ifaces, err := net.Interfaces()
	if err != nil {
		return ipAdress
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ipAdress = append(ipAdress, ipnet.IP.String())
				}
			}
		}
	}
	return ipAdress
}
