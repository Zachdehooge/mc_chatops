package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	h "github.com/zachdehooge/MC-Chatops/functions"
)

var db *sql.DB

// Global Variables
var s *discordgo.Session

func init() {
	godotenv.Load()
	log.Print("Getting bot token from .env file")
	var BotToken = os.Getenv("TOKEN")
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v | Check the .env", err)
	}

}

// Slash Commands
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "botstatus",
			// All commands and options must have a description
			Description: "bot uptime",
		},
		{
			Name:        "serverstatus",
			Description: "server uptime",
		},
		{
			Name:        "serverstart",
			Description: "starts the minecraft server",
		},
		{
			Name:        "serverstop",
			Description: "stops the minecraft server",
		},
		{
			Name:        "serverscale",
			Description: "scales the minecraft server | default is auto",
		},
		{
			Name:        "databasestart",
			Description: "starts the database",
		},
		{
			Name:        "databasestop",
			Description: "stops the database",
		},
		{
			Name:        "listservers",
			Description: "lists the servers in the database",
		},
		{
			Name:        "addserver",
			Description: "adds a server to the database",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ip",
					Description: "server IP to add to the database",
					Required:    true,
				},
			},
		},
		{
			Name:        "removeserver",
			Description: "removes a server to the database",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ip",
					Description: "server IP to remove to the database",
					Required:    true,
				},
			},
		},
		{
			Name:        "help",
			Description: "help",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"botstatus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Bot Uptime",
							Description: fmt.Sprintf("Bot Uptime: %s", h.BotUptime()),
							Color:       0x57F287,
						},
					},
				},
			})
		},
		"serverstatus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Server Uptime",
							Description: fmt.Sprintf("Server Uptime: %s\nServer Status Code: %s", h.ServerUptime(), h.ServerStatus()),
							Color:       h.ColorStatus(),
						},
					},
				},
			})
		},
		"startserver": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			h.StartServer()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Server Start",
							Description: "Starting Server...",
							Color:       0x57F287,
						},
					},
				},
			})
		},
		"stopserver": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			h.StopServer()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Server Stop",
							Description: "Stopping server...",
							Color:       0xFF0000,
						},
					},
				},
			})
		},
		"scaleserver": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Server Scale",
							Description: "Scaling server...",
							Color:       0xADD8E6,
						},
					},
				},
			})
		},
		// TODO!: Depracate this?
		"databasestart": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//h.DatabaseInit()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Initializing Database...",
							// TODO!: Should be auto started by default
							Color: 0xADD8E6,
						},
					},
				},
			})
		},
		"addserver": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ip := i.ApplicationCommandData().Options[0].StringValue()
			h.AddServer(db, ip)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: fmt.Sprintf("Adding %s to the database...", ip),
							// TODO!: Need to make a function in commands.go to fetch the database results and check that against what was passed in the command to ensure the IP successfully was ADDED to the database
							Color: 0x39ff02,
						},
					},
				},
			})
			//commandHandlers["listservers"](s, i)
		},
		"removeserver": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ip := i.ApplicationCommandData().Options[0].StringValue()
			h.RemoveServer(ip)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: fmt.Sprintf("Removing %s from the database...", ip),
							// TODO!: Need to make a function in commands.go to fetch the database results and check that against what was passed in the command to ensure the IP successfully was REMOVED to the database
							Color: 0xff0206,
						},
					},
				},
			})
			//commandHandlers["listservers"](s, i)
		},
		"listservers": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Listing Servers Database...",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Servers",
									Value:  strings.Join(h.GetServers(), "\n"),
									Inline: false,
								},
							},
							Color: 0xADD8E6,
						},
					},
				},
			})
		},
		"databasestop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			h.DatabaseDestroy()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Destroying Database...",
							Color: 0xADD8E6,
						},
					},
				},
			})
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "List of Commands",
							Color: 0xFF0090,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "/botstatus",
									Value:  "Shows bot uptime",
									Inline: false,
								},
								{
									Name:   "/serverstatus",
									Value:  "Shows server uptime and status",
									Inline: false,
								},
								{
									Name:   "/serverstart",
									Value:  "Starts the Minecraft server",
									Inline: false,
								},
								{
									Name:   "/serverstop",
									Value:  "Stops the Minecraft server",
									Inline: false,
								},
								{
									Name:   "/serverscale",
									Value:  "Scales the Minecraft server",
									Inline: false,
								},
								{
									Name:   "/addserver",
									Value:  "Adds a server to the database",
									Inline: false,
								},
								{
									Name:   "/removeserver",
									Value:  "Removes a server from the database",
									Inline: false,
								},
								{
									Name:   "/listservers",
									Value:  "Lists all servers in the database",
									Inline: false,
								},
								{
									Name:   "/databasestart",
									Value:  "Starts the database for the bot",
									Inline: false,
								},
								{
									Name:   "/databasedestroy",
									Value:  "Destroys the database for the bot",
									Inline: false,
								},
							},
						},
						{
							Title: "Servers",
							Color: 0xFF0090,
							Fields: []*discordgo.MessageEmbedField{
								{
									Value:  "SERVER 1",
									Inline: false,
								},
								{
									Value:  "SERVER 2",
									Inline: false,
								},
								{
									Value:  "SERVER 3",
									Inline: false,
								},
							},
						},
					},
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {

	var GuildID = os.Getenv("GuildID")

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "Your Minecraft Server",
					Type: discordgo.ActivityTypeWatching,
				},
			},
			Status: "online",
		})
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	existing, err := s.ApplicationCommands(s.State.User.ID, GuildID)
	if err != nil {
		log.Fatalf("Failed to list existing commands: %v", err)
	}

	for _, cmd := range existing {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, cmd.ID)
		if err != nil {
			log.Printf("Failed to delete old command '%v': %v", cmd.Name, err)
		} else {
			//log.Printf("Deleted old command: %v", cmd.Name)
		}
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	log.Println("Refreshing commands...")
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildID, commands)
	if err != nil {
		log.Fatalf("Cannot refresh commands: %v", err)
	}

	h.SetDB(db)

	db, err := sql.Open("sqlite3", "./servers.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL UNIQUE
	);`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Print("Database initialized successfully!")

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Ready to take commands!")
	log.Println("Press Ctrl+C to exit")
	<-stop
	log.Println("Gracefully shutting down.")
}
