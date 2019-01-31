package passkit

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func testSigner_LoadSigningInformationFromFiles(t *testing.T) {
	signingInfo, err := LoadSigningInformationFromFiles(filepath.Join("test", "passbook", "passkit.p12"), "password", filepath.Join("test", "passbook", "ca.pem"))
	if err != nil {
		t.Errorf("could not load signing info. %v", err)
	}

	_, err = signManifestFile(nil, signingInfo)
	if err == nil {
		t.Errorf("should fail")
	}

	passJson, err := ioutil.ReadFile(filepath.Join("test", "pass2.json"))
	if err != nil {
		t.Errorf("could not load pass json file. %v", err)
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
