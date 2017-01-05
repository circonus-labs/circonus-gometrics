package config

// Key for CheckBundleConfig options
type Key string

// Constants per type as defined in
// https://login.circonus.com/resources/api/calls/check_bundle
const (
	// Endpoint prefixes
	AccountPrefix              = "/account"
	AccountCIDRegex            = "^" + AccountPrefix + "/([0-9]+|current)$"
	AcknowledgementPrefix      = "/acknowledgement"
	AcknowledgementCIDRegex    = "^" + AcknowledgementPrefix + "/[0-9]+$"
	AlertPrefix                = "/alert"
	AlertCIDRegex              = "^" + AlertPrefix + "/[0-9]+$"
	AnnotationPrefix           = "/annotation"
	AnnotationCIDRegex         = "^" + AnnotationPrefix + "/[0-9]+$"
	BrokerPrefix               = "/broker"
	BrokerCIDRegex             = "^" + BrokerPrefix + "/[0-9]+$"
	CheckBundleMetricsPrefix   = "/check_bundle_metrics"
	CheckBundleMetricsCIDRegex = "^" + CheckBundleMetricsPrefix + "/[0-9]+$"
	CheckBundlePrefix          = "/check_bundle"
	CheckBundleCIDRegex        = "^" + CheckBundlePrefix + "/[0-9]+$"
	CheckPrefix                = "/check"
	CheckCIDRegex              = "^" + CheckPrefix + "/[0-9]+$"
	ContactGroupPrefix         = "/contact_group"
	ContactGroupCIDRegex       = "^" + ContactGroupPrefix + "/[0-9]+$"
	DashboardPrefix            = "/dashboard"
	DashboardCIDRegex          = "^" + DashboardPrefix + "/[0-9]+$"
	GraphPrefix                = "/graph"
	GraphCIDRegex              = "^" + GraphPrefix + "/[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$"
	MaintenancePrefix          = "/maintenance"
	MaintenanceCIDRegex        = "^" + MaintenancePrefix + "/[0-9]+$"
	MetricClusterPrefix        = "/metric_cluster"
	MetricClusterCIDRegex      = "^" + MetricClusterPrefix + "/[0-9]+$"
	MetricPrefix               = "/metric"
	MetricCIDRegex             = "^" + MetricPrefix + "/[0-9]+_[^[:space:]]+$"
	OutlierReportPrefix        = "/outlier_report"
	OutlierReportCIDRegex      = "^" + OutlierReportPrefix + "/[0-9]+$"
	ProvisionBrokerPrefix      = "/provision_broker"
	ProvisionBrokerCIDRegex    = "^" + ProvisionBrokerPrefix + "/[a-z0-9]+-[a-z0-9]+$"
	RuleSetGroupPrefix         = "/rule_set_group"
	RulesetGroupCIDRegex       = "^" + RuleSetGroupPrefix + "/[0-9]+$"
	RuleSetPrefix              = "/rule_set"
	RuleSetCIDRegex            = "^" + RuleSetPrefix + "/[0-9]+_.+$"
	UserPrefix                 = "/user"
	UserCIDRegex               = "^" + UserPrefix + "/([0-9]+|current)$"
	WorksheetPrefix            = "/worksheet"
	WorksheetCIDRegex          = "^" + WorksheetPrefix + "/[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$"

	NumSeverityLevels = 5

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
