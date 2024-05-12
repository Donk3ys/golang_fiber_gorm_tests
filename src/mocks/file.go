package mocks_test

import (
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type FileSysetmClient struct {
	mock.Mock
}

func (m *FileSysetmClient) UpdateOrCreateFileFromBase64(existingFileUrl string, filePath string, fileBytes string) error {
	log.Printf("[FILE BASE64]: existingFileUrl:%s, filePath:%s, fileBytes-length:%d\n", existingFileUrl, filePath, len(fileBytes))

	args := m.Called(existingFileUrl, filePath, fileBytes)
	return args.Error(0)
}

func (m *FileSysetmClient) UpdateOrCreateFileFromForm(existingFileUrl string, filePath string, file *multipart.FileHeader, c *fiber.Ctx) error {
	log.Printf("[FILE MULTPART]: existingFileUrl:%s, filePath:%s\n", existingFileUrl, filePath)

	args := m.Called(existingFileUrl, filePath, file, c)
	return args.Error(0)
}

func (m *FileSysetmClient) RemoveFile(url string) error {
	log.Printf("[FILE REMOVE]: url:%s\n", url)

	args := m.Called(url)
	return args.Error(0)
}

func (m *FileSysetmClient) CreateFilePathAndUrl(folderPath string, newFileName string, fileExt string) (string, string) {
	log.Printf("[FILE CREATE URL]: folderPath:%s, newFileName:%s, fileExt:%s\n", folderPath, newFileName, fileExt)

	args := m.Called(folderPath, newFileName, fileExt)
	return args.String(0), args.String(1)
}
