package translation

import (
	_ "embed"

	"github.com/open-cmi/cmmns/pkg/translation"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

//go:embed locales/en-US/out.gotext.json
var usMessages string

//go:embed locales/zh-CN/out.gotext.json
var cnMessages string

func init() {
	usDict, err := translation.InitTransDict(usMessages)
	if err != nil {
		panic(err)
	}

	cnDict, err := translation.InitTransDict(cnMessages)
	if err != nil {
		panic(err)
	}

	dict := map[string]catalog.Dictionary{
		"en_US": usDict,
		"zh_CN": cnDict,
	}
	fallback := language.MustParse("en-US")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}
