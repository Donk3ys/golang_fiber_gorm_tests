package mocks_test

import (
	"log"

	"github.com/stretchr/testify/mock"
)

type SmsCLient struct {
	mock.Mock
}

func (m *SmsCLient) SendSMS(to string, msg string) error {
	log.Printf("[SEND SMS]: to:%s, msg:%s\n", to, msg)

	args := m.Called(to, msg)
	return args.Error(0)
}
