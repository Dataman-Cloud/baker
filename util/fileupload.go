package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// FileUploadRequest is create a file upload http request with optional extra params
func FileUploadRequest(uri string, token, uploadParamName, uploadFilePath string, extraParams map[string]string) (*http.Request, error) {
	file, err := os.Open(uploadFilePath)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(uploadParamName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContent)

	for key, val := range extraParams {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	contentType := writer.FormDataContentType()

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", contentType)
	return req, err
}
