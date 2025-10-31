package regex_s

import (
	"regexp"
)

func RegexRecognize(pattern, str string) (bool, error) {
	match, err := regexp.Match(pattern, []byte(str))
	if !match {
		return match, err
	}
	return match, nil
}
