// Package service section
package service

import (
	"bytes"
	"caos/model"
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
	Input     []string
	Transform model.ChainPrompt
}

// ExecuteChainJob - ExecuteChainJob chain sequence
func (c *Chain) ExecuteChainJob(service Agent, prompt *model.PromptProperties) {
	defer clear()
	c.Input = append(c.Input, prompt.Input...)
	c.onConstructAssemble(service, prompt.Input)
	if strings.Contains(service.EngineProperties.Model, "16k") {
		c.onConstructValidation(service)
	}
}

// clear - Delete existent cookies and ssl certificate
func clear() {
	os.Remove("rootCA.pem")
	os.Remove(".cookies")
}

// setCertificateSSL - Implement SSL
func setCertificateSSL(service Agent, url string) {
	os.Remove("rootCA.pem")
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
	os.Remove(".cookies")
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
	var bufferBody string

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

				if isLinkReference(content) {
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

				if isTagReference(tokenType, content) &&
					len(resultBody) <= 99 {
					context := fmt.Sprintf("%v", content)
					bufferBody += fmt.Sprintf(" %v", context)
				}
			}
		} else {
			resultBody = append(resultBody, bufferBody)
			return urlHeader, resultBody
		}
	}
}

// isTagReference
func isTagReference(tknType html.TokenType, c string) bool {
	if tknType == html.TextToken {
		compound := strings.Split(c, " ")
		if !strings.ContainsAny(c, "</>{}\n") &&
			len(compound) >= 1 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// isLinkReference
func isLinkReference(c string) bool {
	if strings.Contains(c, "href") &&
		!strings.Contains(c, "analytics") &&
		!strings.Contains(c, "google") &&
		!strings.Contains(c, "cdn") {
		return true
	} else {
		return false
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

	_, e := os.OpenFile("rootCA.pem", 0, 0)
	if e == nil {
		easy.Setopt(curl.OPT_USE_SSL, true)
		easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
		easy.Setopt(curl.OPT_SSL_VERIFYHOST, true)
		easy.Setopt(curl.OPT_SSL_VERIFYSTATUS, false)
		easy.Setopt(curl.OPT_CAINFO, "rootCA.pem")
	}
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
func (c *Chain) onConstructAssemble(service Agent, input []string) {
	setCertificateSSL(service, parameters.ExternalSearchBaseURL)
	setCookieJar(service, parameters.ExternalSearchBaseURL)
	context := strings.ReplaceAll(input[0], " ", "+")
	req := fmt.Sprint(parameters.ExternalSearchBaseURL, context)
	reader := bytes.NewReader(setOpt(req, service.preferences.User, service.preferences.Encoding))
	c.Transform.Source, c.Transform.Context = setConstructResults(reader)
}

// onConstructValidation - Validate transformer
func (c *Chain) onConstructValidation(service Agent) {
	var sourceBuffer []string
	var contextBuffer []string
	for i := range c.Transform.Source {
		if i <= 3 {
			setCertificateSSL(service, c.Transform.Source[i])
			setCookieJar(service, c.Transform.Source[i])

			req := fmt.Sprint(c.Transform.Source[i])
			reader := bytes.NewReader(setOpt(req, service.preferences.User, service.preferences.Encoding))
			header, context := setConstructResults(reader)

			sourceBuffer = append(sourceBuffer, header...)
			contextBuffer = append(contextBuffer, context...)
		}
	}

	c.Transform.Source = append(c.Transform.Source, sourceBuffer...)
	c.Transform.Context = append(c.Transform.Context, contextBuffer...)
}
