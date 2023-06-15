package helper

import (
	"POS-PointofSales/app/config"
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
)

func UploadFile(fileContents interface{}, path string) ([]string, error) {
	var urls []string
	switch cnv := fileContents.(type) {
	case []multipart.File:
		for _, content := range cnv {
			uploadResult, err := getLink(content, path)
			if err != nil {
				return nil, err
			}
			urls = append(urls, uploadResult)
		}
	case multipart.File:
		uploadResult, err := getLink(cnv, path)
		if err != nil {
			return nil, err
		}
		urls = append(urls, uploadResult)
	}
	return urls, nil
}

func getLink(content multipart.File, path string) (string, error) {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinaryApiScret)
	if err != nil {
		return "", err
	}
	uploadParams := uploader.UploadParams{
		Folder: config.CloudinaryUploadFolder + path,
	}
	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		content,
		uploadParams)
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func ValidImageFormat(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
