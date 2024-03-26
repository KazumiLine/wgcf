package wgcf

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func trace() error {
	response, err := http.Get("https://cloudflare.com/cdn-cgi/trace")
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Println("Trace result:")
	fmt.Println(strings.TrimSpace(string(bodyBytes)))
	return nil
}
