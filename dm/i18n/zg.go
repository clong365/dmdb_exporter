/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package i18n

import (
	"encoding/json"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
)

type msg struct {
	Id          string `json:"id"`
	Translation string `json:"translation,omitempty"`
}

type i18n struct {
	Language string `json:"language"`
	Messages []msg  `json:"messages"`
}

func readI18nJson(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str

}

func InitConfig(jsonPath string) {

	var i18n i18n
	//str := readI18nJson(jsonPath)
	json.Unmarshal([]byte(jsonPath), &i18n)
	msaArry := i18n.Messages
	tag := language.MustParse(i18n.Language)
	for _, e := range msaArry {
		message.SetString(tag, e.Id, e.Translation)
	}
}

func Get(key string, locale int) string {
	var p *message.Printer

	switch locale {
	case 0:
		p = message.NewPrinter(language.SimplifiedChinese)
	case 1:
		p = message.NewPrinter(language.AmericanEnglish)
	case 2:
		p = message.NewPrinter(language.TraditionalChinese)
	}

	return p.Sprintf(key)
}
