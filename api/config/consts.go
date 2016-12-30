package config

import "github.com/circonus-labs/circonus-gometrics/api"

// Constants per type as defined in
// https://login.circonus.com/resources/api/calls/check_bundle
const (
	// "http"
	AuthMethod   api.CheckBundleConfigKey = "auth_method"
	AuthPassword api.CheckBundleConfigKey = "auth_password"
	AuthUser     api.CheckBundleConfigKey = "auth_user"
	Body         api.CheckBundleConfigKey = "body"
	CAChain      api.CheckBundleConfigKey = "ca_chain"
	CertFile     api.CheckBundleConfigKey = "certificate_file"
	Ciphers      api.CheckBundleConfigKey = "ciphers"
	Code         api.CheckBundleConfigKey = "code"
	Extract      api.CheckBundleConfigKey = "extract"
	// HeaderPrefix is special because the actual key is dynamic and matches:
	// `header_(\S+)`
	HeaderPrefix api.CheckBundleConfigKey = "header_"
	HTTPVersion  api.CheckBundleConfigKey = "http_version"
	KeyFile      api.CheckBundleConfigKey = "key_file"
	Method       api.CheckBundleConfigKey = "method"
	Payload      api.CheckBundleConfigKey = "payload"
	ReadLimit    api.CheckBundleConfigKey = "read_limit"
	Redirects    api.CheckBundleConfigKey = "redirects"
	URL          api.CheckBundleConfigKey = "url"
)
