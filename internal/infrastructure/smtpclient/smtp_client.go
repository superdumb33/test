package smtpclient

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type SMTPClient struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPClient(from string, host string, stringPort string, username, password string) *SMTPClient {
	port, _ := strconv.Atoi(stringPort)
	return &SMTPClient{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (client *SMTPClient) Send(to []string, subject, body string) error {
	message := 
		fmt.Sprintf("From: %s\r\n", client.From) + 
        fmt.Sprintf("To: %s\r\n", to) + 
        fmt.Sprintf("Subject: %s\r\n", subject) + 
        "MIME-Version: 1.0\r\n" + 
        "Content-Type: text/plain; charset=\"UTF-8\"\r\n" + 
        "\r\n" +
		body

		address := fmt.Sprintf("%s:%d", client.Host, client.Port)
	
		return smtp.SendMail(address, nil, client.From, to, []byte(message))
}
