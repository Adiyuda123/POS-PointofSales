package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	JWT                    string
	CloudinaryName         string
	CloudinaryApiKey       string
	CloudinaryApiScret     string
	CloudinaryUploadFolder string
	XenditSecretKey        string
)

type AppConfig struct {
	DBUSER     string
	DBPASSWORD string
	DBHOST     string
	DBPORT     string
	DBNAME     string
	JWT        string
}

func InitConfig() *AppConfig {
	return &AppConfig{}
}

func ReadEnv() *AppConfig {
	app := &AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("DBUSER"); found {
		app.DBUSER = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPASSWORD"); found {
		app.DBPASSWORD = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHOST"); found {
		app.DBHOST = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPORT"); found {
		app.DBPORT = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBNAME"); found {
		app.DBNAME = val
		isRead = false
	}

	if val, found := os.LookupEnv("JWT"); found {
		JWT = val
		isRead = false
	}

	if val, found := os.LookupEnv("CLOUDINARY_CLOUD_NAME"); found {
		CloudinaryName = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_API_KEY"); found {
		CloudinaryApiKey = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_API_SECRET"); found {
		CloudinaryApiScret = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_UPLOAD_FOLDER"); found {
		CloudinaryUploadFolder = val
		isRead = false
	}
	if val, found := os.LookupEnv("XENDIT_SECRET_KEY"); found {
		XenditSecretKey = val
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error reading config:", err.Error())
			return nil
		}

		app.DBUSER = viper.GetString("DBUSER")
		app.DBPASSWORD = viper.GetString("DBPASSWORD")
		app.DBHOST = viper.GetString("DBHOST")
		app.DBPORT = viper.GetString("DBPORT")
		app.DBNAME = viper.GetString("DBNAME")

		JWT = viper.GetString("JWT")

		CloudinaryName = viper.GetString("CLOUDINARY_CLOUD_NAME")
		CloudinaryApiKey = viper.GetString("CLOUDINARY_API_KEY")
		CloudinaryApiScret = viper.GetString("CLOUDINARY_API_SECRET")
		CloudinaryUploadFolder = viper.GetString("CLOUDINARY_UPLOAD_FOLDER")
		XenditSecretKey = viper.GetString("XENDIT_SECRET_KEY")
	}

	return app
}
