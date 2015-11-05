package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/wlMalk/gapi/constants"
	"github.com/wlMalk/gapi/defaults"
)

func WriteEntity(w http.ResponseWriter, status int, value interface{}, encoding string, indent bool) (int, error) {
	if value == nil {
		return 0, nil
	}

	switch encoding {
	case constants.MIME_JSON:
		return WriteAsJson(w, status, value, indent)
	case constants.MIME_XML:
		return WriteAsXml(w, status, value, true, indent)
	}

	return 0, nil
}

func WriteAsXml(w http.ResponseWriter, status int, value interface{}, writeHeader bool, indent bool) (int, error) {
	var output []byte
	var err error

	if value == nil {
		return 0, nil
	}
	if indent {
		output, err = xml.MarshalIndent(value, " ", " ")
	} else {
		output, err = xml.Marshal(value)
	}

	if err != nil {
		return WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set(constants.HEADER_ContentType, constants.MIME_XML)
	w.WriteHeader(status)
	if writeHeader {
		cl, err := w.Write([]byte(xml.Header))
		if err != nil {
			return cl, err
		}
	}
	return w.Write(output)

}

func WriteAsJson(w http.ResponseWriter, status int, value interface{}, indent bool) (int, error) {
	return WriteJson(w, status, value, constants.MIME_JSON, indent)
}

func WriteJson(w http.ResponseWriter, status int, value interface{}, contentType string, indent bool) (int, error) {
	var output []byte
	var err error

	if value == nil {
		return 0, nil
	}
	if indent {
		output, err = json.MarshalIndent(value, " ", " ")
	} else {
		output, err = json.Marshal(value)
	}

	if err != nil {
		return WriteString(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set(constants.HEADER_ContentType, contentType)
	w.WriteHeader(status)
	return w.Write(output)
}

func WriteError(w http.ResponseWriter, status int, e *Error) (int, error) {
	return WriteEntity(w, status, e, defaults.MimeType, defaults.Indent)
}

func WriteString(w http.ResponseWriter, status int, str string) (int, error) {
	w.WriteHeader(status)
	return w.Write([]byte(str))
}
