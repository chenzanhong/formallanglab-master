package regex_s

import (
	"backend/internal/domain/model"
	"regexp"
)

func RegexRecognize(pattern model.Regex, str string) (bool, error) {
	match, err := regexp.Match(string(pattern), []byte(str))
	if !match {
		return match, err
	}
	return match, nil
}
