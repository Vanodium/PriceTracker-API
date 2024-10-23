package notifications

import (
	"log"

	"os"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

func EmailNotification(userEmail string, link string) {
	err := godotenv.Load("../../internal/transport/gmailcfg.env")
	if err != nil {
		log.Fatal(err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "pricetracker.sup@gmail.com")
	message.SetHeader("To", userEmail)
	message.SetHeader("Subject", "Price change detected")
	message.SetBody("text/plain", "Your tracker found price changing at "+link)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SENDER_EMAIL"), os.Getenv("APP_PASSWORD"))
	err = dialer.DialAndSend(message)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Email sent to %s with link %s", userEmail, link)
}
