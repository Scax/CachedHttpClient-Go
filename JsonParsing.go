package CachedHttpClient

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"time"
)

type JsonResponse struct {
	Status           string
	StatusCode       int
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             []byte
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Uncompressed     bool
	Trailer          http.Header
	Request          string
	TLS              *JsonTlsConnectionState
}

func NewJsonResponse(res *http.Response) (*JsonResponse, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	return &JsonResponse{
		Status:           res.Status,
		StatusCode:       res.StatusCode,
		Proto:            res.Proto,
		ProtoMajor:       res.ProtoMajor,
		ProtoMinor:       res.ProtoMinor,
		Header:           res.Header,
		Body:             buf.Bytes(),
		ContentLength:    res.ContentLength,
		TransferEncoding: res.TransferEncoding,
		Close:            res.Close,
		Uncompressed:     res.Uncompressed,
		Trailer:          res.Trailer,
		Request:          "",
		TLS:              NewJsonTlsConnectionState(res.TLS),
	}, nil
}
func (response JsonResponse) ToResponse() *http.Response {

	var res = http.Response{
		Status:           response.Status,
		StatusCode:       response.StatusCode,
		Proto:            response.Proto,
		ProtoMajor:       response.ProtoMajor,
		ProtoMinor:       response.ProtoMinor,
		Header:           response.Header,
		Body:             ioutil.NopCloser(bytes.NewBuffer(response.Body)),
		ContentLength:    response.ContentLength,
		TransferEncoding: response.TransferEncoding,
		Close:            response.Close,
		Uncompressed:     response.Uncompressed,
		Trailer:          response.Trailer,
		Request:          nil,
		TLS:              response.TLS.ToConnectionState(),
	}

	return &res

}

func responseToJSON(res *http.Response) ([]byte, error) {
	response, err := NewJsonResponse(res)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return marshal, nil

}

type JsonTlsConnectionState struct {
	Version                     uint16
	HandshakeComplete           bool
	DidResume                   bool
	CipherSuite                 uint16
	NegotiatedProtocol          string
	NegotiatedProtocolIsMutual  bool
	ServerName                  string
	PeerCertificates            []*JsonX509Certificate
	VerifiedChains              [][]*JsonX509Certificate
	SignedCertificateTimestamps [][]byte
	OCSPResponse                []byte
	TLSUnique                   []byte
}

func NewJsonTlsConnectionState(tls *tls.ConnectionState) *JsonTlsConnectionState {

	if tls == nil {
		return nil
	}

	return &JsonTlsConnectionState{
		Version:                     tls.Version,
		HandshakeComplete:           tls.HandshakeComplete,
		DidResume:                   tls.DidResume,
		CipherSuite:                 tls.CipherSuite,
		NegotiatedProtocol:          tls.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  tls.NegotiatedProtocolIsMutual,
		ServerName:                  tls.ServerName,
		PeerCertificates:            NewJsonX509CertificateArray(tls.PeerCertificates),
		VerifiedChains:              NewJsonX509CertificateArrayArray(tls.VerifiedChains),
		SignedCertificateTimestamps: tls.SignedCertificateTimestamps,
		OCSPResponse:                tls.OCSPResponse,
		TLSUnique:                   tls.TLSUnique,
	}
}
func (state *JsonTlsConnectionState) ToConnectionState() *tls.ConnectionState {

	return &tls.ConnectionState{
		Version:                     state.Version,
		HandshakeComplete:           state.HandshakeComplete,
		DidResume:                   state.DidResume,
		CipherSuite:                 state.CipherSuite,
		NegotiatedProtocol:          state.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  state.NegotiatedProtocolIsMutual,
		ServerName:                  state.ServerName,
		PeerCertificates:            ToX509CertificateArray(state.PeerCertificates),
		VerifiedChains:              ToX509CertificateArrayArray(state.VerifiedChains),
		SignedCertificateTimestamps: state.SignedCertificateTimestamps,
		OCSPResponse:                state.OCSPResponse,
		TLSUnique:                   state.TLSUnique,
	}
}

type JsonX509Certificate struct {
	Raw                         []byte
	RawTBSCertificate           []byte
	RawSubjectPublicKeyInfo     []byte
	RawSubject                  []byte
	RawIssuer                   []byte
	Signature                   []byte
	SignatureAlgorithm          x509.SignatureAlgorithm
	PublicKeyAlgorithm          x509.PublicKeyAlgorithm
	PublicKey                   *rsa.PublicKey
	Version                     int
	SerialNumber                *big.Int
	Issuer                      pkix.Name
	Subject                     pkix.Name
	NotBefore, NotAfter         time.Time
	KeyUsage                    x509.KeyUsage
	Extensions                  []pkix.Extension
	ExtraExtensions             []pkix.Extension
	UnhandledCriticalExtensions []asn1.ObjectIdentifier
	ExtKeyUsage                 []x509.ExtKeyUsage
	UnknownExtKeyUsage          []asn1.ObjectIdentifier
	BasicConstraintsValid       bool
	IsCA                        bool
	MaxPathLen                  int
	MaxPathLenZero              bool
	SubjectKeyId                []byte
	AuthorityKeyId              []byte
	OCSPServer                  []string
	IssuingCertificateURL       []string
	DNSNames                    []string
	EmailAddresses              []string
	IPAddresses                 []net.IP
	URIs                        []*url.URL
	PermittedDNSDomainsCritical bool
	PermittedDNSDomains         []string
	ExcludedDNSDomains          []string
	PermittedIPRanges           []*net.IPNet
	ExcludedIPRanges            []*net.IPNet
	PermittedEmailAddresses     []string
	ExcludedEmailAddresses      []string
	PermittedURIDomains         []string
	ExcludedURIDomains          []string
	CRLDistributionPoints       []string
	PolicyIdentifiers           []asn1.ObjectIdentifier
}

func (certificate *JsonX509Certificate) ToCertificate() *x509.Certificate {

	cert := x509.Certificate{
		Raw:                         certificate.Raw,
		RawTBSCertificate:           certificate.RawTBSCertificate,
		RawSubjectPublicKeyInfo:     certificate.RawSubjectPublicKeyInfo,
		RawSubject:                  certificate.RawSubject,
		RawIssuer:                   certificate.RawIssuer,
		Signature:                   certificate.Signature,
		SignatureAlgorithm:          certificate.SignatureAlgorithm,
		PublicKeyAlgorithm:          certificate.PublicKeyAlgorithm,
		PublicKey:                   certificate.PublicKey,
		Version:                     certificate.Version,
		SerialNumber:                certificate.SerialNumber,
		Issuer:                      certificate.Issuer,
		Subject:                     certificate.Subject,
		NotBefore:                   certificate.NotBefore,
		NotAfter:                    certificate.NotAfter,
		KeyUsage:                    certificate.KeyUsage,
		Extensions:                  certificate.Extensions,
		ExtraExtensions:             certificate.ExtraExtensions,
		UnhandledCriticalExtensions: certificate.UnhandledCriticalExtensions,
		ExtKeyUsage:                 certificate.ExtKeyUsage,
		UnknownExtKeyUsage:          certificate.UnknownExtKeyUsage,
		BasicConstraintsValid:       certificate.BasicConstraintsValid,
		IsCA:                        certificate.IsCA,
		MaxPathLen:                  certificate.MaxPathLen,
		MaxPathLenZero:              certificate.MaxPathLenZero,
		SubjectKeyId:                certificate.SubjectKeyId,
		AuthorityKeyId:              certificate.AuthorityKeyId,
		OCSPServer:                  certificate.OCSPServer,
		IssuingCertificateURL:       certificate.IssuingCertificateURL,
		DNSNames:                    certificate.DNSNames,
		EmailAddresses:              certificate.EmailAddresses,
		IPAddresses:                 certificate.IPAddresses,
		URIs:                        certificate.URIs,
		PermittedDNSDomainsCritical: certificate.PermittedDNSDomainsCritical,
		PermittedDNSDomains:         certificate.PermittedDNSDomains,
		ExcludedDNSDomains:          certificate.ExcludedDNSDomains,
		PermittedIPRanges:           certificate.PermittedIPRanges,
		ExcludedIPRanges:            certificate.ExcludedIPRanges,
		PermittedEmailAddresses:     certificate.PermittedEmailAddresses,
		ExcludedEmailAddresses:      certificate.ExcludedEmailAddresses,
		PermittedURIDomains:         certificate.PermittedURIDomains,
		ExcludedURIDomains:          certificate.ExcludedURIDomains,
		CRLDistributionPoints:       certificate.CRLDistributionPoints,
		PolicyIdentifiers:           certificate.PolicyIdentifiers,
	}

	return &cert

}

func NewJsonX509Certificate(cert *x509.Certificate) *JsonX509Certificate {
	return &JsonX509Certificate{
		Raw:                         cert.Raw,
		RawTBSCertificate:           cert.RawTBSCertificate,
		RawSubjectPublicKeyInfo:     cert.RawSubjectPublicKeyInfo,
		RawSubject:                  cert.RawSubject,
		RawIssuer:                   cert.RawIssuer,
		Signature:                   cert.Signature,
		SignatureAlgorithm:          cert.SignatureAlgorithm,
		PublicKeyAlgorithm:          cert.PublicKeyAlgorithm,
		PublicKey:                   cert.PublicKey.(*rsa.PublicKey),
		Version:                     cert.Version,
		SerialNumber:                cert.SerialNumber,
		Issuer:                      cert.Issuer,
		Subject:                     cert.Subject,
		NotBefore:                   cert.NotBefore,
		NotAfter:                    cert.NotAfter,
		KeyUsage:                    cert.KeyUsage,
		Extensions:                  cert.Extensions,
		ExtraExtensions:             cert.ExtraExtensions,
		UnhandledCriticalExtensions: cert.UnhandledCriticalExtensions,
		ExtKeyUsage:                 cert.ExtKeyUsage,
		UnknownExtKeyUsage:          cert.UnknownExtKeyUsage,
		BasicConstraintsValid:       cert.BasicConstraintsValid,
		IsCA:                        cert.IsCA,
		MaxPathLen:                  cert.MaxPathLen,
		MaxPathLenZero:              cert.MaxPathLenZero,
		SubjectKeyId:                cert.SubjectKeyId,
		AuthorityKeyId:              cert.AuthorityKeyId,
		OCSPServer:                  cert.OCSPServer,
		IssuingCertificateURL:       cert.IssuingCertificateURL,
		DNSNames:                    cert.DNSNames,
		EmailAddresses:              cert.EmailAddresses,
		IPAddresses:                 cert.IPAddresses,
		URIs:                        cert.URIs,
		PermittedDNSDomainsCritical: cert.PermittedDNSDomainsCritical,
		PermittedDNSDomains:         cert.PermittedDNSDomains,
		ExcludedDNSDomains:          cert.ExcludedDNSDomains,
		PermittedIPRanges:           cert.PermittedIPRanges,
		ExcludedIPRanges:            cert.ExcludedIPRanges,
		PermittedEmailAddresses:     cert.PermittedEmailAddresses,
		ExcludedEmailAddresses:      cert.ExcludedEmailAddresses,
		PermittedURIDomains:         cert.PermittedURIDomains,
		ExcludedURIDomains:          cert.ExcludedURIDomains,
		CRLDistributionPoints:       cert.CRLDistributionPoints,
		PolicyIdentifiers:           cert.PolicyIdentifiers,
	}
}
func NewJsonX509CertificateArray(certs []*x509.Certificate) []*JsonX509Certificate {
	if certs == nil {
		return nil
	}
	var array = make([]*JsonX509Certificate, len(certs))
	for k, v := range certs {
		array[k] = NewJsonX509Certificate(v)
	}

	return array

}
func NewJsonX509CertificateArrayArray(certs [][]*x509.Certificate) [][]*JsonX509Certificate {
	if certs == nil {
		return nil
	}
	var array = make([][]*JsonX509Certificate, len(certs))
	for k, v := range certs {
		array[k] = NewJsonX509CertificateArray(v)
	}

	return array

}

func ToX509CertificateArrayArray(certificates [][]*JsonX509Certificate) [][]*x509.Certificate {
	if certificates == nil {
		return nil
	}
	certs := make([][]*x509.Certificate, len(certificates))

	for k, v := range certificates {
		certs[k] = ToX509CertificateArray(v)
	}

	return certs

}
func ToX509CertificateArray(certificates []*JsonX509Certificate) []*x509.Certificate {

	if certificates == nil {
		return nil
	}

	var certs = make([]*x509.Certificate, len(certificates))

	for k, v := range certificates {
		certs[k] = v.ToCertificate()
	}

	return certs
}
