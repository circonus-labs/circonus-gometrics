package circonusgometrics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func (m *CirconusMetrics) submit(output map[string]interface{}, newMetrics map[string]*CheckBundleMetric) {
	if len(newMetrics) > 0 {
		m.addNewCheckMetrics(newMetrics)
	}

	str, err := json.Marshal(output)
	if err != nil {
		m.Log.Printf("Error marshling output %+v", err)
		return
	}

	numStats, err := m.trapCall(str)
	if err != nil {
		m.Log.Printf("[ERROR] %+v\n", err)
		return
	}

	m.Log.Printf("%d stats sent to %s\n", numStats, m.trapUrl)
}

func (m *CirconusMetrics) trapCall(payload []byte) (int, error) {
	m.trapmu.Lock()
	trapUrl := m.trapUrl
	m.trapmu.Unlock()

	dataReader := bytes.NewReader(payload)

	req, err := retryablehttp.NewRequest("PUT", trapUrl, dataReader)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")

	client := retryablehttp.NewClient()
	if m.trapSSL {
		client.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs:    m.certPool,
				ServerName: m.trapCN,
			},
			DisableKeepAlives:   true,
			MaxIdleConnsPerHost: -1,
			DisableCompression:  true,
		}
	} else {
		client.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableKeepAlives:   true,
			MaxIdleConnsPerHost: -1,
			DisableCompression:  true,
		}
	}
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 50 * time.Millisecond
	client.RetryMax = 3
	client.Logger = m.Log

	attempts := -1
	client.RequestLogHook = func(logger *log.Logger, req *http.Request, retryNumber int) {
		attempts = retryNumber
	}

	resp, err := client.Do(req)
	if err != nil {
		if attempts == client.RetryMax {
			if time.Since(m.trapLastUpdate) >= m.MaxSubmissionUrlAge {
				m.trapmu.Lock()
				defer m.trapmu.Unlock()
				m.ready = false
				m.trapUrl = ""
			}
		}
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.Log.Printf("[ERROR] reading body, proceeding. %s\n", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		m.Log.Printf("[ERROR] parsing body, proceeding. %s\n", err)
	}

	if resp.StatusCode != 200 {
		return 0, errors.New("bad response code: " + strconv.Itoa(resp.StatusCode))
	}
	switch v := response["stats"].(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
	}
	return 0, errors.New("bad response type")
}
