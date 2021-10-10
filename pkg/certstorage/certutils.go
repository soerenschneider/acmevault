package certstorage

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	vaultCertKeyPrivateKey = "private_key"
	vaultCertKeyCert       = "cert"
	vaultCertKeyDomain     = "domain"
	vaultCertKeyIssuer     = "dummyIssuer"
	vaultCertKeyUrl        = "url"
	vaultCertKeyStableUrl  = "stable_url"
	vaultCertKeyCsr        = "dummyCsr"

	VaultAccountKeyUri     = "uri"
	VaultAccountKeyEmail   = "email"
	VaultAccountKeyAccount = "account"
	VaultAccountKeyKey     = "key"
)

type CertMetadata struct {
	Expiry time.Time
	Domain string
}

func CertToMap(res *AcmeCertificate) map[string]interface{} {
	if res == nil {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		vaultCertKeyPrivateKey: res.PrivateKey,
		vaultCertKeyCert:       res.Certificate,
		vaultCertKeyDomain:     res.Domain,
		vaultCertKeyIssuer:     res.IssuerCertificate,
		vaultCertKeyUrl:        res.CertURL,
		vaultCertKeyStableUrl:  res.CertStableURL,
		vaultCertKeyCsr:        res.CSR,
	}
}

func MapToCert(data map[string]interface{}) (*AcmeCertificate, error) {
	res := &AcmeCertificate{}
	if data == nil || len(data) < 6 {
		return nil, errors.New("empty/incomplete map provided")
	}

	res.Domain = fmt.Sprint(data[vaultCertKeyDomain])
	res.CertStableURL = fmt.Sprint(data[vaultCertKeyStableUrl])
	res.CertURL = fmt.Sprint(data[vaultCertKeyUrl])

	_, ok := data[vaultCertKeyPrivateKey]; if ok {
		privRaw := fmt.Sprintf("%s", data[vaultCertKeyPrivateKey])
		priv, err := base64.StdEncoding.DecodeString(privRaw)
		if err != nil {
			return nil, fmt.Errorf("can not decode private key: %v", err)
		}
		res.PrivateKey = priv
	}

	certRaw := fmt.Sprintf("%s", data[vaultCertKeyCert])
	cert, err := base64.StdEncoding.DecodeString(certRaw)
	if err != nil {
		return nil, fmt.Errorf("can not decode certificate: %v", err)
	}
	res.Certificate = cert

	issuerRaw := fmt.Sprintf("%s", data[vaultCertKeyIssuer])
	issuer, err := base64.StdEncoding.DecodeString(issuerRaw)
	if err != nil {
		return nil, fmt.Errorf("can not decode issuer cert: %v", err)
	}
	res.IssuerCertificate = issuer

	csrRaw := fmt.Sprintf("%s", data[vaultCertKeyCsr])
	csr, err := base64.StdEncoding.DecodeString(csrRaw)
	if err != nil {
		log.Warn().Msg("Could not decode csr")
	} else {
		res.CSR = csr
	}

	return res, nil
}

func (cert *CertMetadata) GetDurationUntilExpiry() time.Duration {
	return cert.Expiry.Sub(time.Now().UTC())
}

func ConvertToPem(privateKey crypto.PrivateKey) (string, error) {
	pemKey := certcrypto.PEMBlock(privateKey)
	s := &bytes.Buffer{}
	err := pem.Encode(s, pemKey)
	return s.String(), err
}

func FromPem(keyData []byte) (crypto.PrivateKey, error) {
	keyBlock, _ := pem.Decode(keyData)

	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(keyBlock.Bytes)
	}

	return nil, errors.New("unknown private key type")
}
