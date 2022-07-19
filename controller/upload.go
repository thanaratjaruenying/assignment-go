package controller

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type healthResult struct {
	Status bool   `json:"status"`
	Url    string `json:"url"`
}

func validateUrl(rawUrl string) (string, error) {
	re := regexp.MustCompile(`(?s)https?://`)
	found := re.FindStringIndex(rawUrl)
	if len(found) > 0 {
		return rawUrl, nil
	}

	newUrl := strings.Join([]string{"http://", rawUrl}, "")

	_, err := url.ParseRequestURI(newUrl)
	if err != nil {
		return "", errors.New("invalid URL")
	}

	return newUrl, nil
}

func ping(rawUrl string, chnl chan healthResult) {
	validUrl, validErr := validateUrl(rawUrl)
	if validErr != nil {
		chnl <- healthResult{
			Status: false,
			Url:    rawUrl,
		}

		return
	}

	resp, err := http.Get(validUrl)
	if err != nil {
		chnl <- healthResult{
			Status: false,
			Url:    rawUrl,
		}

		return
	}
	defer resp.Body.Close()

	chnl <- healthResult{
		Status: true,
		Url:    rawUrl,
	}
}

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("fileUpload")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	defer buffer.Close()

	fileScanner := bufio.NewScanner(buffer)
	fileScanner.Split(bufio.ScanLines)

	healthChan := make(chan healthResult)

	for fileScanner.Scan() {
		urls := strings.Split(fileScanner.Text(), ",")
		for _, url := range urls {
			if trimedUrl := strings.TrimSpace(url); len(trimedUrl) > 0 {
				go ping(trimedUrl, healthChan)
			}
		}
	}

	result := make([]healthResult, 0)
	for h := range healthChan {
		result = append(result, h)
	}

	return c.JSON(result)
}
