// Package service section
package service

import (
	"bytes"
	"caos/service/parameters"
	"caos/util"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/andelf/go-curl"
	"golang.org/x/net/html"
)

// Chain - Event sequence object
type Chain struct {
	input     []string
	transform *Transformer
}

// ExecuteChainJob - ExecuteChainJob chain sequence
func (c *Chain) ExecuteChainJob(service Agent) {
	c.transform.onConstructAssemble(service, c.input)
	go c.transform.onConstructValidation(service)
}

// Transformer - LLM properties
type Transformer struct {
	source           []string
	ctxSource        []string
	contextualPrompt []string
	validationPrompt []string
}

// setCertificateSSL - Implement SSL
func setCertificateSSL(service Agent, url string) {
	trimL := strings.Split(url, "//")
	if trimL == nil {
		return
	}

	trimR := strings.Split(trimL[1], "/")
	if trimR == nil {
		return
	}

	domain := trimR[0] + ":443"
	conn, _ := tls.Dial("tcp", domain, &tls.Config{})
	conn.Handshake()

	defer conn.Close()
	state := conn.ConnectionState()

	var end string
	var buffer []string
	for i := range state.PeerCertificates {
		publicKey := pem.EncodeToMemory(&pem.Block{Bytes: state.PeerCertificates[i].Raw})
		filter := string(publicKey)
		header := strings.Replace(filter, "BEGIN ", "BEGIN CERTIFICATE", 1)
		end = strings.Replace(header, "END ", "END CERTIFICATE", 1)
		buffer = append(buffer, fmt.Sprint("\n", end))
	}

	os.WriteFile("rootCA.pem", []byte(util.RemoveWrapper(fmt.Sprint(buffer))), 0644)
}

// setCookiejar - Implement cookies
func setCookieJar(service Agent, url string) {
	jar, _ := cookiejar.New(nil)
	service.exClient.Jar = jar

	resp, err := service.exClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var fields string
	for _, c := range service.exClient.Jar.Cookies(resp.Request.URL) {
		fields += fmt.Sprint(c.Name, "=", c.Value, "\n")
	}

	os.WriteFile(".cookies", []byte(fields), 0644)
}

// setConstructResults - Implement transformer
func setConstructResults(reader *bytes.Reader) ([]string, []string) {
	resultBody := []string{}
	urlHeader := []string{}

	node := html.NewTokenizer(reader)
	for {
		i := node.Next()
		if i != html.ErrorToken {
			token := node.Token()
			tokenType := token.Type
			content := token.String()
			if content != " " &&
				!strings.Contains(content, "(function(){") {

				isLinkRef := func(c string) bool {
					if strings.Contains(c, "href") &&
						!strings.Contains(c, "analytics") &&
						!strings.Contains(c, "google") &&
						!strings.Contains(c, "cdn") {
						return true
					} else {
						return false
					}
				}

				if isLinkRef(content) {
					urls := strings.Split(content, "https://")
					o := len(urls)
					if o > 1 {
						var nested []string
						if strings.ContainsAny(urls[1], "%") {
							nested = strings.Split(urls[1], "%")
						} else {
							nested = strings.Split(urls[1], "&amp")
						}
						remote := fmt.Sprintf("https://%v", nested[0])
						shorted := strings.Split(remote, " ")

						if !strings.Contains(fmt.Sprint(urlHeader), shorted[0]) {
							urlHeader = append(urlHeader, shorted[0])
						}
					}
				}

				isTagRef := func(tknType html.TokenType, c string) bool {
					if tknType == html.TextToken {
						compound := strings.Split(c, " ")
						if !strings.ContainsAny(c, "</>{}\n") &&
							len(compound) > 3 &&
							len(c) >= 1 {
							return true
						} else {
							return false
						}
					} else {
						return false
					}
				}

				if isTagRef(tokenType, content) &&
					len(resultBody) <= 99 {
					context := fmt.Sprintf("%v", content)
					resultBody = append(resultBody, context)
				}
			}
		} else {
			return urlHeader, resultBody
		}
	}
}

// setOpt
func setOpt(req string, user string, enc string) []byte {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, req)
	easy.Setopt(curl.OPT_USERAGENT, user)
	easy.Setopt(curl.OPT_ACCEPT_ENCODING, enc)
	easy.Setopt(curl.OPT_HTTP_CONTENT_DECODING, true)
	easy.Setopt(curl.OPT_COOKIESESSION, true)
	easy.Setopt(curl.OPT_COOKIEFILE, ".cookies")
	easy.Setopt(curl.OPT_USE_SSL, true)
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	easy.Setopt(curl.OPT_SSL_VERIFYHOST, true)
	easy.Setopt(curl.OPT_SSL_VERIFYSTATUS, false)
	easy.Setopt(curl.OPT_CAINFO, "rootCA.pem")
	// Create a buffer for storing the response body
	var buffer []byte
	writer := func(data []byte, userdata interface{}) bool {
		buffer = append(buffer, data...)
		return true
	}

	// Set write function to store response body in buffer
	easy.Setopt(curl.OPT_WRITEFUNCTION, writer)

	// Perform the request and check for errors
	_ = easy.Perform()
	return buffer
}

// onConstructAssemble - Assemble transformer
func (c *Transformer) onConstructAssemble(service Agent, input []string) {
	setCertificateSSL(service, parameters.ExternalSearchBaseURL)
	setCookieJar(service, parameters.ExternalSearchBaseURL)
	context := strings.ReplaceAll(input[0], " ", "+")
	req := fmt.Sprint(parameters.ExternalSearchBaseURL, context)
	reader := bytes.NewReader(setOpt(req, service.preferences.User, service.preferences.Encoding))
	c.source, c.contextualPrompt = setConstructResults(reader)
}

// onConstructValidation - Validate transformer
func (c *Transformer) onConstructValidation(service Agent) {
	var sourceBuffer []string
	var contextBuffer []string
	for i := range c.source {
		if i <= 10 {
			setCertificateSSL(service, c.source[i])
			setCookieJar(service, c.source[i])

			req := fmt.Sprint(c.source[i])
			reader := bytes.NewReader(setOpt(req, service.preferences.User, service.preferences.Encoding))
			header, context := setConstructResults(reader)

			sourceBuffer = append(sourceBuffer, header...)
			contextBuffer = append(contextBuffer, context...)
		}
	}

	c.ctxSource = append(c.ctxSource, sourceBuffer...)
	c.validationPrompt = append(c.validationPrompt, contextBuffer...)
}
