package main

import (
	"fmt"
	"net/smtp"
)

func sendMailSimple(subject string, body string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"suryak14919@gmail.com",
		"cpooltnpqnpwfirg",
		"smtp.gmail.com",
	)

	msg := "Subject :" + subject + "\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"suryak14919@gmail.com",
		// []string{"ravitejak1900@gmail.com"},
		// poojasyo@gmail.com
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}

func sendMailSimpleHTML(subject string, html string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"suryak14919@gmail.com",
		"cpooltnpqnpwfirg",
		"smtp.gmail.com",
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject :" + subject + "\n" + headers + "\n\n" + html

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"suryak14919@gmail.com",
		// []string{"ravitejak1900@gmail.com"},
		// poojasyo@gmail.com
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("Hello World...")
	sendMailSimple("Well Wisher", "Hi Get Well Soon Buddy, Keep Smiling, Don't be in a hurry keep calm(Peace)", []string{"basavaraja@synergytechs.net"})

}
