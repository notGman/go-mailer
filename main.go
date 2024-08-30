package main

import (
	"bytes"
	"encoding/csv"
	"log"
	"net/smtp"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

func main() {
	mails := []string{}

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	f, err := os.Open("mails.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rd := csv.NewReader(f)
	records, err := rd.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		mails = append(mails, record[0])
	}

	sender, _ := os.LookupEnv("SENDER")
	pass, _ := os.LookupEnv("PASS")

	for _, mail := range mails {
		from := sender
		appPass := pass
		to := []string{
			mail,
		}
		s := "smtp.gmail.com:587"
		ms, err := parseTemplate("index.html", nil)
		if err != nil {
			log.Fatalln("Could not parse template file", err)
		}

		b := getMessageString(from, to[0], "go mail automater", ms)
		auth := smtp.PlainAuth("", from, appPass, "smtp.gmail.com")
		err = smtp.SendMail(s, auth, from, to, b)
		if err != nil {
			log.Fatalln("Could not send email", err)
		}
		log.Println("Mails sent")
	}
}

func getMessageString(from, to, subject, body string) []byte {
	headers := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	return []byte(headers + body)
}

func parseTemplate(templateName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		log.Fatalln("Could not parse template", err)
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
