package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
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

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	// For macOS and other OSes, return the HOME environment variable
	return os.Getenv("HOME")
}

func GetRenterdDefaultPath() string {
	home := UserHomeDir()
	if runtime.GOOS == "windows" {
		return home + "\\AppData\\Roaming\\Renterd"
	} else if runtime.GOOS == "linux" {
		return home + "/.config/renterd"
	} else if runtime.GOOS == "darwin" {
		return home + "/Library/Application Support/Renterd"
	}
	return ""
}

func GetDefaultSqliteBackupPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	if runtime.GOOS == "windows" {
		fmt.Println(path + "\\backup\\renterd.sqlite3.bak")
		return path + "\\backup\\renterd.sqlite3.bak"
	} else if runtime.GOOS == "linux" {
		return path + "/backup/renterd.sqlite3.bak"
	} else if runtime.GOOS == "darwin" {
		return path + "/backup/renterd.sqlite3.bak"
	}
	return ""
}

func GetSqliteDbDefautPath() string {
	path := GetRenterdDefaultPath()
	if runtime.GOOS == "windows" {
		return path + "\\data\\db\\db.sqlite"
	} else if runtime.GOOS == "linux" {
		return path + "/data/db/db.sqlite"
	} else if runtime.GOOS == "darwin" {
		return path + "/data/db/db.sqlite"
	}
	return ""
}
