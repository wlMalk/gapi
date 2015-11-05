package defaults

import (
	"github.com/wlMalk/gapi/constants"
)

var (
	CacheReadEntityBytes = true

	MimeType = constants.MIME_JSON

	MimeInAccept = true

	Indent = true

	MaxMemory int64 = 32 << 20 // 32 MB

	Schemes  = []string{constants.SCHEME_HTTP, constants.SCHEME_HTTPS}
	Consumes = []string{constants.MIME_JSON}
	Produces = []string{constants.MIME_JSON}
)
