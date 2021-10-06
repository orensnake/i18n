package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var Translation *TTranslation

type fTranslation struct {
	Lang []struct {
		Name string `json:"name"`
		Dict []struct {
			Id  int    `json:"id"`
			Txt string `json:"txt"`
		} `json:"dict"`
	} `json:"lang"`
}

type TTranslation struct {
	Lang  string
	trans fTranslation
}

func (t *TTranslation) getLocale() (string, error) {
	envLang, fnd := os.LookupEnv("LANG")
	if fnd {
		return strings.Split(envLang, ".")[0], nil
	} else {
		return "en_US", nil
	}
}

func (t *TTranslation) loadDict(fName string) fTranslation {
	var dict fTranslation
	configFile, err := os.Open(fName)
	defer configFile.Close()
	if err != nil {
		fmt.Println("Error opening translation file ", err.Error())
		os.Exit(1)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&dict)
	return dict
}

func (t *TTranslation) Init(translationFile string) {
	t.Lang, _ = t.getLocale()
	t.trans = t.loadDict(translationFile)
	// Check system language in translations
	fnd := false
	for i := 0; i < len(t.trans.Lang); i++ {
		if t.Lang == t.trans.Lang[i].Name {
			fnd = true
			break
		}
	}
	if !fnd {
		fmt.Println("Translation to", t.Lang, "not found.")
		if t.trans.Lang != nil {
			t.Lang = t.trans.Lang[0].Name
			fmt.Println("Using", t.Lang, "as default")
		} else {
			fmt.Println("Translations not found. Exiting.")
			os.Exit(-1)
		}
	}
}

func (t TTranslation) SetLang(lang string) {
	t.Lang = lang
}

func (t *TTranslation) GetText(id int) string {
	res := ""
	for l := 0; l < len(t.trans.Lang); l++ {
		if t.trans.Lang[l].Name == t.Lang {
			// Есть перевод с необходимоо языка
			for i := 0; i < len(t.trans.Lang[l].Dict); i++ {
				if t.trans.Lang[l].Dict[i].Id == id {
					res = t.trans.Lang[l].Dict[i].Txt
					return res
				}
			}
		}
	}
	return fmt.Sprintf("???unknown message with id %v ???", id)
}
