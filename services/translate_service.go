package services

import (
	"github.com/OwO-Network/gdeeplx"
	"github.com/abadojack/whatlanggo"
	"github.com/xiaoxuan6/deeplx"
	"strings"
)

func Translate(text string) string {
	lang := whatlanggo.DetectLang(text)
	deepLLang := strings.ToUpper(lang.Iso6391())

	if strings.Contains(deepLLang, "ZH") {
		return text
	}

	result, err := gdeeplx.Translate(text, deepLLang, "ZH", 0)
	if err != nil {
		response := deeplx.Translate(text, deepLLang, "zh")

		if response.Code == 200 {
			return response.Data
		}

		return text
	}

	return result.(map[string]interface{})["data"].(string)
}
