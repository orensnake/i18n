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
	//t.Trans = new(fTranslation)
	t.trans = t.loadDict(translationFile)
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
