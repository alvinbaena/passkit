package passkit

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	BundleIconRetinaHD                = "icon@3x.png"
	BundleIconRetina                  = "icon@2x.png"
	BundleIcon                        = "icon.png"
	BundleLogoRetinaHD                = "logo@3x.png"
	BundleLogoRetina                  = "logo@2x.png"
	BundleLogo                        = "logo.png"
	BundleThumbnailRetinaHD           = "thumbnail@3x.png"
	BundleThumbnailRetina             = "thumbnail@2x.png"
	BundleThumbnail                   = "thumbnail.png"
	BundleStripRetinaHD               = "strip@3x.png"
	BundleStripRetina                 = "strip@2x.png"
	BundleStrip                       = "strip.png"
	BundleBackgroundRetinaHD          = "background@3x.png"
	BundleBackgroundRetina            = "background@2x.png"
	BundleBackground                  = "background.png"
	BundleFooterRetinaHD              = "footer@3x.png"
	BundleFooterRetina                = "footer@2x.png"
	BundleFooter                      = "footer.png"
	BundlePersonalizationLogoRetinaHD = "personalizationLogo@3x.png"
	BundlePersonalizationLogoRetina   = "personalizationLogo@2x.png"
	BundlePersonalizationLogo         = "personalizationLogo.png"
)

type PassTemplate interface {
	ProvisionPassAtDirectory(tmpDirPath string) error
	GetAllFiles() (map[string][]byte, error)
}

type FolderPassTemplate struct {
	templateDir string
}

func NewFolderPassTemplate(templateDir string) *FolderPassTemplate {
	return &FolderPassTemplate{templateDir: templateDir}
}

func (f *FolderPassTemplate) ProvisionPassAtDirectory(tmpDirPath string) error {
	return copyDir(f.templateDir, tmpDirPath)
}

func (f *FolderPassTemplate) GetAllFiles() (map[string][]byte, error) {
	loaded, err := loadDir(f.templateDir)
	if err != nil {
		return nil, err
	}

	ret := make(map[string][]byte)
	for name, data := range loaded {
		ret[filepath.Base(name)] = data
	}

	return ret, err
}

type InMemoryPassTemplate struct {
	files map[string][]byte
	mu    sync.Mutex
}

func NewInMemoryPassTemplate() *InMemoryPassTemplate {
	return &InMemoryPassTemplate{files: make(map[string][]byte)}
}

func (m *InMemoryPassTemplate) ProvisionPassAtDirectory(tmpDirPath string) error {
	dst := filepath.Clean(tmpDirPath)

	_, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(dst, os.ModeDir)
	if err != nil {
		return nil
	}

	for file, d := range m.files {
		err = os.WriteFile(filepath.Join(dst, file), d, 0644)
		if err != nil {
			_ = os.RemoveAll(dst)
			return err
		}
	}

	return nil
}

func (m *InMemoryPassTemplate) GetAllFiles() (map[string][]byte, error) {
	return m.files, nil
}

func (m *InMemoryPassTemplate) AddFileBytes(name string, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[name] = data
}

func (m *InMemoryPassTemplate) AddFileBytesLocalized(name, locale string, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[m.pathForLocale(name, locale)] = data
}

func (m *InMemoryPassTemplate) downloadFile(u url.URL) ([]byte, error) {
	timeout := 10 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (m *InMemoryPassTemplate) AddFileFromURL(name string, u url.URL) error {
	b, err := m.downloadFile(u)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[name] = b
	return nil
}

func (m *InMemoryPassTemplate) AddFileFromURLLocalized(name, locale string, u url.URL) error {
	b, err := m.downloadFile(u)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[m.pathForLocale(name, locale)] = b
	return nil
}

func (m *InMemoryPassTemplate) AddAllFiles(directoryWithFilesToAdd string) error {
	src := filepath.Clean(directoryWithFilesToAdd)
	loaded, err := loadDir(src)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for name, data := range loaded {
		m.files[filepath.Base(name)] = data
	}

	return nil
}

func (m *InMemoryPassTemplate) pathForLocale(name string, locale string) string {
	if strings.TrimSpace(locale) == "" {
		return name
	}

	return filepath.Join(locale+".lproj", name)
}
