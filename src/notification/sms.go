package notification

// import (
// 	"fmt"
// 	"os"
//
// 	"github.com/gofiber/fiber/v2/log"
// 	"github.com/twilio/twilio-go"
// 	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
// )

type SMS struct {
	Client ISmsClient
}

// func (s *SMS) SendUserMobileVerificationCode(mobile string, email string, code uint16) error {
// 	msg := fmt.Sprintf("Your FuelU OTP %d for %s mobile verification", code, email)
// 	return sendSMS(mobile, msg)
// }
//
// func (s *SMS) SendUserLoginVerificationCode(mobile string, email string, code uint16) error {
// 	msg := fmt.Sprintf("Your FuelU OTP %d for %s login", code, email)
// 	return sendSMS(mobile, msg)
// }
//

// Interface for test mock
type ISmsClient interface {
	SendSMS(to string, msg string) error
}

type Twillio struct {
	Username string
	Password string
}

func (t Twillio) SendSMS(to string, msg string) error {
	// client := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: t.Username,
	// 	Password: t.Password,
	// })
	//
	// params := &twilioApi.CreateMessageParams{}
	// params.SetTo(to)
	// // params.SetFrom("+15017250604")
	// params.SetMessagingServiceSid(t.Username)
	// params.SetBody(msg)
	//
	// _, err := client.Api.CreateMessage(params)
	// if err != nil {
	// 	log.Errorf("[Twillio Error]: %s", err)
	// 	return err
	// }
	// // response, _ := json.Marshal(*resp)
	// // fmt.Println("Response: " + string(response))
	return nil
}
