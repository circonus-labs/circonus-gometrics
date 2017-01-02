// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dashboard API support - Fetch, Create, Update, Delete, and Search
// See: https://login.circonus.com/resources/api/calls/dashboard

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// DashboardGridLayout defines layout
type DashboardGridLayout struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

// DashboardAccessConfig defines access config
type DashboardAccessConfig struct {
	BlackDash           bool   `json:"black_dash"`
	Enabled             bool   `json:"enabled"`
	Fullscreen          bool   `json:"fullscreen"`
	FullscreenHideTitle bool   `json:"fullscreen_hide_title"`
	Nickname            string `json:"nickname"`
	ScaleText           bool   `json:"scale_text"`
	SharedID            string `json:"shared_id"`
	TextSize            int    `json:"text_size"`
}

// DashboardOptions defines options
type DashboardOptions struct {
	AccessConfigs       []DashboardAccessConfig `json:"access_configs"`
	FullscreenHideTitle bool                    `json:"fullscreen_hide_title"`
	HideGrid            bool                    `json:"hide_grid"`
	Linkages            [][]string              `json:"linkages"`
	ScaleText           bool                    `json:"scale_text"`
	TextSize            int                     `json:"text_size"`
}

// ChartTextWidgetDatapoint defines datapoints for charts
type ChartTextWidgetDatapoint struct {
	NumericOnly bool   `json:"numeric_only,omitempty"` // metric cluster
	ClusterID   int    `json:"cluster_id,omitempty"`   // metric cluster
	AccountID   string `json:"account_id,omitempty"`   // metric cluster, metric
	Label       string `json:"label,omitempty"`        // metric
	Metric      string `json:"metric,omitempty"`       // metric
}

// ChartWidgetDefinitionLegend defines chart widget definition legend
type ChartWidgetDefinitionLegend struct {
	Show bool   `json:"show,omitempty"`
	Type string `json:"type,omitempty"`
}

// ChartWidgetWedgeLabels defines chart widget wedge labels
type ChartWidgetWedgeLabels struct {
	OnChart  bool `json:"on_chart,omitempty"`
	ToolTips bool `json:"tooltips,omitempty"`
}

// ChartWidgetWedgeValues defines chart widget wedge values
type ChartWidgetWedgeValues struct {
	Angle string `json:"angle,omitempty"`
	Color string `json:"color,omitempty"`
	Show  bool   `json:"show,omitempty"`
}

// ChartWidgtDefinition defines chart widget definition
type ChartWidgtDefinition struct {
	Datasource        string                      `json:"datasource,omitempty"`
	Derive            string                      `json:"derive,omitempty"`
	DisableAutoformat bool                        `json:"disable_autoformat,omitempty"`
	Formula           string                      `json:"formula,omitempty"`
	Legend            ChartWidgetDefinitionLegend `json:"legend,omitempty"`
	Period            string                      `json:"period,omitempty"`
	PopOnHover        bool                        `json:"pop_onhover,omitempty"`
	WedgeLabels       ChartWidgetWedgeLabels      `json:"wedge_labels,omitempty"`
	WedgeValues       ChartWidgetWedgeValues      `json:"wedge_values,omitempty"`
}

// ForecastGaugeWidgetThresholds defines forecast widget thresholds
type ForecastGaugeWidgetThresholds struct {
	Colors []string `json:"colors,omitempty"` // forecasts, gauges
	Values []string `json:"values,omitempty"` // forecasts, gauges
	Flip   bool     `json:"flip,omitempty"`   // gauges
}

// StatusWidgetAgentStatusSettings defines agent status settings
type StatusWidgetAgentStatusSettings struct {
	Search         string `json:"search,omitempty"`
	ShowAgentTypes string `json:"show_agent_types,omitempty"`
	ShowContact    bool   `json:"show_contact,omitempty"`
	ShowFeeds      bool   `json:"show_feeds,omitempty"`
	ShowSetup      bool   `json:"show_setup,omitempty"`
	ShowSkew       bool   `json:"show_skew,omitempty"`
	ShowUpdates    bool   `json:"show_updates,omitempty"`
}

// StatusWidgetHostStatusSettings defines host status settings
type StatusWidgetHostStatusSettings struct {
	LayoutStyle  string   `json:"layout_style,omitempty"`
	Search       string   `json:"search,omitempty"`
	SortBy       string   `json:"sort_by,omitempty"`
	TagFilterSet []string `json:"tag_filter_set,omitempty"`
}

// DashboardWidgetSettings defines settings specific to widget
type DashboardWidgetSettings struct {
	AccountID           int                             `json:"account_id,omitempty"`            // alerts, clusters, gauges, graphs, lists, status
	Acknowledged        string                          `json:"acknowledged,omitempty"`          // alerts
	Cleared             string                          `json:"cleared,omitempty"`               // alerts
	ContactGroups       []int                           `json:"contact_groups,omitempty"`        // alerts
	Dependents          string                          `json:"dependents,omitempty"`            // alerts
	Display             string                          `json:"display,omitempty"`               // alerts
	Maintenance         string                          `json:"maintenance,omitempty"`           // alerts
	MinAge              string                          `json:"min_age,omitempty"`               // alerts
	OffHours            []int                           `json:"off_hours,omitempty"`             // alerts
	Search              string                          `json:"search,omitempty"`                // alerts, lists
	Severity            string                          `json:"severity,omitempty"`              // alerts
	TagFilterSet        []string                        `json:"tag_filter_set,omitempty"`        // alerts
	TimeWindow          string                          `json:"time_window,omitempty"`           // alerts
	Title               string                          `json:"title,omitempty"`                 // alerts, charts, forecasts, gauges, html
	WeekDays            []string                        `json:"weekdays,omitempty"`              // alerts
	ChartType           string                          `json:"chart_type,omitempty"`            // charts
	Datapoints          []ChartTextWidgetDatapoint      `json:"datapoints,omitempty"`            // charts, text
	Definition          ChartWidgtDefinition            `json:"definition,omitempty"`            // charts
	Algorithm           string                          `json:"algorithm,omitempty"`             // clusters
	ClusterID           int                             `json:"cluster_id,omitempty"`            // clusters
	ClusterName         string                          `json:"cluster_name,omitempty"`          // clusters
	Layout              string                          `json:"layout,omitempty"`                // clusters
	Size                string                          `json:"size,omitempty"`                  // clusters
	Threshold           string                          `json:"threshold,omitempty"`             // clusters
	Format              string                          `json:"format,omitempty"`                // forecasts
	ResourceLimit       string                          `json:"resource_limit,omitempty"`        // forecasts
	ResourceUsage       string                          `json:"resource_usage,omitempty"`        // forecasts
	Thresholds          ForecastGaugeWidgetThresholds   `json:"thresholds,omitempty"`            // forecasts, gauges
	Trend               string                          `json:"trend,omitempty"`                 // forecasts
	CheckUUID           string                          `json:"check_uuid,omitempty"`            // gauges
	DisableAutoformat   bool                            `json:"disable_autoformat,omitempty"`    // gauges
	Formula             string                          `json:"formula,omitempty"`               // gauges
	MetricDisplayName   string                          `json:"metric_display_name,omitempty"`   // gauges
	MetricName          string                          `json:"metric_name,omitempty"`           // gauges
	Period              string                          `json:"period,omitempty"`                // gauges, graphs, text
	RangeHigh           string                          `json:"range_high,omitempty"`            // gauges
	RangeLow            string                          `json:"range_low,omitempty"`             // gauges
	Type                string                          `json:"type,omitempty"`                  // gauges, lists
	ValueType           string                          `json:"value_type,omitempty"`            // gauges, text
	DateWindow          string                          `json:"date_window,omitempty"`           // graphs
	GraphUUID           string                          `json:"graph_id,omitempty"`              // graphs
	HideXAxis           bool                            `json:"hide_xaxis,omitempty"`            // graphs
	HideYAxis           bool                            `json:"hide_yaxis,omitempty"`            // graphs
	KeyInline           bool                            `json:"key_inline,omitempty"`            // graphs
	KeyLoc              string                          `json:"key_loc,omitempty"`               // graphs
	KeySize             string                          `json:"key_size,omitempty"`              // graphs
	KeyWrap             bool                            `json:"key_wrap,omitempty"`              // graphs
	Label               string                          `json:"label,omitempty"`                 // graphs
	OverlaySetID        string                          `json:"overlay_set_id,omitempty"`        // graphs
	Realtime            bool                            `json:"realtime,omitempty"`              // graphs
	ShowFlags           bool                            `json:"show_flags,omitempty"`            // graphs
	Markup              string                          `json:"markup,omitempty"`                // html
	Limit               string                          `json:"limit,omitempty"`                 // lists
	AgentStatusSettings StatusWidgetAgentStatusSettings `json:"agent_status_settings,omitempty"` // status
	ContentType         string                          `json:"content_type,omitempty"`          // status
	HostStatusSettings  StatusWidgetHostStatusSettings  `json:"host_status_settings,omitempty"`  // status
	Autoformat          bool                            `json:"autoformat,omitempty"`            // text
	BodyFormat          string                          `json:"body_format,omitempty"`           // text
	TitleFormat         string                          `json:"title_format,omitempty"`          // text
	UseDefault          bool                            `json:"use_default,omitempty"`           // text
}

// DashboardWidget defines widget
type DashboardWidget struct {
	Active   bool                    `json:"active"`
	Height   int                     `json:"height"`
	Name     string                  `json:"name"`
	Origin   string                  `json:"origin"`
	Settings DashboardWidgetSettings `json:"settings"`
	Type     string                  `json:"type"`
	WidgetID string                  `json:"widget_id"`
	Width    int                     `json:"width"`
}

// Dashboard defines a dashboard
type Dashboard struct {
	CID            string              `json:"_cid,omitempty"`
	Active         bool                `json:"_active,omitempty"`
	Created        int                 `json:"_created,omitempty"`
	CreatedBy      string              `json:"_created_by,omitempty"`
	UUID           string              `json:"_dashboard_uuid,omitempty"`
	LastModified   int                 `json:"_last_modified,omitempty"`
	AccountDefault bool                `json:"account_default"`
	GridLayout     DashboardGridLayout `json:"grid_layout"`
	Options        DashboardOptions    `json:"options"`
	Shared         bool                `json:"shared"`
	Title          string              `json:"title"`
	Widgets        []DashboardWidget   `json:"widgets"`
}

const (
	baseDashboardPath = "/dashboard"
	dashboardCIDRegex = "^" + baseDashboardPath + "/[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{8,12}$"
)

// FetchDashboard retrieves a dashboard definition
func (a *API) FetchDashboard(cid CIDType) (*Dashboard, error) {
	if cid == nil || *cid == "" {
		return nil, fmt.Errorf("Invalid dashboard CID [none]")
	}

	dashboardCID := string(*cid)

	matched, err := regexp.MatchString(dashboardCIDRegex, dashboardCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid dashboard CID [%s]", dashboardCID)
	}

	result, err := a.Get(string(*cid))
	if err != nil {
		return nil, err
	}

	dashboard := new(Dashboard)
	if err := json.Unmarshal(result, dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

// FetchDashboards retrieves all dashboards
func (a *API) FetchDashboards() (*[]Dashboard, error) {
	result, err := a.Get(baseDashboardPath)
	if err != nil {
		return nil, err
	}

	var dashboards []Dashboard
	if err := json.Unmarshal(result, &dashboards); err != nil {
		return nil, err
	}

	return &dashboards, nil
}

// UpdateDashboard update dashboard definition
func (a *API) UpdateDashboard(config *Dashboard) (*Dashboard, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid dashboard config [nil]")
	}

	dashboardCID := string(config.CID)

	matched, err := regexp.MatchString(dashboardCIDRegex, dashboardCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid dashboard CID [%s]", dashboardCID)
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result, err := a.Put(dashboardCID, cfg)
	if err != nil {
		return nil, err
	}

	dashboard := &Dashboard{}
	if err := json.Unmarshal(result, dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

// CreateDashboard create a new dashboard
func (a *API) CreateDashboard(config *Dashboard) (*Dashboard, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid dashboard config [nil]")
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result, err := a.Post(baseDashboardPath, cfg)
	if err != nil {
		return nil, err
	}

	dashboard := &Dashboard{}
	if err := json.Unmarshal(result, dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

// DeleteDashboard delete a dashboard
func (a *API) DeleteDashboard(config *Dashboard) (bool, error) {
	if config == nil {
		return false, fmt.Errorf("Invalid dashboard config [none]")
	}
	cid := CIDType(&config.CID)
	return a.DeleteDashboardByCID(cid)
}

// DeleteDashboardByCID delete a dashboard by cid
func (a *API) DeleteDashboardByCID(cid CIDType) (bool, error) {
	if cid == nil || *cid == "" {
		return false, fmt.Errorf("Invalid dashboard CID [none]")
	}

	dashboardCID := string(*cid)

	matched, err := regexp.MatchString(dashboardCIDRegex, dashboardCID)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, fmt.Errorf("Invalid dashboard CID [%s]", dashboardCID)
	}

	_, err = a.Delete(dashboardCID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// SearchDashboards returns list of dashboards matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) SearchDashboards(searchCriteria *SearchQueryType, filterCriteria *SearchFilterType) (*[]Dashboard, error) {
	q := url.Values{}

	if searchCriteria != nil && *searchCriteria != "" {
		q.Set("search", string(*searchCriteria))
	}

	if filterCriteria != nil && len(*filterCriteria) > 0 {
		for filter, criteria := range *filterCriteria {
			for _, val := range criteria {
				q.Add(filter, val)
			}
		}
	}

	if q.Encode() == "" {
		return a.FetchDashboards()
	}

	reqURL := url.URL{
		Path:     baseDashboardPath,
		RawQuery: q.Encode(),
	}

	result, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var dashboards []Dashboard
	if err := json.Unmarshal(result, &dashboards); err != nil {
		return nil, err
	}

	return &dashboards, nil
}
