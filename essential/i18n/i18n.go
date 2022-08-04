package i18n

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printerMapping map[string]*message.Printer = make(map[string]*message.Printer)

var supportedLang []string = []string{"en-US", "zh-CN"}

type Config struct {
	Lang string `json:"lang"`
}

var gConf Config

func GetLang() string {
	return gConf.Lang
}

func ChangeLang(lang string) {
	gConf.Lang = lang
	config.Save()
}

// Sprint is like fmt.Sprint, but using language-specific formatting.
func Sprint(args ...interface{}) string {
	return printerMapping[gConf.Lang].Sprint(args)
}

// Print is like fmt.Print, but using language-specific formatting.
func Print(args ...interface{}) (n int, err error) {
	return printerMapping[gConf.Lang].Print(args)
}

// Sprintln is like fmt.Sprintln, but using language-specific formatting.
func Sprintln(args ...interface{}) string {
	return printerMapping[gConf.Lang].Sprintln(args)
}

// Println is like fmt.Println, but using language-specific formatting.
func Println(args ...interface{}) (n int, err error) {
	return printerMapping[gConf.Lang].Println(args)
}

// Sprintf is like fmt.Sprintf, but using language-specific formatting.
func Sprintf(format string, args ...interface{}) string {
	return printerMapping[gConf.Lang].Sprintf(format, args)
}

// Printf is like fmt.Printf, but using language-specific formatting.
func Printf(format string, args ...interface{}) (n int, err error) {
	return printerMapping[gConf.Lang].Printf(format, args)
}

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	var found bool = false
	for _, lang := range supportedLang {
		if lang == gConf.Lang {
			found = true
		}
		tag := language.MustParse(lang)
		p := message.NewPrinter(tag)
		printerMapping[lang] = p
	}
	if !found {
		gConf.Lang = supportedLang[0]
	}
	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(gConf)
	return raw
}

func init() {
	gConf.Lang = "en-US"
	config.RegisterConfig("locale", Init, Save)
}
