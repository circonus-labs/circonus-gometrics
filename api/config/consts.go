package config

// Key for CheckBundleConfig options and CheckDetails info
type Key string

// Constants per type as defined in
// https://login.circonus.com/resources/api/calls/check_bundle
const (
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
	AuthPassword = Key("auth_password")
	AuthUser     = Key("auth_user")
	CAChain      = Key("ca_chain")
	CertFile     = Key("certificate_file")
	Ciphers      = Key("ciphers")
	KeyFile      = Key("key_file")
	Port         = Key("port")
	Query        = Key("query")
	Secret       = Key("secret")
	URI          = Key("uri")
	URL          = Key("url")
	Username     = Key("username")
	UseSSL       = Key("use_ssl")
	HeaderPrefix = Key("header_")
	HTTPVersion  = Key("http_version")
	Method       = Key("method")
	Payload      = Key("payload")
	ReadLimit    = Key("read_limit")

	//
	// CAQL check
	//
	// Common items:
	// Query

	//
	// Circonus Windows Agent
	//
	// Common items:
	// AuthPassword
	// AuthUser
	// Port
	// URL
	Calculated = Key("calculated")
	Category   = Key("cateogry")

	//
	// Cloudwatch
	//
	// Notes:
	// DimPrefix is special because the actual key is dynamic and matches: `dim_(.+)`
	// Common items:
	// URL
	APIKey            = Key("api_key")
	APISecret         = Key("api_secret")
	CloudwatchMetrics = Key("cloudwatch_metrics")
	DimPrefix         = Key("dim_")
	Granularity       = Key("granularity")
	Namespace         = Key("namespace")
	Statistics        = Key("statistics")
	Version           = Key("version")

	//
	// Collectd
	//
	// Common items:
	// AsyncMetrics
	// Username
	// Secret
	SecurityLevel = Key("security_level")

	//
	// Composite
	//
	CompositeMetricName = Key("composite_metric_name")
	Formula             = Key("formula")

	//
	// DHCP
	//
	HardwareAddress = Key("hardware_addr")
	HostIP          = Key("host_ip")
	RequestType     = Key("request_type")
	SendPort        = Key("send_port")

	//
	// DNS
	//
	// Common items:
	// Query
	CType      = Key("ctype")
	Nameserver = Key("nameserver")
	RType      = Key("rtype")

	//
	// EC Console
	//
	// Common items:
	// Port
	Command            = Key("command")
	Objects            = Key("objects")
	SASLAuthentication = Key("sasl_authentication")
	SASLUser           = Key("sasl_user")
	XPath              = Key("xpath")

	//
	// Elastic Search
	//
	// Common items:
	// Port
	// URL

	//
	// Ganglia
	//
	// Common items:
	// AsyncMetrics

	//
	// Google Analytics
	//
	// Common items:
	// Username
	OAuthToken       = Key("oauth_token")
	OAuthTokenSecret = Key("oauth_token_secret")
	OAuthVersion     = Key("oauth_version")
	Password         = Key("password")
	TableID          = Key("table_id")
	UseOAuth         = Key("use_oauth")

	//
	// HA Proxy
	//
	// Common items:
	// AuthPassword
	// AuthUser
	// Port
	// UseSSL
	Host   = Key("host")
	Select = Key("select")

	//
	// HTTP
	//
	// Notes:
	// HeaderPrefix is special because the actual key is dynamic and matches: `header_(\S+)`
	// Common items:
	// AuthPassword
	// AuthUser
	// CAChain
	// CertFile
	// Ciphers
	// KeyFile
	// URL
	// HeaderPrefix
	// HTTPVersion
	// Method
	// Payload
	// ReadLimit
	AuthMethod = Key("auth_method")
	Body       = Key("body")
	Code       = Key("code")
	Extract    = Key("extract")
	Redirects  = Key("redirects")

	//
	// HTTPTRAP
	//
	// Common items:
	// AsyncMetrics
	// Secret

	//
	// IMAP
	//
	// Common items:
	// AuthPassword
	// AuthUser
	// CAChain
	// CertFile
	// Ciphers
	// KeyFile
	// Port
	// UseSSL
	Fetch      = Key("fetch")
	Folder     = Key("folder")
	HeaderHost = Key("header_Host")
	Search     = Key("search")

	//
	// JMX
	//
	// Common items:
	// Password
	// Port
	// URI
	// Username
	MbeanDomains = Key("mbean_domains")

	//
	// JSON
	//
	// Common items:
	// AuthMethod
	// AuthPassword
	// AuthUser
	// CAChain
	// CertFile
	// Ciphers
	// HeaderPrefix
	// HTTPVersion
	// KeyFile
	// Method
	// Payload
	// Port
	// ReadLimit
	// URL

	//
	// Keynote
	//
	// Notes:
	// SlotAliasPrefix is special because the actual key is dynamic and matches: `slot_alias_(\d+)`
	// Common items:
	// APIKey
	BaseURL         = Key("base_url")
	PageComponent   = Key("pagecomponent")
	SlotAliasPrefix = Key("slot_alias_")
	SlotIDList      = Key("slot_id_list")
	TransPageList   = Key("transpagelist")

	//
	// reserved - config option(s) can't actually be set - here for r/o access
	//
	ReverseSecretKey = Key("reverse:secret_key")
	SubmissionURL    = Key("submission_url")

	//
	// Endpoint prefix & cid regex
	//
	DefaultCIDRegex            = "[0-9]+"
	DefaultUUIDRegex           = "[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}"
	AccountPrefix              = "/account"
	AccountCIDRegex            = "^(" + AccountPrefix + "/(" + DefaultCIDRegex + "|current))$"
	AcknowledgementPrefix      = "/acknowledgement"
	AcknowledgementCIDRegex    = "^(" + AcknowledgementPrefix + "/(" + DefaultCIDRegex + "))$"
	AlertPrefix                = "/alert"
	AlertCIDRegex              = "^(" + AlertPrefix + "/(" + DefaultCIDRegex + "))$"
	AnnotationPrefix           = "/annotation"
	AnnotationCIDRegex         = "^(" + AnnotationPrefix + "/(" + DefaultCIDRegex + "))$"
	BrokerPrefix               = "/broker"
	BrokerCIDRegex             = "^(" + BrokerPrefix + "/(" + DefaultCIDRegex + "))$"
	CheckBundleMetricsPrefix   = "/check_bundle_metrics"
	CheckBundleMetricsCIDRegex = "^(" + CheckBundleMetricsPrefix + "/(" + DefaultCIDRegex + "))$"
	CheckBundlePrefix          = "/check_bundle"
	CheckBundleCIDRegex        = "^(" + CheckBundlePrefix + "/(" + DefaultCIDRegex + "))$"
	CheckPrefix                = "/check"
	CheckCIDRegex              = "^(" + CheckPrefix + "/(" + DefaultCIDRegex + "))$"
	ContactGroupPrefix         = "/contact_group"
	ContactGroupCIDRegex       = "^(" + ContactGroupPrefix + "/(" + DefaultCIDRegex + "))$"
	DashboardPrefix            = "/dashboard"
	DashboardCIDRegex          = "^(" + DashboardPrefix + "/(" + DefaultCIDRegex + "))$"
	GraphPrefix                = "/graph"
	GraphCIDRegex              = "^(" + GraphPrefix + "/(" + DefaultUUIDRegex + "))$"
	MaintenancePrefix          = "/maintenance"
	MaintenanceCIDRegex        = "^(" + MaintenancePrefix + "/(" + DefaultCIDRegex + "))$"
	MetricClusterPrefix        = "/metric_cluster"
	MetricClusterCIDRegex      = "^(" + MetricClusterPrefix + "/(" + DefaultCIDRegex + "))$"
	MetricPrefix               = "/metric"
	MetricCIDRegex             = "^(" + MetricPrefix + "/((" + DefaultCIDRegex + ")_([^[:space:]]+)))$"
	OutlierReportPrefix        = "/outlier_report"
	OutlierReportCIDRegex      = "^(" + OutlierReportPrefix + "/(" + DefaultCIDRegex + "))$"
	ProvisionBrokerPrefix      = "/provision_broker"
	ProvisionBrokerCIDRegex    = "^(" + ProvisionBrokerPrefix + "/([a-z0-9]+-[a-z0-9]+))$"
	RuleSetGroupPrefix         = "/rule_set_group"
	RuleSetGroupCIDRegex       = "^(" + RuleSetGroupPrefix + "/(" + DefaultCIDRegex + "))$"
	RuleSetPrefix              = "/rule_set"
	RuleSetCIDRegex            = "^(" + RuleSetPrefix + "/((" + DefaultCIDRegex + ")_([^[:space:]]+)))$"
	UserPrefix                 = "/user"
	UserCIDRegex               = "^(" + UserPrefix + "/(" + DefaultCIDRegex + "|current))$"
	WorksheetPrefix            = "/worksheet"
	WorksheetCIDRegex          = "^(" + WorksheetPrefix + "/(" + DefaultUUIDRegex + "))$"
	// contact group serverity levels
	NumSeverityLevels = 5
)
