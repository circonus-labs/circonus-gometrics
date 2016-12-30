package config

// Key key for CheckBundleConfig
type Key string

// Constants per type as defined in
// https://login.circonus.com/resources/api/calls/check_bundle
const (
	AsyncMetrics     = Key("async_metrics")
	ReverseSecretKey = Key("reverse:secret_key")
	SecretKey        = Key("secret")
	SubmissionURL    = Key("submission_url")

	// "http"
	AuthMethod   = Key("auth_method")
	AuthPassword = Key("auth_password")
	AuthUser     = Key("auth_user")
	Body         = Key("body")
	CAChain      = Key("ca_chain")
	CertFile     = Key("certificate_file")
	Ciphers      = Key("ciphers")
	Code         = Key("code")
	Extract      = Key("extract")
	// HeaderPrefix is special because the actual key is dynamic and matches:
	// `header_(\S+)`
	HeaderPrefix = Key("header_")
	HTTPVersion  = Key("http_version")
	KeyFile      = Key("key_file")
	Method       = Key("method")
	Payload      = Key("payload")
	ReadLimit    = Key("read_limit")
	Redirects    = Key("redirects")
	URL          = Key("url")
)
