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
var databaseName = "servers.sqlite3"

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

func DatabaseInit() {
	os.Create(databaseName)

	var present = checkFileExists(databaseName)

	if present {
		log.Print("Database Found!")
	} else {
		log.Print("Database Not Found")
	}

	db, err := sql.Open("sqlite3", "./"+databaseName)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	log.Printf("Conencted to database %s", databaseName)
}

func CheckDBHealth((h *Handler) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	return nil
}

// TODO: Set up table for servers and return them in the help command for servers that are connected

func AddServer(ip string) {
	log.Printf("Adding Server IP: %s", ip)
}

func RemoveServer(ip string) {
	log.Printf("Removing Server IP: %s", ip)
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
