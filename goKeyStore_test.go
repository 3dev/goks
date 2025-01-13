package goKeyStore

import (
	"encoding/hex"
	"fmt"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestOpen(t *testing.T) {

	ks, err := Open("nonexistenfile.goKey", "fakepasskey")
	assert.Nil(t, ks)
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {

	ks, err := New("test1.goKey", "test1passkey")
	assert.NotNil(t, ks)
	assert.NoError(t, err)

	ks, err = Open("test1.goKey", "test1passkey")
	assert.Nil(t, err)
	assert.NotNil(t, ks)
	assert.Equal(t, 0, ks.Count())
}

func TestPut(t *testing.T) {

	ks, err := New("test2.goKey", "test2passkey")
	assert.NotNil(t, ks)
	assert.NoError(t, err)

	err = ks.Put("systemData", []byte("linux kernel v4.0 should be out soon"))
	assert.NoError(t, err)
	err = ks.Put("certificate_1",
		[]byte(`-----BEGIN CERTIFICATE-----
MIIDyjCCArKgAwIBAgIUaWGGGVs2rqcWqQVoHvm4m8wOqKYwDQYJKoZIhvcNAQEL
BQAwfTELMAkGA1UEBhMCTkcxDjAMBgNVBAgMBUxhZ29zMRgwFgYDVQQHDA9WaWN0
b3JpYSBJc2xhbmQxFDASBgNVBAoMC0ludGVyc3dpdGNoMREwDwYDVQQLDAhTeXN0
ZWdyYTEbMBkGA1UEAwwSSW50ZXJzd2l0Y2ggUktEIENBMB4XDTI1MDExMDIzMjM0
NFoXDTM1MDEwODIzMjM0NFowZTELMAkGA1UEBhMCTkcxDjAMBgNVBAgMBUxhZ29z
MQ4wDAYDVQQHDAVMYWdvczEUMBIGA1UECgwLSW50ZXJzd2l0Y2gxETAPBgNVBAsM
CFN5c3RlZ3JhMQ0wCwYDVQQDDARIT1NUMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAwF+OitaumF+VdOtfQjMd3OGviaB6uUc8kTWS7dnEExxELgLX8KBX
wuTx9oxtBPRR/bai9Szw46mmZiGD+x2kEwk6HpCbe4oxKRhEqQp/OmcCeGDu+9kV
C1t3SoKiog56dvBxdf5+C2BI0BLrxXNc0/XwC0Y83olBfG9dlkSDxZhvmvxlSt2d
PD7O/tOBxN+bE92/0BYzWJ7OtqCS+Ktt8WYJESkUBT84chGa/A/k5zSkXhGaXevP
Bsvp14YPrdw4sEQhNa3TltFsdDDAskzcHuNqa5rnZGmNfCR+ONUVyB0x5S5YI8pa
EQdmxQynsvlDAJzuxI+EBAiWyk+XW1RIPQIDAQABo1owWDAdBgNVHQ4EFgQUkjga
MrhFbL8DZtNJ5+Yp2VkGq9swHwYDVR0jBBgwFoAUDjUSsOZrmPHgmXZCm4dNcQ0A
dyAwCQYDVR0TBAIwADALBgNVHQ8EBAMCA6gwDQYJKoZIhvcNAQELBQADggEBAD31
Wa5Yf+KwqjhQV39XvTDtnIrxtzpeMI1Jd/hXTyE5hOQIMA7vNnnuHHcXBfc6LqtP
4NA5PJjESuJDnuwWbRuxBYJcpYcJ9ai7tFH2o55wTSoy43jZKK7n/WE5EZSZXTPn
CsYKBHXOrRFQMF90Y5lM3flU5KGoL2i2+/Z844jkSIizbaJPsOXan5ewK28BO8EE
5V5/QTGDZq2ohFSBgY2AJixB5DrhD0of9EN5GNsOVq6Xv/MQxE20WXXtLW+he4SL
G6hrEOwS/ajPFVhZ9ky4hnFxXq19czr0n/4YeYa4Q9/Crl5hKxB1dRtxqJRTV6fP
gv+NyO14VaJIMazvGLk=
-----END CERTIFICATE-----`))
	assert.NoError(t, err)
	assert.Equal(t, 2, ks.Count())
	fmt.Println(hex.Dump(ks.fileHeader.Bytes()))

	fmt.Println(ks.Keys())
}
