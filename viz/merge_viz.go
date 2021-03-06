package viz

import (
	"strings"
)

var MergeHeaderFunc = func(input string) string {
	tmp := strings.Split(input, ".")
	if len(tmp) > 1 {
		return strings.Join(tmp[0:len(tmp)-1], ".")
	}
	return input
}

var MergePackageFunc = func(input string) string {
	split := "/"
	if !strings.Contains(input, split) {
		split = "."
	}
	if !strings.Contains(input, split) {
		split = "::"
	}
	tmp := strings.Split(input, split)
	packageName := tmp[0]
	if packageName == input {
		packageName = "main"
	}
	if len(tmp) > 2 {
		packageName = strings.Join(tmp[0:len(tmp)-1], split)
	}

	return packageName
}
