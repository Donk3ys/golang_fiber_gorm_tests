package storage

import (
	"api/src/constants"
	"api/src/util"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FileSystem struct {
	Client IFileSysetmClient
}

type IFileSysetmClient interface {
	UpdateOrCreateFileFromBase64(existingFileUrl string, filePath string, fileBytes string) error
	UpdateOrCreateFileFromForm(existingFileUrl string, filePath string, file *multipart.FileHeader, c *fiber.Ctx) error
	RemoveFile(url string) error
	CreateFilePathAndUrl(folderPath string, newFileName string, fileExt string) (string, string)
}

const USER_PROFILE_PICS_PATH = "/public/users/profile_images/"

type LocalStorage struct{}

// Constructor to make folders
func NewLocalStorage() LocalStorage {
	os.MkdirAll("."+USER_PROFILE_PICS_PATH, os.ModePerm)

	return LocalStorage{}
}

func (ls LocalStorage) UpdateOrCreateFileFromBase64(existingFileUrl string, filePath string, fileBytes string) error {
	bytes, err := base64.StdEncoding.DecodeString(fileBytes)
	if err != nil {
		return err
	}

	// Create file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	// Remove old file if already exists
	if existingFileUrl != "" {
		ls.RemoveFile(existingFileUrl)
	}

	return nil
}

func (ls LocalStorage) UpdateOrCreateFileFromForm(existingFileUrl string, filePath string, file *multipart.FileHeader, c *fiber.Ctx) error {
	err := c.SaveFile(file, "."+filePath)
	if err != nil {
		return err
	}

	// Remove old file if already exists
	if existingFileUrl != "" {
		ls.RemoveFile(existingFileUrl)
	}

	return nil
}

func (LocalStorage) RemoveFile(url string) error {
	params := strings.Split(url, "/")
	exFilename := params[len(params)-1]
	exFolderPath := params[len(params)-2]
	exFilePath := fmt.Sprintf("%s%s", exFolderPath, exFilename)
	return os.Remove(exFilePath)
}

func (LocalStorage) CreateFilePathAndUrl(folderPath string, newFileName string, fileExt string) (string, string) {
	uniqueCode := util.GenerateCode()
	filePath := fmt.Sprintf("%s/%s--%d.%s", folderPath, newFileName, uniqueCode, fileExt)
	newUrl := fmt.Sprintf("%s%s", constants.API_URL, filePath)
	return filePath, newUrl
}
