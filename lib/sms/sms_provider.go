package sms

import "github.com/NikosSiak/Open-Banking-API/lib"

type SMSProvider interface {
	SendVerificationCode(to string) (string, error)
	VerifyCode(verificationId string, code string) (bool, error)
}

func NewSMSProvider(env lib.Env) SMSProvider {
	if env.IsProduction() {
		return newTwilioProvider(env.TwilioCredentials)
	}

	return newMockProvider()
}
