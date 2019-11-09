package CachedHttpClient

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestNewFileCache(t *testing.T) {
	fileCache, err := NewFileCache("tmp/request.cache")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transport := &CachedTransport{Cache: fileCache}
	DefaultCashedClient.Transport = transport
	client := DefaultCashedClient
	startTestServerTLS()
	client = client
	request, err := http.NewRequest("GET", serverTLS, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	responseNotCached, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(DefaultCashedClient.Transport.(*CachedTransport).Cache.(*FileCache).cache) != 1 {
		t.Error("request was not save to cache")
		t.FailNow()
	}

	responseCached, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	notCached, err := ioutil.ReadAll(responseNotCached.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cached, err := ioutil.ReadAll(responseCached.Body)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(notCached) != string(cached) {
		t.Error(string(notCached), "!=", string(cached))

	}

}

func TestAppendFileCache(t *testing.T) {
	cacheFile := "tmp/append.request.cache"

	fileContent := `{"Request":"GET / HTTP/1.1\r\nHost: localhost:8081\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n","Response":{"Status":"200 OK","StatusCode":200,"Proto":"HTTP/1.1","ProtoMajor":1,"ProtoMinor":1,"Header":{"Content-Length":["19"],"Content-Type":["text/plain; charset=utf-8"],"Date":["Sat, 09 Nov 2019 02:41:51 GMT"]},"Body":"ODY3NDY2NTIyMzA4MjE1MzU1MQ==","ContentLength":19,"TransferEncoding":null,"Close":false,"Uncompressed":false,"Trailer":null,"Request":"","TLS":{"Version":772,"HandshakeComplete":true,"DidResume":false,"CipherSuite":4865,"NegotiatedProtocol":"","NegotiatedProtocolIsMutual":true,"ServerName":"","PeerCertificates":[{"Raw":"MIIC+TCCAeGgAwIBAgIQJ9phBHlJ/3w9cKMe1HoruTANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQKEwdBY21lIENvMB4XDTE5MTEwODE3MDcxOVoXDTIwMTEwNzE3MDcxOVowEjEQMA4GA1UEChMHQWNtZSBDbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMZ9LLXONHURuLVmYgW+ZEvgKvCGcju905hazdaiQMQypCa9T17NiVzuBxeKQzRc3SdyxL/gAp94YwyRWddXYY1WVLo7VH1dY3BPo2A7rZwrCpKvP9ubLkaUkgfPyCk3sS6pug/+A9RgmquHc6lm4QSGr5v6AWmF2ZY1IiEVl/N37jPtAyavgWMgXXe8pHt5S36ci2z79EfonkRBAX/MWJEqjL7BaF9CSupxji2pgd3GDyUQAWGJKwYPxqQOqPYD3XLYbPi/VvXWKalsc/d9I6ZhPfye2f2W9feQzkPIzzsuPRUXdKKyM5E+rq8VR9RYOU+Iwfy96m3LfLnGcOguDm8CAwEAAaNLMEkwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwFAYDVR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQAlp4i253gCadP+eJtqVuvt+IL1DIvNu36xiPYj3fw9hs0TnGhyu0ckbXpMksyDVF9TONpYkS6EgrHGViKHUaJljxe3BCbugZvDcNUA5Kz8PPaRkbPlB3sUDcZPAnzzhWwruhfYv7w2DTT6Px35dJKYmiS3ZS63RDSru1eF4sV3oAXEmow1gEeZiKkcxYMjKlLtlJ2J/rIv1+KB0eQ5MlQXiymvb9XqNX+RosKXN3nUYT9Zdqp449ogeeMeibMe21gnkDfBNMGnMLCr/PSdzsVtYFSsRSZXyyR6/G0tFq+XZ7oNqgO+otEooGHHL7FQFnpcR702UqpnwAsZPnIyJwhs","RawTBSCertificate":"MIIB4aADAgECAhAn2mEEeUn/fD1wox7Ueiu5MA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNVBAoTB0FjbWUgQ28wHhcNMTkxMTA4MTcwNzE5WhcNMjAxMTA3MTcwNzE5WjASMRAwDgYDVQQKEwdBY21lIENvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxn0stc40dRG4tWZiBb5kS+Aq8IZyO73TmFrN1qJAxDKkJr1PXs2JXO4HF4pDNFzdJ3LEv+ACn3hjDJFZ11dhjVZUujtUfV1jcE+jYDutnCsKkq8/25suRpSSB8/IKTexLqm6D/4D1GCaq4dzqWbhBIavm/oBaYXZljUiIRWX83fuM+0DJq+BYyBdd7yke3lLfpyLbPv0R+ieREEBf8xYkSqMvsFoX0JK6nGOLamB3cYPJRABYYkrBg/GpA6o9gPdcths+L9W9dYpqWxz930jpmE9/J7Z/Zb195DOQ8jPOy49FRd0orIzkT6urxVH1Fg5T4jB/L3qbct8ucZw6C4ObwIDAQABo0swSTAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADAUBgNVHREEDTALgglsb2NhbGhvc3Q=","RawSubjectPublicKeyInfo":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxn0stc40dRG4tWZiBb5kS+Aq8IZyO73TmFrN1qJAxDKkJr1PXs2JXO4HF4pDNFzdJ3LEv+ACn3hjDJFZ11dhjVZUujtUfV1jcE+jYDutnCsKkq8/25suRpSSB8/IKTexLqm6D/4D1GCaq4dzqWbhBIavm/oBaYXZljUiIRWX83fuM+0DJq+BYyBdd7yke3lLfpyLbPv0R+ieREEBf8xYkSqMvsFoX0JK6nGOLamB3cYPJRABYYkrBg/GpA6o9gPdcths+L9W9dYpqWxz930jpmE9/J7Z/Zb195DOQ8jPOy49FRd0orIzkT6urxVH1Fg5T4jB/L3qbct8ucZw6C4ObwIDAQAB","RawSubject":"MBIxEDAOBgNVBAoTB0FjbWUgQ28=","RawIssuer":"MBIxEDAOBgNVBAoTB0FjbWUgQ28=","Signature":"JaeItud4AmnT/nibalbr7fiC9QyLzbt+sYj2I938PYbNE5xocrtHJG16TJLMg1RfUzjaWJEuhIKxxlYih1GiZY8XtwQm7oGbw3DVAOSs/Dz2kZGz5Qd7FA3GTwJ884VsK7oX2L+8Ng00+j8d+XSSmJokt2Uut0Q0q7tXheLFd6AFxJqMNYBHmYipHMWDIypS7ZSdif6yL9figdHkOTJUF4spr2/V6jV/kaLClzd51GE/WXaqeOPaIHnjHomzHttYJ5A3wTTBpzCwq/z0nc7FbWBUrEUmV8skevxtLRavl2e6DaoDvqLRKKBhxy+xUBZ6XEe9NlKqZ8ALGT5yMicIbA==","SignatureAlgorithm":4,"PublicKeyAlgorithm":1,"PublicKey":{"N":25056910303322939806583737109066884128144601853459127274697308916781949953377221483643154774177588940904379509181716401514416650623217069460729444857057615083081331363638758249729080407640027970863576709940108814737745511078397909809351720896613772748200709286330407151844569287737450280018138569479167668104908020255387095437799528742067315022017830712804762585236364659341877595921567763479736770106360791752526434484928885751234829250812425793997603305574706701161538359280936252114787878128057147267011450045207254151857807491048447080779224578342720538337928623052627161367343127581834241292250703952028422245999,"E":65537},"Version":3,"SerialNumber":52973780298953660003847832739734236089,"Issuer":{"Country":null,"Organization":["Acme Co"],"OrganizationalUnit":null,"Locality":null,"Province":null,"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"","Names":[{"Type":[2,5,4,10],"Value":"Acme Co"}],"ExtraNames":null},"Subject":{"Country":null,"Organization":["Acme Co"],"OrganizationalUnit":null,"Locality":null,"Province":null,"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"","Names":[{"Type":[2,5,4,10],"Value":"Acme Co"}],"ExtraNames":null},"NotBefore":"2019-11-08T17:07:19Z","NotAfter":"2020-11-07T17:07:19Z","KeyUsage":5,"Extensions":[{"Id":[2,5,29,15],"Critical":true,"Value":"AwIFoA=="},{"Id":[2,5,29,37],"Critical":false,"Value":"MAoGCCsGAQUFBwMB"},{"Id":[2,5,29,19],"Critical":true,"Value":"MAA="},{"Id":[2,5,29,17],"Critical":false,"Value":"MAuCCWxvY2FsaG9zdA=="}],"ExtraExtensions":null,"UnhandledCriticalExtensions":null,"ExtKeyUsage":[1],"UnknownExtKeyUsage":null,"BasicConstraintsValid":true,"IsCA":false,"MaxPathLen":-1,"MaxPathLenZero":false,"SubjectKeyId":null,"AuthorityKeyId":null,"OCSPServer":null,"IssuingCertificateURL":null,"DNSNames":["localhost"],"EmailAddresses":null,"IPAddresses":null,"URIs":null,"PermittedDNSDomainsCritical":false,"PermittedDNSDomains":null,"ExcludedDNSDomains":null,"PermittedIPRanges":null,"ExcludedIPRanges":null,"PermittedEmailAddresses":null,"ExcludedEmailAddresses":null,"PermittedURIDomains":null,"ExcludedURIDomains":null,"CRLDistributionPoints":null,"PolicyIdentifiers":null}],"VerifiedChains":null,"SignedCertificateTimestamps":null,"OCSPResponse":null,"TLSUnique":null}}}`

	file, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = file.WriteString(fileContent)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fileCache, err := OpenFileCache(cacheFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(fileCache.cache) != 1 {
		t.Error("map from file not correctly loaded")
		t.FailNow()
	}

	transport := *DefaultCachedTransport
	transport.Cache = fileCache
	DefaultCashedClient.Transport = &transport

	startTestServerTLS()

	request, err := http.NewRequest("GET", serverTLS, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	responseNotCached, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(DefaultCashedClient.Transport.(*CachedTransport).Cache.(*FileCache).cache) != 1 {
		t.Error("request was not save to cache")
		t.FailNow()
	}

	responseCached, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	notCached, err := ioutil.ReadAll(responseNotCached.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cached, err := ioutil.ReadAll(responseCached.Body)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(notCached) != string(cached) {
		t.Error(string(notCached), "!=", string(cached))

	}
}
