package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/assignment-go/controller"
	"github.com/assignment-go/mock"
	"github.com/stretchr/testify/assert"
)

func init() {
	controller.Client = &mock.MockClient{}
}

func Test_main(t *testing.T) {
	// Create our test table
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		pindResponse  int
		expectedBody  string
	}{
		{
			description:   "ping gets 200",
			route:         "/upload",
			expectedError: false,
			expectedCode:  200,
			pindResponse:  200,
			expectedBody:  `[{"status":true,"url":"https://www.test.com"},{"status":true,"url":"http://www.google.com"},{"status":true,"url":"https://www.youtube.com"}]`,
		},
		{
			description:   "ping gets 404 error",
			route:         "/upload",
			expectedError: false,
			expectedCode:  200,
			pindResponse:  404,
			expectedBody:  `[{"status":false,"url":"https://www.test.com"},{"status":false,"url":"http://www.google.com"},{"status":false,"url":"https://www.youtube.com"}]`,
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

	// build response JSON
	json := ``
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	for _, tt := range tests {

		// create a new reader with that JSON
		mock.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: tt.pindResponse,
				Body:       r,
			}, nil
		}

		// Create a new http request with the route
		filePath := "mock/site.csv"
		fieldName := "file"
		bodyBuf := new(bytes.Buffer)

		mw := multipart.NewWriter(bodyBuf)

		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}

		w, err := mw.CreateFormFile(fieldName, filePath)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.Copy(w, file); err != nil {
			t.Fatal(err)
		}

		// close the writer before making the request
		mw.Close()

		req, _ := http.NewRequest(
			"POST",
			tt.route,
			bodyBuf,
		)
		// req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
		req.Header.Add("Content-Type", mw.FormDataContentType())

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, tt.expectedError, err != nil, tt.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if tt.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, tt.expectedCode, res.StatusCode, tt.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, tt.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, tt.expectedBody, string(body), tt.description)
	}
}
