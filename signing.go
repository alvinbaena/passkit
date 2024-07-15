package passkit

import (
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mozilla.org/pkcs7"
	"software.sslmate.com/src/go-pkcs12"
)

const (
	manifestJsonFileName        = "manifest.json"
	passJsonFileName            = "pass.json"
	personalizationJsonFileName = "personalization.json"
	signatureFileName           = "signature"
)

type Signer interface {
	CreateSignedAndZippedPassArchive(p *Pass, t PassTemplate, i *SigningInformation) ([]byte, error)
	CreateSignedAndZippedPersonalizedPassArchive(p *Pass, pz *Personalization, t PassTemplate, i *SigningInformation) ([]byte, error)
	SignManifestFile(manifestJson []byte, i *SigningInformation) ([]byte, error)
}

type SigningInformation struct {
	signingCert     *x509.Certificate
	appleWWDRCACert *x509.Certificate
	privateKey      interface{}
}

func LoadSigningInformationFromFiles(pkcs12KeyStoreFilePath, keyStorePassword, appleWWDRCAFilePath string) (*SigningInformation, error) {
	p12, err := os.ReadFile(pkcs12KeyStoreFilePath)
	if err != nil {
		return nil, err
	}

	ca, err := os.ReadFile(appleWWDRCAFilePath)
	if err != nil {
		return nil, err
	}

	return LoadSigningInformationFromBytes(p12, keyStorePassword, ca)
}

func LoadSigningInformationFromBytes(pkcs12KeyStoreFile []byte, keyStorePassword string, appleWWDRCAFile []byte) (*SigningInformation, error) {
	info := &SigningInformation{}

	pk, cer, err := pkcs12.Decode(pkcs12KeyStoreFile, keyStorePassword)
	if err != nil {
		return nil, err
	}

	if err := verify(cer); err != nil {
		return nil, fmt.Errorf("error decoding pkcs12: %w", err)
	}

	wwdrca, err := x509.ParseCertificate(appleWWDRCAFile)
	if err != nil {
		return nil, err
	}

	if err := verify(wwdrca); err != nil {
		return nil, fmt.Errorf("error verifying Apple WWDRCAFile: %w", err)
	}

	info.privateKey = pk
	info.signingCert = cer
	info.appleWWDRCACert = wwdrca

	return info, nil
}

// verify checks if a certificate has expired
func verify(cert *x509.Certificate) error {
	_, err := cert.Verify(x509.VerifyOptions{Roots: x509.NewCertPool()})
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case x509.CertificateInvalidError:
		switch e.Reason {
		case x509.Expired:
			return errors.New("certificate has expired or is not yet valid")
		default:
			return err
		}
	case x509.UnknownAuthorityError:
		// Apple cert isn't in the cert pool
		// ignoring this error
		return nil
	default:
		return err
	}
}

func signManifestFile(manifestJson []byte, i *SigningInformation) ([]byte, error) {
	if manifestJson == nil {
		return nil, fmt.Errorf("manifestJson has to be present")
	}

	s, err := pkcs7.NewSignedData(manifestJson)
	if err != nil {
		return nil, err
	}

	s.AddCertificate(i.appleWWDRCACert)

	signingTimeAttr, err := createSigningTimeAttribute()
	if err != nil {
		return nil, err
	}

	signerInfoConfig := pkcs7.SignerInfoConfig{
		ExtraSignedAttributes: []pkcs7.Attribute{signingTimeAttr},
	}

	err = s.AddSigner(i.signingCert, i.privateKey, signerInfoConfig)
	if err != nil {
		return nil, err
	}

	s.Detach()
	return s.Finish()
}

func createSigningTimeAttribute() (pkcs7.Attribute, error) {
	signingTimeBytes, err := asn1.Marshal(time.Now().UTC())
	if err != nil {
		return pkcs7.Attribute{}, err
	}

	return pkcs7.Attribute{
		Type:  asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2},
		Value: signingTimeBytes,
	}, nil
}
