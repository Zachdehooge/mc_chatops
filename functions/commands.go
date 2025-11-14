package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"errors"

	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

var startTime = time.Now()
var serverStartTime time.Time
var serverRunning bool
var databaseName = "servers.db"
var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func BotUptime() string {
	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func ServerStatus() string {
	godotenv.Load()
	log.Print("Getting bot token from .env file")
	server := os.Getenv("SERVERADD")

	log.Print("Fetching server information")
	url := fmt.Sprintf("http://%s", server)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("could not create request:", err)
		return "error"
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("could not fetch server:", err)
		return "error"
	}
	defer resp.Body.Close()

	return fmt.Sprintf("%d", resp.StatusCode)
}

func StartServer() string {
	serverStartTime = time.Now()
	serverRunning = true
	return "Starting Server..."
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

// TODO!: Set up table for servers and return them in the help command for servers that are connected

func AddServer(db *sql.DB, ip string) {
	if db == nil {
		log.Println("Database not initialized. Call SetDB() first.")
		return
	}

	log.Printf("Adding server IP: %s", ip)

	stmt, err := db.Prepare("INSERT OR IGNORE INTO servers (ip) VALUES (?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(ip)
	if err != nil {
		log.Printf("Failed to insert server IP: %v", err)
		return
	}

	log.Printf("Server %s added successfully!", ip)
}

func RemoveServer(ip string) {
	log.Printf("Removing Server IP: %s", ip)
}

func GetServers() []string {
	if db == nil {
		return []string{}
	}

	rows, err := db.Query("SELECT ip FROM servers")
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var servers []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return []string{}
		}
		servers = append(servers, ip)
	}

	if err := rows.Err(); err != nil {
		return []string{}
	}

	return servers
}

func DatabaseDestroy() {
	time.Sleep(10 * time.Second)

	doesFileExist := checkFileExists(databaseName)

	if doesFileExist {
		log.Printf("Tearing down database %s...", databaseName)
		os.Remove(databaseName)
		log.Printf("Database %s destroyed successfully!", databaseName)
	} else {
		log.Printf("Database %s does not exist", databaseName)
	}
}
