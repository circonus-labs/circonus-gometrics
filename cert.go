package circonusgometrics

import (
	"log"
)

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

func loadCACert() {
	cert := circonusCA

	caDetails := apiCall("/v2/pki/ca.crt")
	val, ok := caDetails["contents"]
	if ok {
		cert = val.([]byte)
	} else {
		log.Printf("Error fetching ca.crt, using default\n")
	}

	rootCA.AppendCertsFromPEM(cert)

	log.Print("CA cert loaded")

}

/*
func parsePEMCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		var der *pem.Block
		der, pemData = pem.Decode(pemData)
		if der == nil {
			break
		}
		if der.Type == "CERTIFICATE" {
			dcerts, err := x509.ParseCertificates(der.Bytes)
			if err != nil {
				return nil, err
			}
			certs = append(certs, dcerts...)
		}
	}
	return certs, nil
}
func getCertChain() {
	caDetails := apiCall("/v2/pki/ca.crt")
	val, ok := caDetails["contents"]
	if !ok {
		log.Printf("Error fetching ca.crt\n")
		setRootCA(circonusCA)
		return
	}

	setRootCA([]byte(val.(string)))
	log.Print("Circonusgometrics fetched CA.")
}

func setRootCA(val []byte) {
	certs, err := parsePEMCertificates(val)
	if err != nil {
		return
	}
	for _, cert := range certs {
		rootCA.AddCert(cert)
	}
}
*/
