package file

import (
	"mime/multipart"
	"strings"
)

func IsImage(file *multipart.FileHeader) bool {
	return strings.Split(file.Header.Get("Content-Type"), "/")[0] == "image"
}
