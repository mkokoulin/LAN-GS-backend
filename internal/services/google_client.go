package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"
)

func NewGoogleClient(ctx context.Context, secretPath string, scope ...string) (*http.Client, error) {
	// data, err := os.ReadFile(secretPath)
    // if err != nil {
	// 	return nil, fmt.Errorf("Unable to read client secret file: %v", err)
    // }

	sc := map[string] string {
		"type": "service_account",
		"project_id": "lan-site-94255",
		"private_key_id": "b56432cf737d7862233bb99e9aa07c68e1b105ea",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCznc0tzT9WCNf+\nABR8o6kZnsgKy3Tm4ud6vvP46yJqAFtajHHpLZbniOf+hA/7OUnAiCPdhzdcHNZe\ndgiih04+GumALtJQkXqd/Lgf6Ef0qxQtbHU9aJDGmQWxsv63hWhKBlzXstQrKDb6\nt7p+vwppBBaDBV9+EP8uupTf0fEsH10z+Ru0uRDc4ZJPxq85jpmca5gbOMiyazC7\n0v5X3pLb/Fh6gYh/sfVnFfJx9tS8deXvP7oHhqeTtm/N+PPkQJ1g/sbQykD1AmRS\n1z27Ct6woiOn8k+2P9PIfQgUiGqqJpXpVwRd0g508M8fYYuEfC+rFIQsVtFlnpF9\nDVFjKS1rAgMBAAECggEAIUBa3hpLmd1MQyD3qfVQVkZfSSROMXobGU2lc9Tvy4ji\nYJp0chL3B0kAjc0b+kmqPiLV5Opl1L8f7l/SjGgZE4U+6fSBCdpMxVemLH3/aCuN\nsjUMZUBo4OMuOy3kWswvONkBsgrQnWa1+uctds5eGn/fvQSvH3L9EsUGz4KZr62U\nwCpym+LEyAHFsGCTRTaIcmbKeor9B2kQi4gHq/hHYmxFJb/BRTfUNrvOraVZaXsz\nuPMHIgakXKQP9ADWGKAWgklG3AqmmVC7ZUlNXUL5p4nfSEM1jCFftxFfG/6I+X4s\n2WszI9iKIXzCbBmQ9xIoMHRnry4SbizeHFhQEx9ugQKBgQDye0HtbyT70ez4OkuT\n/sqGvvVyv599pt8D8oZHwGOxVyb9gfAjZDfwETPX94xlC/pUrqypAPtL4r35Sh1e\nFmOwSzuvRg1f2OKmPYy7DUERyEsUdouMHbaio3XocOF2GSeMA2WF8ZICfWdRnmZU\nCpbhkQoliGVKVh+NbOYL++7AgQKBgQC9oVJWHN4u39X1xbdKwIO+80ucMtABPPvU\nEWB64LMCzOixPUXS2c9tCtmUq9pyPwCqAqPBUV73JGQdbZHWXIooO2DC4XBd0v5m\nRg9yiSaiGDzAQPPbBUc5mLgewYGQ8MR/w0ZAsVxPqAMdVNbxntxn4k1MDJiGIy/i\nLpefIkL36wKBgQC4BdUV5djSiBHonQ1IpwB63KeYS1c2XBM6gq9n+tlt+C9uC1P8\n+Az/035d89AHy9xSsjH1HPqaL91vONEq26ESZTZJocd6qzXvZhzMxJVScoiQYhsr\n3k0CBz3vhuOE5jg+KUG+MoRWAWgM6ELOmy7Ax3tE2svMa6oMgc3g4HTkgQKBgAUi\nDsaB9YmzSWljtrhxSZ+rmkpaHcNK0U5GQiRRXMcgoNPbYr54YuMCvi0GEd2x0uTH\nOYOMHlP2Sjd5tc7lpl+8a7wauh3wDi7aiqSBDeipW0ug9njhRbJLbgB3IHi567fB\no28w3dzSIXNzznWv5StytsDuPlqzLSKkPDp0hPeNAoGBANRqfALlFO3hz16Fqvxv\nH5TORkzqHHz5q6QzSuvteOTOn9Za0J+u2qy1ov7+xQOkWSMWHJWDFuHnwhNfvljZ\nR9Lg74LVTC9SU0YIOrT752ZjM/Ip07S6EicoXR5bmFMs1U85ntAVqNyCzAwGwPUO\nI1jAsGoZCGyhbRTCxcwBojJR\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-c8zis@lan-site-94255.iam.gserviceaccount.com",
		"client_id": "115996013580444267171",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-c8zis%40lan-site-94255.iam.gserviceaccount.com",
	}

	b, err := json.Marshal(sc)
    if err != nil {
        panic(err)
    }

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(b, scope...)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// create client with config and context
	return config.Client(ctx), nil
}