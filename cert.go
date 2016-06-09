package circonusgometrics

import (
	"encoding/json"
	"fmt"
)

type CACert struct {
	Contents string `json:"contents"`
}

// Load the CA cert for the broker designated by the submission url
func (m *CirconusMetrics) loadCACert() {
	if m.cert != nil {
		return
	}

	if !m.trapSSL {
		// ssl not needed
		return
	}

	cert, err := m.fetchCert()
	if err != nil {
		if m.Debug {
			m.Log.Printf("[DEBUG] Unable to fetch ca.crt, using default. %+v\n", err)
		}
	}

	if cert == nil {
		cert = circonusCA
	}

	m.cert = cert
	m.certPool.AppendCertsFromPEM(cert)
}

// Use the Circonus API to retrieve the CA certificate
func (m *CirconusMetrics) fetchCert() ([]byte, error) {
	response, err := m.apiCall("GET", "/v2/pki/ca.crt", nil)
	if err != nil {
		return nil, err
	}

	cadata := new(CACert)
	err = json.Unmarshal(response, cadata)
	if err != nil {
		return nil, err
	}

	if cadata.Contents == "" {
		return nil, fmt.Errorf("[ERROR] Unable to find ca cert %+v", cadata)
	}

	return []byte(cadata.Contents), nil
}

// Default Circonus CA certificate
var circonusCA []byte = []byte(`-----BEGIN CERTIFICATE-----
MIID4zCCA0ygAwIBAgIJAMelf8skwVWPMA0GCSqGSIb3DQEBBQUAMIGoMQswCQYD
VQQGEwJVUzERMA8GA1UECBMITWFyeWxhbmQxETAPBgNVBAcTCENvbHVtYmlhMRcw
FQYDVQQKEw5DaXJjb251cywgSW5jLjERMA8GA1UECxMIQ2lyY29udXMxJzAlBgNV
BAMTHkNpcmNvbnVzIENlcnRpZmljYXRlIEF1dGhvcml0eTEeMBwGCSqGSIb3DQEJ
ARYPY2FAY2lyY29udXMubmV0MB4XDTA5MTIyMzE5MTcwNloXDTE5MTIyMTE5MTcw
NlowgagxCzAJBgNVBAYTAlVTMREwDwYDVQQIEwhNYXJ5bGFuZDERMA8GA1UEBxMI
Q29sdW1iaWExFzAVBgNVBAoTDkNpcmNvbnVzLCBJbmMuMREwDwYDVQQLEwhDaXJj
b251czEnMCUGA1UEAxMeQ2lyY29udXMgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MR4w
HAYJKoZIhvcNAQkBFg9jYUBjaXJjb251cy5uZXQwgZ8wDQYJKoZIhvcNAQEBBQAD
gY0AMIGJAoGBAKz2X0/0vJJ4ad1roehFyxUXHdkjJA9msEKwT2ojummdUB3kK5z6
PDzDL9/c65eFYWqrQWVWZSLQK1D+v9xJThCe93v6QkSJa7GZkCq9dxClXVtBmZH3
hNIZZKVC6JMA9dpRjBmlFgNuIdN7q5aJsv8VZHH+QrAyr9aQmhDJAmk1AgMBAAGj
ggERMIIBDTAdBgNVHQ4EFgQUyNTsgZHSkhhDJ5i+6IFlPzKYxsUwgd0GA1UdIwSB
1TCB0oAUyNTsgZHSkhhDJ5i+6IFlPzKYxsWhga6kgaswgagxCzAJBgNVBAYTAlVT
MREwDwYDVQQIEwhNYXJ5bGFuZDERMA8GA1UEBxMIQ29sdW1iaWExFzAVBgNVBAoT
DkNpcmNvbnVzLCBJbmMuMREwDwYDVQQLEwhDaXJjb251czEnMCUGA1UEAxMeQ2ly
Y29udXMgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MR4wHAYJKoZIhvcNAQkBFg9jYUBj
aXJjb251cy5uZXSCCQDHpX/LJMFVjzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEB
BQUAA4GBAAHBtl15BwbSyq0dMEBpEdQYhHianU/rvOMe57digBmox7ZkPEbB/baE
sYJysziA2raOtRxVRtcxuZSMij2RiJDsLxzIp1H60Xhr8lmf7qF6Y+sZl7V36KZb
n2ezaOoRtsQl9dhqEMe8zgL76p9YZ5E69Al0mgiifTteyNjjMuIW
-----END CERTIFICATE-----`)
