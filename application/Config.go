package application

import (
	"strconv"

	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-common/client"
	"github.com/smarest/smarest-common/infrastructure/persistence"
	"github.com/smarest/smarest-common/util"
)

type Bean struct {
	FileUploadService *FileUploadService
}

func InitBean() (*Bean, error) {
	maxFileSize, err := strconv.ParseInt(util.GetEnvDefault("POS_FILE_UPLOAD_SIZE", "10"), 0, 64)
	if err != nil {
		return nil, err
	}

	userTimeout, err := strconv.Atoi(util.GetEnvDefault("POS_USER_TIMEOUT", "5000"))
	if err != nil {
		return nil, err
	}

	loginClient := client.NewLoginClient(
		util.GetEnvDefault("POS_USER_HOST", "http://localhost:8080"),
		userTimeout,
	)

	loginService := application.NewLoginService(
		util.GetEnvDefault("POS_LOGIN_URL", "http://localhost:8080/login"),
		util.GetEnvDefault("POS_LOGIN_TOKEN", "pos_access_token"),
		persistence.NewLoginRepository(loginClient))

	fileUploadService := NewFileUploadService(
		loginService,
		util.GetEnvDefault("POS_FILE_UPLOAD_DIRECTORY", "./files"),
		maxFileSize)
	return &Bean{FileUploadService: fileUploadService}, nil
}
