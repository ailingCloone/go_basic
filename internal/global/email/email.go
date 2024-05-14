package email

import (
	"fmt"
	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/model/otp"
	"nrs_customer_module_backend/internal/model/ui"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/conf"
	"gopkg.in/gomail.v2"
)

// func SendEmail(to []string, body, subject, bodyType string, fileToAttach []string) error {
// 	if bodyType == "" {
// 		bodyType = "text/html"
// 	}

// 	var c config.Config
// 	configFile := "etc/api.yaml"
// 	if err := conf.LoadConfig(configFile, &c); err != nil {
// 		return nil
// 	}

// 	smtpDetails := c.SMTPEmail

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", smtpDetails.From)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody(bodyType, body)

// 	d := gomail.NewDialer(smtpDetails.Host, smtpDetails.Port, smtpDetails.Username, smtpDetails.Password)

// 	for _, i := range to {
// 		// Send the email
// 		m.SetHeader("To", i)

// 		for _, j := range fileToAttach {
// 			m.Attach(j)
// 		}

// 		if err := d.DialAndSend(m); err != nil {
// 			return err
// 		}
// 	}

// 	return nil

// }

// func SendOTPEmail(email, otp string, body, subject, bodyType string) (err error) {
// 	if bodyType == "" {
// 		bodyType = "text/plain"
// 	}

// 	var c config.Config
// 	configFile := "etc/api.yaml"
// 	if err := conf.LoadConfig(configFile, &c); err != nil {
// 		return err
// 	}

// 	smtpDetails := c.SMTPEmail

// 	// m := gomail.NewMessage()
// 	// m.SetHeader("From", smtpDetails.From)
// 	// m.SetHeader("Subject", subject)
// 	// m.SetBody(bodyType, body)

// 	// d := gomail.NewDialer(smtpDetails.Host, smtpDetails.Port, smtpDetails.Username, smtpDetails.Password)

// 	// m.SetHeader("To", email)

// 	// // Send the email
// 	// if err := d.DialAndSend(m); err != nil {
// 	// 	fmt.Println("the err", err)
// 	// 	return err
// 	// }
// 	msg := "From: " + smtpDetails.From + "\n" +
// 		"To: " + email + "\n" +
// 		"Subject: " + subject + "\n\n" +
// 		body

// 	addr := smtpDetails.Host + ":" + fmt.Sprint(smtpDetails.Port)

// 	// err = smtp.SendMail("smtp.gmail.com:587",
// 	err = smtp.SendMail(addr,
// 		smtp.PlainAuth("", smtpDetails.From, smtpDetails.Password, smtpDetails.Host),
// 		smtpDetails.From, []string{email}, []byte(msg))

// 	if err != nil {
// 		return err
// 	}

//		fmt.Println("Email sent successfully")
//		return err
//	}
func SendOTPEmail(email string, otp int64, body, subject, bodyType string) (err error) {
	if bodyType == "" {
		bodyType = "text/html"
	}

	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return err
	}

	smtpDetails := c.SMTPEmail

	// Create a new message
	message := gomail.NewMessage()
	message.SetHeader("From", smtpDetails.From)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)
	message.SetBody(bodyType, body)

	// Set up the dialer
	dialer := gomail.NewDialer(smtpDetails.Host, smtpDetails.Port, smtpDetails.From, smtpDetails.Password)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	fmt.Println("Email sent successfully")
	return nil

}

func EmailSentResponse(otpRecord *otp.Otp, sendStatus string) map[string]interface{} {
	dialog := ui.Dialog{
		Content: ui.DialogContent{
			ImageUrl:    "http://aaa",
			WebImageUrl: "http://aaa",
			Title:       "The password reset link has been sent successfully to your email. Please check your email to change the password.",
			Subtitle:    "If you did not receive the reset password email, or if you are having trouble resetting your password, please contact our customer service at +6016 299 13898 from 9.30AM to 5.30PM Monday to Friday.",
		},
		Button: []ui.Button{
			{
				Text:   "OK",
				Action: "/",
			},
		},
	}

	sendTimeStr := global.FormatMDYTime(otpRecord.SendTime)
	sendTimePtr := &sendTimeStr

	email_delivery := &types.EmailDelivery{
		Contact:    &otpRecord.Value,
		SendStatus: &sendStatus,
		SendTime:   sendTimePtr,
	}
	content := types.OtpDetailsContent{
		AuthType:      "EMAIL", //forgot password only email
		EmailDelivery: email_delivery,
		Message:       "The password reset link has been sent successfully to your email. Please check your email to change the password.",
		PhoneDelivery: &types.PhoneDelivery{}, //forgot password only use email for now.
		TxID:          otpRecord.Guid,
	}

	response := map[string]interface{}{
		"content": content,
		"dialog":  dialog,
	}

	return response

}
