package endpoint

import "net/url"

func NewEndpoint(scheme, host string) *url.URL {
	return &url.URL{Scheme: scheme, Host: host}
}

func Scheme(scheme string, isSecure bool) string {
	if isSecure {
		return scheme + "s"
	}

	return scheme
}
