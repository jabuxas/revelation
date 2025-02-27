package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/atotto/clipboard"
)

var pasteURL = os.Getenv("PASTEBIN_URL")
var key = os.Getenv("AUTH_KEY")

func main() {
	if pasteURL == "" || key == "" {
		fmt.Println("Please set PASTEBIN_URL and AUTH_KEY environment variables")
		return
	}

	file := strings.Split(SelectFile(), "file://")[1]
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

	if err := clipboard.WriteAll(string(respBody)); err != nil {
		panic(err)
	}
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
	writer.Close()

	// create the HTTP request
	req, err := http.NewRequest("POST", pasteURL, data)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// set headers
	req.Header.Set("X-Auth", key)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
