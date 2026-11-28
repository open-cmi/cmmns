package portal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
)

const defaultPortalDir = "/etc/eyas/portal.d"

// Product represents a portal product entry
type Product struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

// ScanProducts scans /etc/eyas/portal.d/*.json and returns product entries
func ScanProducts(dir string) ([]Product, error) {
	if dir == "" {
		dir = defaultPortalDir
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	products := make([]Product, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		portalFile := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(portalFile)
		if err != nil {
			logger.Debugf("read %s failed: %s\n", portalFile, err.Error())
			continue
		}

		var item Product
		if err := json.Unmarshal(data, &item); err != nil {
			logger.Warnf("invalid json in %s: %s\n", portalFile, err.Error())
			continue
		}
		if item.Name == "" || item.Prefix == "" {
			continue
		}

		products = append(products, item)
	}

	return products, nil
}
