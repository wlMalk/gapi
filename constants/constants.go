package constants

const (
	TYPE_STRING = iota
	TYPE_INT
	TYPE_INT64
	TYPE_FLOAT
	TYPE_FLOAT64
	TYPE_BOOL

	IN_PATH = iota
	IN_QUERY
	IN_HEADER
	IN_BODY

	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PUT    = "PUT"
	METHOD_DELETE = "DELETE"

	SCHEME_HTTP  = "http"
	SCHEME_HTTPS = "https"

	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"

	HEADER_Allow                         = "Allow"
	HEADER_Accept                        = "Accept"
	HEADER_Origin                        = "Origin"
	HEADER_ContentType                   = "Content-Type"
	HEADER_LastModified                  = "Last-Modified"
	HEADER_AcceptEncoding                = "Accept-Encoding"
	HEADER_ContentEncoding               = "Content-Encoding"
	HEADER_AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HEADER_AccessControlRequestMethod    = "Access-Control-Request-Method"
	HEADER_AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HEADER_AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HEADER_AccessControlMaxAge           = "Access-Control-Max-Age"

	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"
)
