package utils

import (
	"crypto/tls"
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", viper.GetString("EMAIL_ID"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, viper.GetString("EMAIL_ID"), viper.GetString("EMAIL_PASSWORD"))

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
