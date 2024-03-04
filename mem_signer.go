package passkit

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

type memorySigner struct {
}

func NewMemoryBasedSigner() Signer {
	return &memorySigner{}
}

func (m *memorySigner) CreateSignedAndZippedPassArchive(p *Pass, t PassTemplate, i *SigningInformation) ([]byte, error) {
	return m.CreateSignedAndZippedPersonalizedPassArchive(p, nil, t, i)
}

func (m *memorySigner) CreateSignedAndZippedPersonalizedPassArchive(p *Pass, pz *Personalization, t PassTemplate, i *SigningInformation) ([]byte, error) {
	originalFiles, err := t.GetAllFiles()
	if err != nil {
		return nil, err
	}
	files := m.makeFilesCopy(originalFiles)

	if !p.IsValid() {
		return nil, fmt.Errorf("%v", p.GetValidationErrors())
	}

	pb, err := p.toJSON()
	if err != nil {
		return nil, err
	}

	files[passJsonFileName] = pb

	if pz != nil {
		if !pz.IsValid() {
			return nil, fmt.Errorf("%v", pz.GetValidationErrors())
		}

		pzb, err := pz.toJSON()
		if err != nil {
			return nil, err
		}

		files[personalizationJsonFileName] = pzb
	}

	msftHash, err := m.hashFiles(files)
	if err != nil {
		return nil, err
	}

	mfst, err := json.Marshal(msftHash)
	if err != nil {
		return nil, err
	}

	files[manifestJsonFileName] = mfst

	signedMfst, err := signManifestFile(mfst, i)
	if err != nil {
		return nil, err
	}

	files[signatureFileName] = signedMfst

	z, err := m.createZipFile(files)
	if err != nil {
		return nil, err
	}

	return z, nil
}

func (m *memorySigner) SignManifestFile(manifestJson []byte, i *SigningInformation) ([]byte, error) {
	return signManifestFile(manifestJson, i)
}

func (m *memorySigner) hashFiles(files map[string][]byte) (map[string]string, error) {
	ret := make(map[string]string)
	for name, data := range files {
		hash := sha1.Sum(data)
		ret[name] = fmt.Sprintf("%x", hash)
	}

	return ret, nil
}

func (m *memorySigner) createZipFile(files map[string][]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for name, data := range files {
		f, err := w.Create(name)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(data)
		if err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *memorySigner) makeFilesCopy(files map[string][]byte) map[string][]byte {
	filesCopy := make(map[string][]byte, len(files))
	for k := range files {
		filesCopy[k] = files[k]
	}

	return filesCopy
}
