package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/", handleHealth)
	r.POST("/sendmail", handleSendEmail)

	port := ":" + os.Getenv("PORT")
	r.Run(port)

  
}

type sendMailRequest struct {
	DestEmail string `json:"dest_email"`
	Details orderDetails `json:"details"`

}

type orderDetails struct {
	ID string `json:"id"`
	Amount string `json:"amount"`
	Name string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Adress string `json:"address"`
	ItemCount string `json:"item_count"`
	Items []string `json:"items"`

}

func handleHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "met me well!")
}

func handleSendEmail(ctx *gin.Context) {

	var req sendMailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := sendMail(req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	
	ctx.JSON(http.StatusOK, res)

}

func sendMail(req sendMailRequest) (string, error) {
	// Sender data.
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("PASSWORD")
  
	// Receiver email address.
	to := []string{
	  req.DestEmail,
	}
  
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
  
	req.Details.ID, _ = newUUID()

	// Message.
	message := []byte(req.Details.Amount)
	
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
	  return "", err
	}

	return req.Details.ID, err
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}