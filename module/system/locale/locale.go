package locale

import (
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/initial"
)

func SetLocale(lang string) error {
	m := Get()
	if m == nil {
		m = New()
	}
	m.Lang = lang
	err := i18n.SetLang(lang)
	if err != nil {
		return err
	}
	return m.Save()
}

func GetLocale() string {
	m := Get()
	if m == nil {
		return i18n.GetLang()
	}
	return m.Lang
}

func Init() error {
	m := Get()
	if m == nil {
		i18n.SetLang("zh-CN")
		return nil
	}
	i18n.SetLang(m.Lang)
	return nil
}

func init() {
	initial.Register("locale", initial.PhaseDefault, Init)
}
