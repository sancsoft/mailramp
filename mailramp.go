package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

type Config struct {
	Sender string `json:"sender"`
	Rate   int    `json:"rate"`
	Server struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"server"`
	Subject string `json:"subject"`
}

func readFileToArray(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func loadConfig() (Config, error) {

	buffer, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func main() {

	fmt.Println("MAILRAMP: Warm up an email domain with messages to a set of email adddresses.")
	fmt.Println(")|( Sanctuary Software Studio, Inc. - All rights reserved.")
	fmt.Println("")
	// load the configuration from JSON
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error reading config: ", err)
	}

	// read the collection of email addresses
	recipients, err := readFileToArray("recipients.txt")
	if err != nil {
		fmt.Println("Error reading recipients: ", err)
	}
	if len(recipients) == 0 {
		fmt.Println("Error no recipients loaded")
		return
	}
	fmt.Printf("Loaded %d recipients\n", len(recipients))

	// read in the message text
	body, err := os.ReadFile("body.txt")
	if err != nil {
		fmt.Println("Error reading message body:", err)
		return
	}

	// construct the message
	message := fmt.Sprintf("From: %s\r\n", config.Sender)
	message += fmt.Sprintf("Subject: %s\r\n", config.Subject)
	message += "\r\n" + string(body[:])

	// Authentication
	auth := smtp.PlainAuth("", config.Server.User, config.Server.Password, config.Server.Host)
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	// send the message to each recipient
	count := 0
	for _, recipient := range recipients {
		count += 1
		fmt.Printf("%s: %05d: %s\n", time.Now().Format(time.RFC3339), count, recipient)
		msg := fmt.Sprintf("To: %s\r\n", recipient) + message
		err = smtp.SendMail(addr, auth, config.Sender, []string{recipient}, []byte(msg))
		if err != nil {
			fmt.Println("Error sending email:", err)
			os.Exit(1)
		}
	}
	fmt.Printf("%d email messages sent successfully", count)
}
