package controller

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type healthResult struct {
	Status bool   `json:"status"`
	Url    string `json:"url"`
}

func GetHttpClient() http.Client {
	httpClient := http.Client{
		Timeout: 2 * time.Second,
	}

	return httpClient
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

func ping(rawUrl string) bool {
	validUrl, validErr := validateUrl(rawUrl)
	if validErr != nil {
		return false
	}

	req, err := http.NewRequest("HEAD", validUrl, nil)
	if err != nil {
		return false
	}

	client := GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()

	if conditional := resp.StatusCode >= 200 && resp.StatusCode < 300; conditional {
		return true
	}

	return false
}

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")

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

	var wg sync.WaitGroup
	var healthResults []healthResult
	for fileScanner.Scan() {
		wg.Add(1)

		urls := strings.Split(fileScanner.Text(), ",")
		for _, url := range urls {
			if trimedUrl := strings.TrimSpace(url); len(trimedUrl) > 0 {
				go func() {
					defer wg.Done()

					status := ping(trimedUrl)
					healthResults = append(healthResults, healthResult{
						Status: status,
						Url:    trimedUrl,
					})
				}()
			}
		}
	}
	wg.Wait()

	return c.JSON(healthResults)
}
