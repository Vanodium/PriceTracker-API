package notifications

import (
	"log"

	"os"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

func EmailNotification(userEmail string, link string) {
	err := godotenv.Load("../../internal/transport/enviroment.env")
	if err != nil {
		log.Fatal(err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("GMAIL_SENDER"))
	message.SetHeader("To", userEmail)
	message.SetHeader("Subject", "Price change detected")
	message.SetBody("text/plain", "Your tracker found price changing at "+link)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("GMAIL_SENDER"), os.Getenv("GMAIL_APP_PASSWORD"))
	err = dialer.DialAndSend(message)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Email sent to %s with link %s", userEmail, link)
}
