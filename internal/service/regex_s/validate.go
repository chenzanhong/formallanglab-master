package regex_s

import (
	"fmt"
	"regexp"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func RegexValidate(regex model.Regex) error {
	return regex.IsValid()
}

var alphaNumRegexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

func RegexValidString(str string) (bool, error) {
	if str == string(model.Epsilon) || str == string(model.EmptyLanguageToken) {
		return true, nil
	}

	if !alphaNumRegexp.MatchString(str) {
		return false, fmt.Errorf("字符串只能包含字母和数字")
	}

	return true, nil
}
