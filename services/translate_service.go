package services

import (
	"github.com/OwO-Network/gdeeplx"
	"github.com/abadojack/whatlanggo"
	"github.com/xiaoxuan6/deeplx"
	"regexp"
	"strings"
)

func Translate(text string) string {
	if len(text) < 1 {
		return ""
	}

	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, "-")
	text = strings.ReplaceAll(text, "|", "")

	lang := whatlanggo.DetectLang(text)
	deepLLang := strings.ToUpper(lang.Iso6391())
	if strings.Contains(deepLLang, "ZH") {
		return text
	}

	response := deeplx.Translate(text, deepLLang, "zh")
	if response.Code == 200 {
		return response.Data
	}

	result, err := gdeeplx.Translate(text, deepLLang, "ZH", 0)
	if err == nil {
		return result.(map[string]interface{})["data"].(string)
	}

	return text
}
