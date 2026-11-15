package prod

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/eyas"
)

type ProdModel struct {
	Name                string `json:"name"`
	NameDesc            string `json:"name_desc"`
	Footer              string `json:"footer"`
	LogoExt             string `json:"logo_ext,omitempty"`
	LogoUpdatedAt       int64  `json:"logo_updated_at,omitempty"`
	FaviconExt          string `json:"favicon_ext,omitempty"`
	FaviconUpdatedAt    int64  `json:"favicon_updated_at,omitempty"`
	BackgroundExt       string `json:"background_ext,omitempty"`
	BackgroundUpdatedAt int64  `json:"background_updated_at,omitempty"`
	isNew               bool
}

func (m *ProdModel) Key() string {
	return "prod-model-info"
}

func (m *ProdModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *ProdModel) Save() error {
	db := sqldb.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"key", "value"}
		values := []string{"$1", "$2"}

		insertClause := fmt.Sprintf("insert into k_v_table(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)
		_, err := db.Exec(insertClause, m.Key(), m.Value())
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		updateClause := "update k_v_table set value=$1 where key=$2"
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.Exec(updateClause, m.Value(), m.Key())
		if err != nil {
			logger.Errorf("update model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	return nil
}

func GetProdModel() *ProdModel {
	var m ProdModel

	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, m.Key())
	if row == nil {
		return nil
	}

	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Debugf("prod model row scan failed: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("Prod model unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}

func NewProdModel() *ProdModel {
	return &ProdModel{
		isNew: true,
	}
}

func GetProdBasisInfo() (ProdModel, bool) {
	m := GetProdModel()
	if m == nil {
		return ProdModel{}, false
	}
	return *m, true
}

type ProdInfoSetRequest struct {
	Name     string `json:"name"`
	NameDesc string `json:"name_desc"`
	Footer   string `json:"footer"`
}

func SetProdBasisInfo(req *ProdInfoSetRequest) error {
	m := GetProdModel()
	if m == nil {
		m = NewProdModel()
	}
	m.Name = req.Name
	m.Footer = req.Footer
	m.NameDesc = req.NameDesc
	return m.Save()
}

func GetBrandingDir() string {
	dir := filepath.Join(eyas.GetDataDir(), "branding")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func GetLogoPath() (string, bool) {
	m := GetProdModel()
	if m != nil && m.LogoExt != "" {
		p := filepath.Join(GetBrandingDir(), "logo"+m.LogoExt)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	// Fallback to legacy logo.png if present
	legacy := filepath.Join(eyas.GetDataDir(), "logo.png")
	if _, err := os.Stat(legacy); err == nil {
		return legacy, true
	}
	return "", false
}

func SetLogoExt(ext string) error {
	m := GetProdModel()
	if m == nil {
		m = NewProdModel()
	}
	if m.LogoExt != "" && m.LogoExt != ext {
		_ = os.Remove(filepath.Join(GetBrandingDir(), "logo"+m.LogoExt))
	}
	m.LogoExt = ext
	m.LogoUpdatedAt = time.Now().Unix()
	return m.Save()
}

func DeleteLogo() error {
	m := GetProdModel()
	if m == nil {
		return nil
	}
	if m.LogoExt != "" {
		_ = os.Remove(filepath.Join(GetBrandingDir(), "logo"+m.LogoExt))
		m.LogoExt = ""
		m.LogoUpdatedAt = time.Now().Unix()
	}
	_ = os.Remove(filepath.Join(eyas.GetDataDir(), "logo.png"))
	return m.Save()
}

func GetFaviconPath() (string, bool) {
	m := GetProdModel()
	if m != nil && m.FaviconExt != "" {
		p := filepath.Join(GetBrandingDir(), "favicon"+m.FaviconExt)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

func SetFaviconExt(ext string) error {
	m := GetProdModel()
	if m == nil {
		m = NewProdModel()
	}
	if m.FaviconExt != "" && m.FaviconExt != ext {
		_ = os.Remove(filepath.Join(GetBrandingDir(), "favicon"+m.FaviconExt))
	}
	m.FaviconExt = ext
	m.FaviconUpdatedAt = time.Now().Unix()
	return m.Save()
}

func DeleteFavicon() error {
	m := GetProdModel()
	if m == nil {
		return nil
	}
	if m.FaviconExt != "" {
		_ = os.Remove(filepath.Join(GetBrandingDir(), "favicon"+m.FaviconExt))
		m.FaviconExt = ""
		m.FaviconUpdatedAt = time.Now().Unix()
	}
	return m.Save()
}

func GetBackgroundPath() (string, bool) {
	m := GetProdModel()
	if m != nil && m.BackgroundExt != "" {
		p := filepath.Join(GetBrandingDir(), "background"+m.BackgroundExt)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	// Fallback to legacy background.png if present
	legacy := filepath.Join(eyas.GetDataDir(), "background.png")
	if _, err := os.Stat(legacy); err == nil {
		return legacy, true
	}
	return "", false
}

func SetBackgroundExt(ext string) error {
	m := GetProdModel()
	if m == nil {
		m = NewProdModel()
	}
	if m.BackgroundExt != "" && m.BackgroundExt != ext {
		_ = os.Remove(filepath.Join(GetBrandingDir(), "background"+m.BackgroundExt))
	}
	m.BackgroundExt = ext
	m.BackgroundUpdatedAt = time.Now().Unix()
	return m.Save()
}

func DeleteBackground() error {
	m := GetProdModel()
	if m == nil {
		return nil
	}
	if m.BackgroundExt != "" {
		_ = os.Remove(filepath.Join(GetBrandingDir(), "background"+m.BackgroundExt))
		m.BackgroundExt = ""
		m.BackgroundUpdatedAt = time.Now().Unix()
	}
	_ = os.Remove(filepath.Join(eyas.GetDataDir(), "background.png"))
	return m.Save()
}

func ResetProdBasisInfo() error {
	_ = DeleteLogo()
	_ = DeleteFavicon()
	_ = DeleteBackground()
	_ = os.RemoveAll(GetBrandingDir())

	db := sqldb.GetDB()
	_, err := db.Exec("delete from k_v_table where key=$1", "prod-model-info")
	return err
}
