package passkit

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

type folderPassTemplate struct {
	templateDir string
}

func NewFolderPassTemplate(templateDir string) *folderPassTemplate {
	return &folderPassTemplate{templateDir: templateDir}
}

func (f *folderPassTemplate) ProvisionPassAtDirectory(tmpDirPath string) error {
	return copyDir(f.templateDir, tmpDirPath)
}

func (f *folderPassTemplate) GetAllFiles() (map[string][]byte, error) {
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

type inMemoryPassTemplate struct {
	files map[string][]byte
}

	func NewInMemoryPassTemplate() *inMemoryPassTemplate {
	return &inMemoryPassTemplate{files: make(map[string][]byte)}
}

func (m *inMemoryPassTemplate) ProvisionPassAtDirectory(tmpDirPath string) error {
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
		err = ioutil.WriteFile(filepath.Join(dst, string(file)), d, 0644)
		if err != nil {
			os.RemoveAll(dst)
			return err
		}
	}

	return nil
}

func (m *inMemoryPassTemplate) GetAllFiles() (map[string][]byte, error) {
	return m.files, nil
}

func (m *inMemoryPassTemplate) AddFileBytes(name string, data []byte) {
	m.files[name] = data
}

func (m *inMemoryPassTemplate) AddFileBytesLocalized(name, locale string, data []byte) {
	m.files[m.pathForLocale(name, locale)] = data
}

func (m *inMemoryPassTemplate) downloadFile(u url.URL) ([]byte, error) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (m *inMemoryPassTemplate) AddFileFromURL(name string, u url.URL) error {
	b, err := m.downloadFile(u)
	if err != nil {
		return err
	}

	m.files[name] = b
	return nil
}

func (m *inMemoryPassTemplate) AddFileFromURLLocalized(name, locale string, u url.URL) error {
	b, err := m.downloadFile(u)
	if err != nil {
		return err
	}

	m.files[m.pathForLocale(name, locale)] = b
	return nil
}

func (m *inMemoryPassTemplate) AddAllFiles(directoryWithFilesToAdd string) error {
	src := filepath.Clean(directoryWithFilesToAdd)
	loaded, err := loadDir(src)
	if err != nil {
		return err
	}

	for name, data := range loaded {
		m.files[filepath.Base(name)] = data
	}

	return nil
}

func (m *inMemoryPassTemplate) pathForLocale(name string, locale string) string {
	if strings.TrimSpace(locale) == "" {
		return name
	}

	return filepath.Join(locale+".lproj", name)
}
