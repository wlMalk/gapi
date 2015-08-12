package wrapper

import (
	"strings"
)

const PathParamsKey = "path_params"

func CurlyToColon(path string) string {
	path = strings.Replace(path, "{", ":", -1)
	path = strings.Replace(path, "}", "", -1)
	return path
}
