package passkit

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

func TestMemorySigner_ConcurrentSignAndTemplateMutation(t *testing.T) {
	signingInfo, err := LoadSigningInformationFromFiles(filepath.Join("test", "passbook", "passkit.p12"), "password", filepath.Join("test", "passbook", "ca.pem"))
	if err != nil {
		t.Fatalf("could not load signing information. %v", err)
	}

	template := NewInMemoryPassTemplate()
	template.AddFileBytes(BundleIcon, []byte("icon"))

	pass := &Pass{
		FormatVersion:      1,
		SerialNumber:       "1234",
		PassTypeIdentifier: "pass.test.concurrency",
		TeamIdentifier:     "TEAMID",
		Description:        "concurrency regression test",
		OrganizationName:   "test org",
		Generic:            NewGenericPass(),
	}
	if !pass.IsValid() {
		t.Fatalf("test pass is not valid. %v", pass.GetValidationErrors())
	}

	signer := NewMemoryBasedSigner()

	const (
		workers    = 4
		iterations = 25
	)

	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(2)

		go func(n int) {
			defer wg.Done()
			for j := range iterations {
				template.AddFileBytes(fmt.Sprintf("file-%d-%d.png", n, j), []byte("data"))
			}
		}(i)

		go func() {
			defer wg.Done()
			for range iterations {
				z, err := signer.CreateSignedAndZippedPassArchive(pass, template, signingInfo)
				if err != nil {
					t.Errorf("could not sign pass. %v", err)
					return
				}
				if len(z) == 0 {
					t.Error("signed pass archive is empty")
					return
				}
			}
		}()
	}

	wg.Wait()
}
