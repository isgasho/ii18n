package ii18n

import (
	"strings"
	"errors"
)

type TMsgs map[string]string

// Source interface
type Source interface {
	Translate(category string, message string, lang string) (string, error)
	TranslateMsg(category string, message string, lang string) (string, error)
	GetMsgFilePath(category string, lang string) string
	LoadMsgs(category string, lang string) (TMsgs, error)
	LoadFallbackMsgs(category string, fallbackLang string, msgs TMsgs, originalMsgFile string) (TMsgs, error)
}

// MessageSource
type MessageSource struct {
	SourceLang       string
	ForceTranslation bool
	BasePath         string
	FileMap          map[string]string
	fileSuffix       string
	loadFunc         func(filename string) (TMsgs, error)
	messages         map[string]TMsgs
}

// translate
func (ms *MessageSource) Translate(category string, message string, lang string) (string, error) {
	if ms.ForceTranslation || lang != ms.SourceLang {
		return ms.TranslateMsg(category, message, lang)
	}
	return "", nil
}

// translate
func (ms *MessageSource) TranslateMsg(category string, message string, lang string) (string, error) {
	cates := strings.Split(category, ".")
	key := cates[0] + "/" + lang + "/" + cates[1]
	if _, ok := ms.messages[key]; !ok {
		val, err := ms.LoadMsgs(category, lang)
		if err != nil {
			return "", err
		}
		ms.messages[key] = val
	}
	if msg, ok := ms.messages[key][message]; ok && msg != "" {
		return msg, nil
	}

	ms.messages[key] = TMsgs{message: ""}
	return "", nil
}

// Get messages file path.
func (ms *MessageSource) GetMsgFilePath(category string, lang string) string {
	suffix := strings.Split(category, ".")[1]
	path := ms.BasePath + "/" + lang + "/"
	if v, ok := ms.FileMap[suffix]; !ok {
		path += v
	} else {
		path += strings.Replace(suffix, "\\", "/", -1)
		if ms.fileSuffix != "" {
			path += "." + ms.fileSuffix
		}
	}
	return path
}

// Loads the message translation for the specified $language and $category.
// If translation for specific locale code such as `en-US` isn't found it
// tries more generic `en`. When both are present, the `en-US` messages will be merged
// over `en`. See [[loadFallbackTMsgs]] for details.
// If the lang is less specific than [[sourceLang]], the method will try to
// load the messages for [[sourceLang]]. For example: [[sourceLang]] is `en-GB`,
// language is `en`. The method will load the messages for `en` and merge them over `en-GB`.
func (ms *MessageSource) LoadMsgs(category string, lang string) (TMsgs, error) {
	msgFile := ms.GetMsgFilePath(category, lang)
	msgs, err := ms.loadFunc(msgFile)
	if err != nil {
		return nil, err
	}
	fallbackLang := lang[0:2]
	fallbackSourceLang := ms.SourceLang[0:2]
	if lang != fallbackLang {
		msgs, err = ms.LoadFallbackMsgs(category, fallbackLang, msgs, msgFile)
	} else if lang == fallbackSourceLang {
		msgs, err = ms.LoadFallbackMsgs(category, ms.SourceLang, msgs, msgFile)
	} else {
		if msgs == nil {
			return nil, errors.New("the message file for category " + category + " does not exist: " + msgFile)
		}
	}
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// Loads the message translation for the specified $language and $category.
// If translation for specific locale code such as `en-US` isn't found it
// tries more generic `en`. When both are present, the `en-US` messages will be merged
func (ms *MessageSource) LoadFallbackMsgs(category string, fallbackLang string, msgs TMsgs, originalMsgFile string) (TMsgs, error) {
	fallbackMsgFile := ms.GetMsgFilePath(category, fallbackLang)
	fallbackMsgs, _ := ms.loadFunc(fallbackMsgFile)
	if msgs == nil && fallbackMsgs == nil &&
		fallbackLang != ms.SourceLang &&
		fallbackLang != ms.SourceLang[0:2] {
		return nil, errors.New("The message file for category " + category + " does not exist: " + originalMsgFile + " Fallback file does not exist as well: " + fallbackMsgFile)
	} else if msgs == nil {
		return fallbackMsgs, nil
	} else if fallbackMsgs != nil {
		for key, val := range fallbackMsgs {
			v, ok := msgs[key]
			if val != "" && (!ok || v == "") {
				msgs[key] = fallbackMsgs[key]
			}
		}
	}

	return msgs, nil
}

// Get messages file path.
func LoadMsgsFromFile(filename string) (TMsgs, error) {
	return nil, nil
}
