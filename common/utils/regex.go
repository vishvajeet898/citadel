package utils

import "regexp"

var (
	NumericResultTypeRegex          = regexp.MustCompile(`^\d+(\.\d+)?$`)
	TextualResultTypeRegex          = regexp.MustCompile(`^[a-zA-Z0-9\(\)\.\+\-:&_,%"<>=\s,/;#@\[\]{}!?’‘“”'&nbsp;]+$`)
	SemiQuantitativeResultTypeRegex = regexp.MustCompile(`^(?:[<>]=?|=)?\s*\d+(?:\.\d+)?(?:[:-]\d+)?$`)
)
