package articles

import (
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
	"regexp"
)

var uriFormatRe = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

type URI struct {
	s string
}

func NewUri(s string) (URI, error) {
	if !uriFormatRe.MatchString(s) {
		return URI{}, errors.NewDomainError("uri只能包含字母、数字、下划线或连字符")
	}
	return URI{s: s}, nil
}

func (u URI) String() string {
	return u.s
}
