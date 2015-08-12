package middleware

import (
	"net/http"

	"github.com/wlMalk/gapi/constants"
	"github.com/wlMalk/gapi/internal/util"
	"github.com/wlMalk/gapi/param"
	"github.com/wlMalk/gapi/request"
	"github.com/wlMalk/gapi/validation"
)

var (
	MaxMemory int64 = 32 << 20 // 32 MB
)

func CheckSchemes(schemes []string) MiddlewareFunc {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *request.Request) {

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

func ValidateParams(params *param.Params) MiddlewareFunc {
	return MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *request.Request) {

			q := r.URL.Query()
			h := r.Header
			pathParams, ok := r.Attributes["path_params"].(map[string]string)
			if !ok {
				panic("Path params not set correctly.")
			}

			if params.ContainsFiles() {
				r.ParseMultipartForm(MaxMemory)
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
