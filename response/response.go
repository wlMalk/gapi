package response

import (
	"net/http"
)

type Response struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
	Encoding      string
	indent        bool
}

func New(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
	}
}

func (res *Response) WriteEntity(status int, value interface{}) error {
	_, err := WriteEntity(res, status, value, res.Encoding, res.indent)
	return err
}

func (res *Response) WriteAsXml(status int, value interface{}, writeHeader bool) error {
	_, err := WriteAsXml(res, status, value, writeHeader, res.indent)
	return err
}

func (res *Response) WriteAsJson(status int, value interface{}) error {
	_, err := WriteAsJson(res, status, value, res.indent)
	return err
}

func (res *Response) WriteJson(status int, value interface{}, contentType string) error {
	_, err := WriteJson(res, status, value, contentType, res.indent)
	return err
}

func (res *Response) WriteError(status int, e *Error) error {
	err := res.WriteEntity(status, e)
	return err
}

func (res *Response) WriteString(status int, str string) error {
	_, err := WriteString(res, status, str)
	return err
}

func (res *Response) WriteHeader(httpStatus int) {
	if res.statusCode == 0 {
		if httpStatus == 0 {
			httpStatus = http.StatusOK
		}
		res.statusCode = httpStatus
		res.ResponseWriter.WriteHeader(httpStatus)
	}
}

func (res *Response) StatusCode() int {
	if res.statusCode == 0 {
		return http.StatusOK
	}
	return res.statusCode
}

//
// func (res *Response) WriteString(str string) (int, error) {
// 	written, err := res.ResponseWriter.Write([]byte(str))
// 	res.contentLength += written
// 	return written, err
// }

func (res *Response) Write(bytes []byte) (int, error) {
	written, err := res.ResponseWriter.Write(bytes)
	res.contentLength += written
	return written, err
}

func (res *Response) Indented(indent bool) {
	res.indent = indent
}

func (res *Response) ContentLength() int {
	return res.contentLength
}

func (res *Response) Flush() {
	if f, ok := res.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (res *Response) CloseNotify() <-chan bool {
	return res.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
