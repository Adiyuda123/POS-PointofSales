package helper

import (
	"POS-PointofSales/app/config"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
