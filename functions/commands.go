package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var startTime = time.Now()
var serverStartTime time.Time
var serverRunning bool

/* var databaseName = "servers.db"
var db *sql.DB */

type IPStore struct {
	IPs []string `json:"ips"`
}

func BotStart() string {
	return startTime.Format(time.DateTime)
}

func BotUptime() string {
	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func Ping(ip string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.Contains(string(out), "TTL=")
	}
	return strings.Contains(string(out), "ttl=")
}

func ServerStatus() string {

	store, err := Load()
	if err != nil {
		return "Error loading IPs..."
	}
	ips := store.GetIPs()
	if len(ips) == 0 {
		return "No IPs Stored..."
	}

	result := ""
	for _, ip := range ips {
		if Ping(ip) {
			result += fmt.Sprintf("🟢 %s\n", ip)
		} else {
			result += fmt.Sprintf("🔴 %s\n", ip)
		}
	}
	return result
}

func StartServer() string {
	serverStartTime = time.Now()
	serverRunning = true
	return "Starting Server..."
}

func RestartServer() string {
	log.Println("Stopping Server...")
	serverRunning = false
	log.Println("Starting Server...")
	serverRunning = true
	return "Restarted Server Successfully..."
}

func StopServer() string {
	serverRunning = false
	return "Stopping Server..."
}

func ServerUptime() string {
	if serverRunning {
		uptime := time.Since(serverStartTime)
		hours := int(uptime.Hours())
		minutes := int(uptime.Minutes()) % 60
		seconds := int(uptime.Seconds()) % 60
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		return "server is down..."
	}
}

func ColorStatus() int {
	if ServerStatus() == "200" {
		return 0x57F287
	} else {
		return 0xFF0000
	}
}

func Load() (*IPStore, error) {
	data, err := os.ReadFile("ips.json")
	if err != nil {
		if os.IsNotExist(err) {
			return &IPStore{IPs: []string{}}, nil
		}
		return nil, err
	}

	var store IPStore
	err = json.Unmarshal(data, &store)
	return &store, err
}

func (s *IPStore) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("ips.json", data, 0644)
}

func (s *IPStore) AddIP(ip string) {
	s.IPs = append(s.IPs, ip)
}

func (s *IPStore) RemoveIP(ip string) {
	newList := []string{}
	for _, v := range s.IPs {
		if v != ip {
			newList = append(newList, v)
		}
	}
	s.IPs = newList
}

func (s *IPStore) GetIPs() []string {
	return s.IPs
}
