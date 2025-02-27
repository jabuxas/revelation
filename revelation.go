package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	// "strings"
)

const pasteURL = "https://paste.jabuxas.com"

var key = os.Getenv("AUTH_KEY")

func main() {
	// file := strings.Split(SelectFile(), "file://")[1]
	file := "/yang/tmp/netin.json"
	request, err := uploadFile(file)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(respBody))
}

func uploadFile(file string) (*http.Request, error) {
	// open the file
	body, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer body.Close()

	// prepare multipart form data
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	defer writer.Close()

	// create the form file part
	part, err := writer.CreateFormFile("file", path.Base(file))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}

	// copy file content to form data
	_, err = io.Copy(part, body)
	if err != nil {
		return nil, fmt.Errorf("error copying file content: %v", err)
	}

	// create the HTTP request
	req, err := http.NewRequest("POST", pasteURL, data)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// set headers
	req.Header.Set("X-Auth", key)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(data.Len())

	return req, nil
}
