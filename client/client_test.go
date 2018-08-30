package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dhaifley/dlib/dauth"
)

// The certificate and key data in this file was created for testing purposes
// only. They are not real or used in any actual projects.

const testCrt = `
-----BEGIN CERTIFICATE-----
MIIFNTCCBB2gAwIBAgIIUdhUn2xoanEwDQYJKoZIhvcNAQELBQAwgbQxCzAJBgNV
BAYTAlVTMRAwDgYDVQQIEwdBcml6b25hMRMwEQYDVQQHEwpTY290dHNkYWxlMRow
GAYDVQQKExFHb0RhZGR5LmNvbSwgSW5jLjEtMCsGA1UECxMkaHR0cDovL2NlcnRz
LmdvZGFkZHkuY29tL3JlcG9zaXRvcnkvMTMwMQYDVQQDEypHbyBEYWRkeSBTZWN1
cmUgQ2VydGlmaWNhdGUgQXV0aG9yaXR5IC0gRzIwHhcNMTcwNTA1MTU1MzAwWhcN
MjAwNTA1MTU1MzAwWjA+MSEwHwYDVQQLExhEb21haW4gQ29udHJvbCBWYWxpZGF0
ZWQxGTAXBgNVBAMMECoucm95YWxmYXJtcy5jb20wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQCs9UN/tRNf0LKwgvSqbvvZU31K4H/VOq6sywP2cCM1+ZZs
r87tjm3+HOJWnXKJjgnpXrC29pM1rlyx/EwLPoEMy0BooApHlmTq4VO+vrnTLkj1
+27a4rAjpT3kzwTul3mVjdGVf+xJfzPOkX0iu1+wxEherPDfCBCqS5FCJBFTvggk
bHADL2S4ScX6rBIU4/YlPM+55pWFVIuuX5Va8iFpHbtB75ZiQ1vF7cJ7T4x+dwze
aPCrsTqQLlS4x06NvISP1Wn+jCy1Ii3vGua4u1o//KrkWquaQfXftTEz+jRaP4Qq
Tiv2JSdIyanLc/8bWFtAWY6d5torWM40HrKbjF4rAgMBAAGjggG+MIIBujAMBgNV
HRMBAf8EAjAAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAOBgNVHQ8B
Af8EBAMCBaAwNwYDVR0fBDAwLjAsoCqgKIYmaHR0cDovL2NybC5nb2RhZGR5LmNv
bS9nZGlnMnMxLTUwMi5jcmwwXQYDVR0gBFYwVDBIBgtghkgBhv1tAQcXATA5MDcG
CCsGAQUFBwIBFitodHRwOi8vY2VydGlmaWNhdGVzLmdvZGFkZHkuY29tL3JlcG9z
aXRvcnkvMAgGBmeBDAECATB2BggrBgEFBQcBAQRqMGgwJAYIKwYBBQUHMAGGGGh0
dHA6Ly9vY3NwLmdvZGFkZHkuY29tLzBABggrBgEFBQcwAoY0aHR0cDovL2NlcnRp
ZmljYXRlcy5nb2RhZGR5LmNvbS9yZXBvc2l0b3J5L2dkaWcyLmNydDAfBgNVHSME
GDAWgBRAwr0njsw0gzCiM9f7bLPwtCyAzjArBgNVHREEJDAighAqLnJveWFsZmFy
bXMuY29tgg5yb3lhbGZhcm1zLmNvbTAdBgNVHQ4EFgQUWN1+hI1gcguUsEvDQDrS
l3mTyCowDQYJKoZIhvcNAQELBQADggEBAEyAI8qQwmIZn1eeeuwmqXGbhK7/gXFh
XEzHsEUXrMi0t/wFcAG5gDgn+YZK9AD/lJkLLMnlcIjAzkqtwXbqdV3Nm/yNfq9C
hc2IzqHQ/7L3vGY6DTVHjEdOWGTJkGt0hqbF+LO+tM+wBx+G57E7EzGGTR7wTSQE
oRQLmesKFKbdrX22UANpULKlf1yfun/wz5Zear0RAIVfb/xBR/Bfhto6R38dv4Da
4QCAZTPQyC2xPzgYpfZJuzHIdMw6/vRPIqxUEHSHVsoE0MeTHS49VhEfAW8/p5X9
2CKOcjBYk9xQ2kKCYUZtpGyJJ318qiYckG+eu2dLvUo61OvboABTgiY=
-----END CERTIFICATE-----`

func TestRESTClientLoginLogout(t *testing.T) {
	fs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"id":1,"user":"test"}`)
	}))

	defer fs.Close()
	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(fs.URL)
		},
	}

	hc := &http.Client{Transport: tr}
	rc, err := NewRESTClient(fs.URL, fs.URL, testCrt)
	if err != nil {
		t.Fatal(err)
	}

	rc.Client = hc
	ch := rc.Login(&dauth.User{User: "test", Pass: "test"})
	for res := range ch {
		if res.Err != nil {
			t.Fatal(res.Err)
		}
	}

	ch = rc.Logout()
	for res := range ch {
		if res.Err != nil {
			t.Fatal(res.Err)
		}
	}
}

func TestRPCClientLoginLogout(t *testing.T) {
	rpc, err := NewRPCClient("test", "test", testCrt)
	if err != nil {
		t.Fatal(err)
	}

	ch := rpc.Login(&dauth.User{User: "test", Pass: "test"})
	for res := range ch {
		if res.Err != nil {
			t.Fatal(res.Err)
		}
	}

	ch = rpc.Logout()
	for res := range ch {
		if res.Err != nil {
			t.Fatal(res.Err)
		}
	}

	err = rpc.Close()
	if err != nil {
		t.Fatal(err)
	}
}
