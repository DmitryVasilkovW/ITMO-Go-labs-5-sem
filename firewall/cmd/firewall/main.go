//go:build !solution

package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	confPath     = flag.String("conf", "./configs/example.yaml", "path to config file")
	firewallAddr = flag.String("addr", ":8081", "address of firewall")
	serviceAddr  = flag.String("service-addr", "http://eu.httpbin.org/", "address of protected service")
)

var _ http.RoundTripper = (*firewall)(nil)

type firewall struct {
	config *Rules
}

func NewFirewall(config *Rules) http.RoundTripper {
	f := &firewall{
		config: config,
	}
	return f
}

func main() {
	flag.Parse()

	rules, parsedURL := getParams()
	reverseProxy := httputil.NewSingleHostReverseProxy(parsedURL)
	reverseProxy.Transport = NewFirewall(rules)

	http.HandleFunc("/", reverseProxy.ServeHTTP)
	err := http.ListenAndServe(*firewallAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getParams() (*Rules, *url.URL) {
	confData, err := os.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	rules, err := ParseR(confData)
	if err != nil {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse(*serviceAddr)
	if err != nil {
		log.Fatal(err)
	}

	return rules, parsedURL
}

func (f *firewall) RoundTrip(request *http.Request) (*http.Response, error) {
	rule := f.getRule(request.RequestURI)

	if !f.isRequestAllowed(rule, request) {
		return f.forbiddenResponse(), nil
	}

	res, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if !f.isResponseAllowed(rule, res) {
		return f.forbiddenResponse(), nil
	}

	return res, nil
}

func (f *firewall) isRequestAllowed(rule *Rule, request *http.Request) bool {
	return rule.reqAllowed(request)
}

func (f *firewall) isResponseAllowed(rule *Rule, res *http.Response) bool {
	return rule.resAllowed(res)
}

func (f *firewall) forbiddenResponse() *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBufferString("Forbidden")),
		StatusCode: 403,
	}
}

func (f *firewall) getRule(reqURI string) *Rule {
	for _, rule := range f.config.Rules {
		if rule.Endpoint == reqURI {
			return &rule
		}
	}
	return nil
}

func (r *Rule) reqAllowed(req *http.Request) bool {
	if r == nil {
		return true
	}

	if !r.isUserAgentAllowed(req.UserAgent()) {
		return false
	}

	if !r.hasRequiredHeaders(req) {
		return false
	}

	if !r.hasForbiddenHeaders(req) {
		return false
	}

	if req.Body != nil {
		if !r.isRequestBodyAllowed(req) {
			return false
		}
	}

	return true
}

func (r *Rule) isUserAgentAllowed(userAgent string) bool {
	for _, v := range r.ForbiddenUserAgents {
		matchString, err := regexp.MatchString(v, userAgent)
		if matchString || err != nil {
			return false
		}
	}
	return true
}

func (r *Rule) hasRequiredHeaders(req *http.Request) bool {
	for _, v := range r.RequiredHeaders {
		if req.Header.Get(v) == "" {
			return false
		}
	}
	return true
}

func (r *Rule) hasForbiddenHeaders(req *http.Request) bool {
	for _, v := range r.ForbiddenHeaders {
		fhPair := strings.SplitN(v, ": ", 2)
		fieldName, fieldRegex := fhPair[0], fhPair[1]
		matchString, err := regexp.MatchString(fieldRegex, req.Header.Get(fieldName))
		if matchString || err != nil {
			return false
		}
	}
	return true
}

func (r *Rule) isRequestBodyAllowed(req *http.Request) bool {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return false
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	return !(BEL(body, r.MaxRequestLengthBytes) || BForbidden(string(body), r.ForbiddenRequestRe))
}

func (r *Rule) resAllowed(response *http.Response) bool {
	if r == nil {
		return true
	}

	if !r.isResponseCodeAllowed(response.StatusCode) {
		return false
	}

	if response.Body != nil {
		if !r.isResponseBodyAllowed(response) {
			return false
		}
	}

	return true
}

func (r *Rule) isResponseCodeAllowed(statusCode int) bool {
	for _, v := range r.ForbiddenResponseCodes {
		if statusCode == v {
			return false
		}
	}
	return true
}

func (r *Rule) isResponseBodyAllowed(response *http.Response) bool {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false
	}

	response.Body = io.NopCloser(bytes.NewBuffer(body))

	return !(BEL(body, r.MaxResponseLengthBytes) || BForbidden(string(body), r.ForbiddenResponseRe))
}

func BEL(body []byte, limit int) bool {
	if limit <= 0 {
		return false
	}

	return len(body) > limit
}

func BForbidden(body string, ForbiddenRes []string) bool {
	if len(body) <= 0 {
		return false
	}
	for _, v := range ForbiddenRes {
		matchString, err := regexp.MatchString(v, body)
		if matchString || err != nil {
			return true
		}
	}

	return false
}
