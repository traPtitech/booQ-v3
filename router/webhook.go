package router

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	"fmt"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func calcHMACSHA1(message string) string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// PostMessage Webhookでメッセージの投稿
func PostMessage(c echo.Context, message string, isBihin bool) error {
	url := ""
	if isBihin {
		url = "https://q.trap.jp/api/v3/webhooks/" + os.Getenv("TRAQ_BIHIN_WEBHOOK_ID")
	} else {
		url = "https://q.trap.jp/api/v3/webhooks/" + os.Getenv("TRAQ_ITEM_WEBHOOK_ID")
	}
	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	req.Header.Set("X-TRAQ-Signature", calcHMACSHA1(message))

	query := netUrl.Values{}
	query.Add("embed", "1")
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := make([]byte, 512)
	resp.Body.Read(response)

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, response)

	return nil
}
