package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const (
	WeiboToken    = "2.00PfzJxC5FNO8D0dbb450c3bTtl59E"
	WeiboShareURL = "https://api.weibo.com/2/statuses/share.json"
)

func Upload(WeiboShareURL, fileNameOrPath, shareContent string) (err error) {
	// Prepare a form that you will submit to that URL.
	var postMultiFormBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&postMultiFormBody)
	// Add your image fileNameOrPath
	f, err := os.Open(fileNameOrPath)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := multipartWriter.CreateFormFile("pic", fileNameOrPath)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	multipartWriter.WriteField("access_token", WeiboToken)
	multipartWriter.WriteField("status", url.PathEscape(shareContent))

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	multipartWriter.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", WeiboShareURL, &postMultiFormBody)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	bodyString := string(bodyBytes)
	if res.StatusCode != http.StatusOK {

		log.Printf("bad status: %s,\n%s", res.Status, bodyString)

	}

	log.Println(bodyString)
	return
}

func main() {
	Upload(WeiboShareURL, "/Users/ericzhou/go/src/github.com/fogleman/primitive/examples/monalisa.png", "hello_world jjj !")
}
