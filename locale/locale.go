package locale

import (
	nativeLog "log"
)

type Language string

const (
	English Language = "en"
	Persian Language = "fa"
	Arabic  Language = "ar"
	Turkish Language = "tr"
)

type Locale interface {
	Message(lang Language, key string, params ...interface{}) string
}

var loc Locale

func Load(dir string) {
	// Read locale
	l, err := NewFromDir(dir)
	if err != nil {
		nativeLog.Printf("error (%v) on setup locale_keys service", err)
		return
	}
	loc = l
}

func Message(lang string, key string) string {
	if loc != nil {
		return loc.Message(Language(lang), key)
	}
	return ""
}
