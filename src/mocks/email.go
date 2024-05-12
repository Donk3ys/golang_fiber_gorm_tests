package mocks_test

import (
	"log"

	"github.com/stretchr/testify/mock"
)

type EmailCLient struct {
	mock.Mock
}

func (m *EmailCLient) SendMail(from string, to []string, subject string, templatePath string, data interface{}) error {
	log.Printf("[SEND MAIL]: from:%s, to:%v, subject:%s, tmpPath:%s, data:%v\n", from, to, subject, templatePath, data)

	args := m.Called(from, to, subject, templatePath, data)
	return args.Error(0)
}
