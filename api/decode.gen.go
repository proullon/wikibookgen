// Code generated by protoc-gen-gotemplate
package wikibookgen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

var _uselessIO = io.EOF
var _uselessJSON = json.Valid(nil)

var ErrBadRequest = fmt.Errorf("Bad Request")

// httpDecoder is used to decode http request
// If needed, default decoder function can be overwritten
// by implementing it directly as httpDecoder method
type httpDecoder struct {
	defaultDecoder
}

type defaultDecoder struct {
}

func (d *defaultDecoder) DecodeStatusRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &Void{}

	return req, nil
}

func (d *defaultDecoder) DecodeCompleteRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &CompleteRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("POST /complete: %s\n", err)
		return nil, ErrBadRequest
	}

	if stringvar, ok := vars["value"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Value); err != nil {
			return nil, fmt.Errorf("%s: %s", "value", err)
		}
	} else if stringvar := r.FormValue("value"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Value); err != nil {
			return nil, fmt.Errorf("%s: %s", "value", err)
		}
	}

	if stringvar, ok := vars["language"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Language); err != nil {
			return nil, fmt.Errorf("%s: %s", "language", err)
		}
	} else if stringvar := r.FormValue("language"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Language); err != nil {
			return nil, fmt.Errorf("%s: %s", "language", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodeOrderRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &OrderRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("POST /order: %s\n", err)
		return nil, ErrBadRequest
	}

	if stringvar, ok := vars["subject"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Subject); err != nil {
			return nil, fmt.Errorf("%s: %s", "subject", err)
		}
	} else if stringvar := r.FormValue("subject"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Subject); err != nil {
			return nil, fmt.Errorf("%s: %s", "subject", err)
		}
	}

	if stringvar, ok := vars["model"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Model); err != nil {
			return nil, fmt.Errorf("%s: %s", "model", err)
		}
	} else if stringvar := r.FormValue("model"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Model); err != nil {
			return nil, fmt.Errorf("%s: %s", "model", err)
		}
	}

	if stringvar, ok := vars["language"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Language); err != nil {
			return nil, fmt.Errorf("%s: %s", "language", err)
		}
	} else if stringvar := r.FormValue("language"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Language); err != nil {
			return nil, fmt.Errorf("%s: %s", "language", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodeOrderStatusRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &OrderStatusRequest{}

	if stringvar, ok := vars["uuid"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	} else if stringvar := r.FormValue("uuid"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodeListWikibookRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &ListWikibookRequest{}

	if stringvar, ok := vars["page"]; ok {
		if err := convertTYPE_INT64(stringvar, &req.Page); err != nil {
			return nil, fmt.Errorf("%s: %s", "page", err)
		}
	} else if stringvar := r.FormValue("page"); stringvar != "" {
		if err := convertTYPE_INT64(stringvar, &req.Page); err != nil {
			return nil, fmt.Errorf("%s: %s", "page", err)
		}
	}

	if stringvar, ok := vars["size"]; ok {
		if err := convertTYPE_INT64(stringvar, &req.Size); err != nil {
			return nil, fmt.Errorf("%s: %s", "size", err)
		}
	} else if stringvar := r.FormValue("size"); stringvar != "" {
		if err := convertTYPE_INT64(stringvar, &req.Size); err != nil {
			return nil, fmt.Errorf("%s: %s", "size", err)
		}
	}

	if stringvar, ok := vars["language"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Language); err != nil {
			return nil, fmt.Errorf("%s: %s", "language", err)
		}
	} else if stringvar := r.FormValue("language"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Language); err != nil {
			return nil, fmt.Errorf("%s: %s", "language", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodeGetWikibookRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &GetWikibookRequest{}

	if stringvar, ok := vars["uuid"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	} else if stringvar := r.FormValue("uuid"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodeGetAvailableDownloadFormatRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &GetAvailableDownloadFormatRequest{}

	if stringvar, ok := vars["uuid"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	} else if stringvar := r.FormValue("uuid"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodeDownloadWikibookRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &DownloadWikibookRequest{}

	if stringvar, ok := vars["uuid"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	} else if stringvar := r.FormValue("uuid"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	}

	if stringvar, ok := vars["format"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Format); err != nil {
			return nil, fmt.Errorf("%s: %s", "format", err)
		}
	} else if stringvar := r.FormValue("format"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Format); err != nil {
			return nil, fmt.Errorf("%s: %s", "format", err)
		}
	}

	return req, nil
}

func (d *defaultDecoder) DecodePrintWikibookRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	_ = vars
	r.ParseForm()

	req := &PrintWikibookRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && err != io.EOF {
		log.Errorf("POST /wikibook/{uuid}/print/{format}: %s\n", err)
		return nil, ErrBadRequest
	}

	if stringvar, ok := vars["uuid"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	} else if stringvar := r.FormValue("uuid"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Uuid); err != nil {
			return nil, fmt.Errorf("%s: %s", "uuid", err)
		}
	}

	if stringvar, ok := vars["format"]; ok {
		if err := convertTYPE_STRING(stringvar, &req.Format); err != nil {
			return nil, fmt.Errorf("%s: %s", "format", err)
		}
	} else if stringvar := r.FormValue("format"); stringvar != "" {
		if err := convertTYPE_STRING(stringvar, &req.Format); err != nil {
			return nil, fmt.Errorf("%s: %s", "format", err)
		}
	}

	return req, nil
}

func convertTYPE_DOUBLE(in string, out *float64) error {
	f, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}

	*out = f
	return nil
}

func convertTYPE_BOOL(in string, out *bool) error {
	b, err := strconv.ParseBool(in)
	if err != nil {
		return err
	}

	*out = b
	return nil
}

func convertTYPE_STRING(in string, out *string) error {
	*out = in
	return nil
}

func convertTYPE_INT64(in string, out *int64) error {
	n, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return fmt.Errorf("got '%s': %s", in, err)
	}

	*out = n
	return nil
}
