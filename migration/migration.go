package migration

import (
	"github.com/open-cmi/migrate"
)

func Migrate() {
	migrate.Init("cmmns")

	migrate.Run()
}
