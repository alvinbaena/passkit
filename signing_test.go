package passkit

import (
	"os"
	"path/filepath"
	"testing"

	"software.sslmate.com/src/go-pkcs12"
)

func TestSigner_LoadSigningInformationFromFiles(t *testing.T) {
	signingInfo, err := LoadSigningInformationFromFiles(filepath.Join("test", "passbook", "passkit.p12"), "password", filepath.Join("test", "passbook", "ca.pem"))
	if err != nil {
		t.Fatalf("could not load signing info. %v", err)
	}

	_, err = signManifestFile(nil, signingInfo)
	if err == nil {
		t.Errorf("should fail")
	}

	passJson, err := os.ReadFile(filepath.Join("test", "pass2.json"))
	if err != nil {
		t.Fatalf("could not load pass json file. %v", err)
	}

	_, err = signManifestFile(passJson, signingInfo)
	if err != nil {
		t.Errorf("could not sign manifest. %v", err)
	}
}

func TestSigner_LoadSigningInformationFromFilesPaths(t *testing.T) {
	_, err := LoadSigningInformationFromFiles(filepath.Join("test", "passbook", "xxxx"), "xxxxx", filepath.Join("test", "passbook", "AppleWWDRCA.cer"))
	if err == nil {
		t.Errorf("loading cert should fail.")
	}

	_, err = LoadSigningInformationFromFiles(filepath.Join("test", "passbook", "passkit.p12"), "xxxxx", filepath.Join("test", "passbook", "xxxx.cer"))
	if err == nil {
		t.Errorf("loading cert should fail.")
	}
}

func TestSigner_ValidCerts(t *testing.T) {
	_, err := LoadSigningInformationFromFiles(filepath.Join("test", "passbook", "passkit.p12"), "password", filepath.Join("test", "passbook", "ca-chain.cert.pem"))
	if err == nil {
		t.Errorf("should fail")
	}
}

func TestSigner_LoadSigningInformationFromModernPKCS12(t *testing.T) {
	// Re-encode the legacy test keystore with modern algorithms (PBES2, AES-256-CBC, SHA-256 MAC),
	// as OpenSSL 3 exports by default.
	legacy, err := os.ReadFile(filepath.Join("test", "passbook", "passkit.p12"))
	if err != nil {
		t.Fatalf("could not read the legacy keystore. %v", err)
	}
	key, cert, _, err := pkcs12.DecodeChain(legacy, "password")
	if err != nil {
		t.Fatalf("could not decode the legacy keystore. %v", err)
	}
	modern, err := pkcs12.Modern.Encode(key, cert, nil, "password")
	if err != nil {
		t.Fatalf("could not encode the modern keystore. %v", err)
	}

	ca, err := os.ReadFile(filepath.Join("test", "passbook", "ca.pem"))
	if err != nil {
		t.Fatalf("could not read the CA file. %v", err)
	}

	signingInfo, err := LoadSigningInformationFromBytes(modern, "password", ca)
	if err != nil {
		t.Fatalf("could not load signing info from a modern keystore. %v", err)
	}

	passJson, err := os.ReadFile(filepath.Join("test", "pass2.json"))
	if err != nil {
		t.Fatalf("could not load pass json file. %v", err)
	}
	if _, err := signManifestFile(passJson, signingInfo); err != nil {
		t.Errorf("could not sign manifest. %v", err)
	}
}
