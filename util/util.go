package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

// fileUploadRequest is create a file upload http request with optional extra params
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

// helper function to convert a string slice to a map.
func SliceToMap(s []string) map[string]bool {
	v := map[string]bool{}
	for _, ss := range s {
		v[ss] = true
	}
	return v
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0777)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, 0777)
			if err != nil {
				logrus.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
