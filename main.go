package main

import (
	"errors"
	"fmt"
	"gopkg.in/resty.v1"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var publicIp string
var hostname string

func main() {
	resty.SetTimeout(5 * time.Second)
	resty.SetDebug(true)

	publicIp, _ = getPublicIp()
	hostname, _ = os.Hostname()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Public ip: %s\n", publicIp)
		_, _ = fmt.Fprintf(w, "Hostname: %s\n", hostname)
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}

//Get our public ip from an external service
func getPublicIp() (string, error) {
	urls := [...]string{"ifconfig.me", "icanhazip.com", "ipecho.net/plain", "ifconfig.co"}

	for _, url := range urls {
		resp, err := resty.R().
			SetHeader("User-Agent", "curl/7.58.0").
			Get("https://" + url)
		if err != nil || !resp.IsSuccess() {
			//fmt.Printf("Err %s", err)
			continue
		}
		return strings.TrimSpace(string(resp.Body())), nil

	}

	return "", errors.New("couldn't fetch public ip from any known service")
}
