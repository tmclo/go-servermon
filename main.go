package main

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"
)

const user = "b3e751a23e3605"
const password = "d5c7f01a2aa78d"
const sender = "no-reply@example.com"

var to = []string{
	"tom@springsy.co.uk",
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func main() {
	fmt.Println("Server Running...")
	for range time.Tick(time.Minute * 2) {
		go func() {
			fmt.Println("Checking hosts...")
			ReadHosts("hosts.txt")
		}()
	}
}

func ReadHosts(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Ping(scanner.Text())
		time.Sleep(10 * time.Second)
	}
	return nil, scanner.Err()
}

func Ping(address string) {
	out, _ := exec.Command("ping", address, "-c 1", "-i 3", "-w 10").Output()
	if strings.Contains(string(out), "Destination Host Unreachable") {
		fmt.Println("Server offline: " + address)
		subject := "SERVER OFFLINE: " + address
		body := "<p>Warning, we have detected your server <b>" + address + "</b> is currently offline.</p>"
		SendMail(subject, body)
	} else {
		fmt.Println("Server " + address + " Online, all checks passed...")
	}
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func SendMail(subject string, body string) {

	request := Mail{
		Sender:  sender,
		To:      to,
		Subject: subject,
		Body:    body,
	}
	addr := "smtp.mailtrap.io:2525"
	host := "smtp.mailtrap.io"

	msg := BuildMessage(request)
	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(addr, auth, sender, to, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
