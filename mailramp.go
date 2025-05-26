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

func sendEmail(host string, auth smtp.Auth, from string, to string, message string) error {
	msg := fmt.Sprintf("To: %s\r\n", to) + message
	err := smtp.SendMail(host, auth, from, []string{to}, []byte(msg))
	return err
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
	host := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	// convert the messages per hour in to full and partial messages per minute
	mpm := config.Rate / 60
	ppm := config.Rate % 60
	ppmerror := (ppm << 1) - 60

	count := 0

	for {
		// capture the start time and set the end time
		startTime := time.Now()
		endDuration, err := time.ParseDuration("60s")
		if err != nil {
			fmt.Println("Error parsing duration:", err)
			os.Exit(1)
		}

		// send the base messages per minute
		for m := 0; m < mpm; m++ {
			recipient := recipients[count%len(recipients)]
			count += 1
			fmt.Printf("%s: %05d: %s\n", time.Now().Format(time.RFC3339), count, recipient)
			sendEmail(host, auth, config.Sender, recipient, message)
			if err != nil {
				fmt.Println("Error sending email:", err)
				os.Exit(1)
			}
		}

		// send another message if error has accumulated based on ppm
		if ppmerror > 0 {
			recipient := recipients[count%len(recipients)]
			count += 1
			fmt.Printf("%s: %05d: %s\n", time.Now().Format(time.RFC3339), count, recipient)
			sendEmail(host, auth, config.Sender, recipient, message)
			if err != nil {
				fmt.Println("Error sending email:", err)
				os.Exit(1)
			}
			ppmerror -= (60 << 1)
		}

		// increase the ppm error
		ppmerror += (ppm << 1)

		// wait out the rest of the minute
		for {
			if time.Since(startTime) > endDuration {
				break
			}
			time.Sleep(time.Second)
		}
	}
}
