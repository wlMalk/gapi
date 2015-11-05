package middleware

import (
	"net/http"
	"strings"

	"github.com/wlMalk/gapi/constants"
	"github.com/wlMalk/gapi/defaults"
	"github.com/wlMalk/gapi/internal/util"
	"github.com/wlMalk/gapi/param"
	"github.com/wlMalk/gapi/request"
	"github.com/wlMalk/gapi/response"
	"github.com/wlMalk/gapi/validation"
)

func CheckSchemes(schemes []string) MiddlewareFunc {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(w *response.Response, r *request.Request) {

			var schemeAccepted bool
			if r.URL.Scheme == "" {
				r.URL.Scheme = constants.SCHEME_HTTP
			}

			schemeAccepted = util.ContainsString(schemes, r.URL.Scheme)
			if !schemeAccepted {
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}

func getEncodingFromAccept(produces []string, r *request.Request) (string, *response.Error) {
	var encoding string

	for _, acceptMime := range strings.Split(r.Header.Get(constants.HEADER_Accept), ",") {
		mime := strings.Trim(strings.Split(acceptMime, ";")[0], " ")
		if 0 == len(mime) || mime == "*/*" {
			if len(produces) == 0 {
				encoding = defaults.MimeType
				break
			} else {
				encoding = produces[0]
				break
			}
		} else {
			if util.ContainsString(produces, mime) {
				encoding = mime
				break
			}
		}
	}

	if len(encoding) == 0 {
		return encoding, response.NewError(406, "Request", "Encoding requested not valid.")
	}

	return encoding, nil
}

func CheckProduces(produces []string, mimeInAccept bool) MiddlewareFunc {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(w *response.Response, r *request.Request) {
			if mimeInAccept {
				enc, err := getEncodingFromAccept(produces, r)
				if err != nil {
					w.WriteError(http.StatusBadRequest, err)
				}
				w.Encoding = enc
			} else {

			}

			next.ServeHTTP(w, r)
		})
	})
}

func ValidateParams(params *param.Params) MiddlewareFunc {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(w *response.Response, r *request.Request) {

			q := r.URL.Query()
			h := r.Header
			pathParams, ok := r.Attributes["path_params"].(map[string]string)
			if !ok {
				panic("Path params not set correctly.")
			}

			if params.ContainsFiles() {
				r.ParseMultipartForm(defaults.MaxMemory)
			} else if params.ContainsBodyParams() {
				r.ParseForm()
			}
			for _, p := range params.Get() {
				var pv *validation.Value
				if p.IsIn(constants.IN_PATH) {

					v, ok := pathParams[p.GetName()]
					if !ok {
						return
					}
					pv = validation.NewValue(p.GetName(), v, "path", p.GetAs())
				} else if p.IsIn(constants.IN_QUERY) {
					v, ok := q[p.GetName()]
					if !ok {
						if p.IsRequired() {
							return
						}
					} else {
						if !p.IsMultiple() {
							pv = validation.NewValue(p.GetName(), v[0], "query", p.GetAs())
						} else {
							pv = validation.NewMultipleValue(p.GetName(), v, "query", p.GetAs())
						}
					}
				} else if p.IsIn(constants.IN_HEADER) {
					v, ok := h[p.GetName()]
					if !ok {
						if p.IsRequired() {
							return
						}
					} else {
						if !p.IsMultiple() {
							pv = validation.NewValue(p.GetName(), v[0], "header", p.GetAs())
						} else {
							pv = validation.NewMultipleValue(p.GetName(), v, "header", p.GetAs())
						}
					}
				} else if p.IsIn(constants.IN_BODY) { // decide what to do when content type is form-encoded
					if p.IsFile() {
						_, ok := r.MultipartForm.File[p.GetName()]
						if !ok {
							if p.IsRequired() {
								return
							}
						} else {
							//pv = NewFileParamValue(p.name, v[0], "header")
						}
					} else if !p.IsFile() && params.ContainsFiles() {
						_, ok := r.MultipartForm.Value[p.GetName()]
						if !ok {
							if p.IsRequired() {
								return
							}
						} else {
							//pv = NewFileParamValue(p.name, v[0], "header")
						}
					} else {
						v, ok := r.Form[p.GetName()]
						if !ok {
							if p.IsRequired() {
								return
							}
						} else {
							if !p.IsMultiple() {
								pv = validation.NewValue(p.GetName(), v[0], "body", p.GetAs())
							} else {
								pv = validation.NewMultipleValue(p.GetName(), v, "body", p.GetAs())
							}
						}
					}
				}
				if pv != nil {
					r.Input[p.GetName()] = pv
				}
			}
			for _, p := range params.Get() {
				pv, ok := r.Input[p.GetName()]
				if ok {
					err := p.Validate(pv, r)
					if err != nil {
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	})
}
