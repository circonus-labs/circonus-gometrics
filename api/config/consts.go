package config

// Key for CheckBundleConfig options
type Key string

// Constants per type as defined in
// https://login.circonus.com/resources/api/calls/check_bundle
const (
	// Endpoint prefixes
	AccountPrefix            = "/account"
	AcknowledgementPrefix    = "/acknowledgement"
	AlertPrefix              = "/alert"
	AnnotationPrefix         = "/annotation"
	BrokerPrefix             = "/broker"
	CheckBundlePrefix        = "/check_bundle"
	CheckBundleMetricsPrefix = "/check_bundle_metrics"
	CheckPrefix              = "/check"
	ContactGroupPrefix       = "/contact_group"
	DashboardPrefix          = "/dashboard"
	GraphPrefix              = "/graph"
	MaintenancePrefix        = "/maintenance"
	MetricClusterPrefix      = "/metric_cluster"
	MetricPrefix             = "/metric"
	OutlierReportPrefix      = "/outlier_report"
	ProvisionBrokerPrefix    = "/provision_broker"
	RuleSetGroupPrefix       = "/rule_set_group"
	RuleSetPrefix            = "/rule_set"
	UserPrefix               = "/user"
	WorksheetPrefix          = "/worksheet"

	//
	// default settings for api.NewCheckBundle()
	//
	DefaultCheckBundleMetricLimit = -1 // unlimited
	DefaultCheckBundleStatus      = "active"
	DefaultCheckBundlePeriod      = 60
	DefaultCheckBundleTimeout     = 10

	//
	// common (apply to more than one check type)
	//
	AsyncMetrics = Key("async_metrics")

	//
	// httptrap
	//
	SecretKey = Key("secret")

	//
	// "http"
	//
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

	//
	// reserved - config option(s) can't actually be set - here for r/o access
	//
	ReverseSecretKey = Key("reverse:secret_key")
	SubmissionURL    = Key("submission_url")
)
