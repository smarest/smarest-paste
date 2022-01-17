package application

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

// FileUploadService upload file
type FileUploadService struct {
	*application.LoginService
	Directory   string
	MaxFileSize int64
	FileType    []string
}

// NewFileUploadService create FileUploadService
func NewFileUploadService(loginService *application.LoginService, directory string, maxFileSize int64) *FileUploadService {

	return &FileUploadService{
		loginService, directory, maxFileSize * 1000000,
		[]string{".jpg", ".JPG", ".png", ".PNG", ".jpeg", ".JPEG", ".gif", ".GIF"}}
}

// Post post
func (s *FileUploadService) Post(c *gin.Context) {
	_, lErr := s.CheckCookie(c)
	if lErr != nil {
		log.Print(lErr.ErrorMessage)
		c.JSON(http.StatusUnauthorized, exception.CreateError(exception.CodeSignatureInvalid, "Access denied."))
		return
	}

	directory := c.PostForm("directory")
	if directory != "" {
		directory = filepath.Base(directory)
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "file required."))
		return
	}

	//check file Size
	if file.Size > s.MaxFileSize {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, fmt.Sprintf("Sorry, your file is too large. < %fMB are allowed: ", float64(s.MaxFileSize/1000000))))
		return
	}

	//check file name

	baseFileName := filepath.Base(file.Filename)
	if s.IsInvalidFileExtension(baseFileName) {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeSignatureInvalid, "Sorry, only JPG, JPEG, PNG & GIF files are allowed."))
		return
	}

	filename := filepath.Join(s.Directory, directory, baseFileName)
	err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, exception.CreateError(exception.CodeUnknown, "can not create file folder."))
		return
	}

	//check file is exits
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "Sorry, file already exists."))
		return
	}

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, exception.CreateError(exception.CodeUnknown, fmt.Sprintf("upload file err: %s", err.Error())))
		return
	}
	c.JSON(http.StatusOK, gin.H{"filePath": c.Request.URL.String() + "/view/" + directory + "/" + baseFileName, "message": fmt.Sprintf("'%s' uploaded!", file.Filename)})
}

func (s *FileUploadService) IsInvalidFileExtension(fileName string) bool {
	ext := filepath.Ext(fileName)
	for _, v := range s.FileType {
		if v == ext {
			return false
		}
	}
	return true
}
