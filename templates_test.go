package passkit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestInMemoryPassTemplate_ConcurrentAccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("remote file data"))
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("could not parse test server url. %v", err)
	}

	template := NewInMemoryPassTemplate()

	const (
		workers    = 4
		iterations = 50
	)

	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(6)

		go func(n int) {
			defer wg.Done()
			for j := range iterations {
				template.AddFileBytes(fmt.Sprintf("file-%d-%d.png", n, j), []byte("data"))
			}
		}(i)

		go func(n int) {
			defer wg.Done()
			for j := range iterations {
				template.AddFileBytesLocalized(fmt.Sprintf("file-%d-%d.strings", n, j), "en", []byte("data"))
			}
		}(i)

		go func() {
			defer wg.Done()
			for range iterations {
				if err := template.AddFileFromURL(BundleIcon, *u); err != nil {
					t.Errorf("could not add file from url. %v", err)
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			for range iterations {
				if err := template.AddAllFiles(filepath.Join("test", "StoreCard.raw")); err != nil {
					t.Errorf("could not add all files. %v", err)
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			for range iterations {
				files, err := template.GetAllFiles()
				if err != nil {
					t.Errorf("could not get all files. %v", err)
					return
				}
				for name := range files {
					if len(name) == 0 {
						t.Error("unexpected empty template entry name")
						return
					}
				}
			}
		}()

		go func(n int) {
			defer wg.Done()
			base := t.TempDir()
			for j := range 10 {
				dir := filepath.Join(base, fmt.Sprintf("pass-%d-%d", n, j))
				if err := template.ProvisionPassAtDirectory(dir); err != nil {
					t.Errorf("could not provision pass at directory. %v", err)
					return
				}
				_ = os.RemoveAll(dir)
			}
		}(i)
	}

	wg.Wait()
}

func TestInMemoryPassTemplate_GetAllFilesReturnsSnapshot(t *testing.T) {
	template := NewInMemoryPassTemplate()
	template.AddFileBytes(BundleIcon, []byte("icon"))

	files, err := template.GetAllFiles()
	if err != nil {
		t.Fatalf("could not get all files. %v", err)
	}

	files["injected.png"] = []byte("injected")
	template.AddFileBytes(BundleLogo, []byte("logo"))

	current, err := template.GetAllFiles()
	if err != nil {
		t.Fatalf("could not get all files. %v", err)
	}

	if _, ok := current["injected.png"]; ok {
		t.Error("mutating the map returned by GetAllFiles should not affect the template")
	}

	if _, ok := files[BundleLogo]; ok {
		t.Error("files added after GetAllFiles should not show up in a previously returned map")
	}
}
