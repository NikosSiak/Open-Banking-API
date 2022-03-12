package sms

import "github.com/google/uuid"

const mockCode = "1312"

type mockProvider struct {
	code string
}

func newMockProvider() mockProvider {
	return mockProvider{
		code: mockCode,
	}
}

func (m mockProvider) SendVerificationCode(to string) (string, error) {
	return uuid.NewString(), nil
}

func (m mockProvider) VerifyCode(verificationId string, code string) (bool, error) {
	return m.code == code, nil
}
