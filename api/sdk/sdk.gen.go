// Code generated by protoc-gen-gotemplate
package wikibookgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type businessError struct {
	CodeError string `json:"code_error"`
	Message   string `json:"message"`
	Root      string `json:"root"`
}

type Mode int

const (
	unset Mode = iota
	HTTP
	GRPC
)

var (
	mode     Mode
	endpoint string
	token    string
)

func Initialize(e string, t string) {
	endpoint = e
	token = t
	mode = HTTP
}

func request(method, path string, in, out interface{}) error {
	baseAddress := endpoint

	resp, err := execrequest(baseAddress, token, method, path, in)
	if err != nil {
		return err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}

func execrequest(baseAddress, token, method, path string, in interface{}) (io.ReadCloser, error) {
	var req *http.Request
	var err error

	if in != nil {
		args, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, baseAddress+path, bytes.NewReader(args))
	} else {
		req, err = http.NewRequest(method, baseAddress+path, nil)
	}

	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var berr businessError
			jerr := json.Unmarshal(body, &berr)
			if jerr == nil && berr.CodeError != "" {
				return nil, fmt.Errorf("[%s] %s (%s)", berr.CodeError, berr.Message, berr.Root)
			}
		}
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// Status : GET /status
func Status() (status []string, err error) {

	in := &Void{}

	var out *StatusResponse
	_ = out

	switch mode {
	case HTTP:
		out, err = httpStatus(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return out.Status, nil

}

// Complete : POST /complete
func Complete(value_ string, language_ string) (titles []string, err error) {

	in := &CompleteRequest{
		Value:    value_,
		Language: language_,
	}

	var out *CompleteResponse
	_ = out

	switch mode {
	case HTTP:
		out, err = httpComplete(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return out.Titles, nil

}

// Order : POST /order
func Order(subject_ string, model_ string, language_ string) (uuid string, err error) {

	in := &OrderRequest{
		Subject:  subject_,
		Model:    model_,
		Language: language_,
	}

	var out *OrderResponse
	_ = out

	switch mode {
	case HTTP:
		out, err = httpOrder(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return out.Uuid, nil

}

// OrderStatus : GET /order/{id}
func OrderStatus(uuid_ string) (status string, wikibookUuid string, err error) {

	in := &OrderStatusRequest{
		Uuid: uuid_,
	}

	var out *OrderStatusResponse
	_ = out

	switch mode {
	case HTTP:
		out, err = httpOrderStatus(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return out.Status, out.WikibookUuid, nil

}

// ListWikibook : GET /wikibook?page={page}&size={size}&language={language}
func ListWikibook(page_ int64, size_ int64, language_ string) (wikibooks []*Wikibook, err error) {

	in := &ListWikibookRequest{
		Page:     page_,
		Size:     size_,
		Language: language_,
	}

	var out *ListWikibookResponse
	_ = out

	switch mode {
	case HTTP:
		out, err = httpListWikibook(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return out.Wikibooks, nil

}

// GetWikibook : GET /wikibook/{id}
func GetWikibook(uuid_ string) (wikibook *Wikibook, err error) {

	in := &GetWikibookRequest{
		Uuid: uuid_,
	}

	var out *GetWikibookResponse
	_ = out

	switch mode {
	case HTTP:
		out, err = httpGetWikibook(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return out.Wikibook, nil

}

// DownloadWikibook : GET /wikibook/{id}/download
func DownloadWikibook(uuid_ string, format_ string) (err error) {

	in := &DownloadWikibookRequest{
		Uuid:   uuid_,
		Format: format_,
	}

	var out *Void
	_ = out

	switch mode {
	case HTTP:
		out, err = httpDownloadWikibook(in, endpoint, token)
	}
	if err != nil {
		return
	}

	return nil

}

// ------------------------- HTTP SDK -----------------------------

// Status : GET /status
func httpStatus(in *Void, baseAddress, token string) (out *StatusResponse, err error) {

	path := "/status"
	placeholders := map[string]string{}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &StatusResponse{}
	err = request("GET", path, in, out)
	return
}

// Complete : POST /complete
func httpComplete(in *CompleteRequest, baseAddress, token string) (out *CompleteResponse, err error) {

	path := "/complete"
	placeholders := map[string]string{
		"value":    fmt.Sprintf("%v", in.Value),
		"language": fmt.Sprintf("%v", in.Language),
	}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &CompleteResponse{}
	err = request("POST", path, in, out)
	return
}

// Order : POST /order
func httpOrder(in *OrderRequest, baseAddress, token string) (out *OrderResponse, err error) {

	path := "/order"
	placeholders := map[string]string{
		"subject":  fmt.Sprintf("%v", in.Subject),
		"model":    fmt.Sprintf("%v", in.Model),
		"language": fmt.Sprintf("%v", in.Language),
	}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &OrderResponse{}
	err = request("POST", path, in, out)
	return
}

// OrderStatus : GET /order/{id}
func httpOrderStatus(in *OrderStatusRequest, baseAddress, token string) (out *OrderStatusResponse, err error) {

	path := "/order/{id}"
	placeholders := map[string]string{
		"uuid": fmt.Sprintf("%v", in.Uuid),
	}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &OrderStatusResponse{}
	err = request("GET", path, in, out)
	return
}

// ListWikibook : GET /wikibook?page={page}&size={size}&language={language}
func httpListWikibook(in *ListWikibookRequest, baseAddress, token string) (out *ListWikibookResponse, err error) {

	path := "/wikibook?page={page}&size={size}&language={language}"
	placeholders := map[string]string{
		"page":     fmt.Sprintf("%v", in.Page),
		"size":     fmt.Sprintf("%v", in.Size),
		"language": fmt.Sprintf("%v", in.Language),
	}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &ListWikibookResponse{}
	err = request("GET", path, in, out)
	return
}

// GetWikibook : GET /wikibook/{id}
func httpGetWikibook(in *GetWikibookRequest, baseAddress, token string) (out *GetWikibookResponse, err error) {

	path := "/wikibook/{id}"
	placeholders := map[string]string{
		"uuid": fmt.Sprintf("%v", in.Uuid),
	}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &GetWikibookResponse{}
	err = request("GET", path, in, out)
	return
}

// DownloadWikibook : GET /wikibook/{id}/download
func httpDownloadWikibook(in *DownloadWikibookRequest, baseAddress, token string) (out *Void, err error) {

	path := "/wikibook/{id}/download"
	placeholders := map[string]string{
		"uuid":   fmt.Sprintf("%v", in.Uuid),
		"format": fmt.Sprintf("%v", in.Format),
	}
	for k, v := range placeholders {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	out = &Void{}
	err = request("GET", path, in, out)
	return
}

type Wikibook struct {
	Uuid     string    `json:"uuid"`
	Subject  string    `json:"subject"`
	Model    string    `json:"model"`
	Language string    `json:"language"`
	Title    string    `json:"title"`
	Pages    int64     `json:"pages"`
	Volumes  []*Volume `json:"volumes"`
}

type Volume struct {
	Title    string     `json:"title"`
	Chapters []*Chapter `json:"chapters"`
}

type Chapter struct {
	Title    string  `json:"title"`
	Articles []*Page `json:"articles"`
}

type Page struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type StatusResponse struct {
	Status []string `json:"status"`
}

type Void struct {
}

type OrderRequest struct {
	Subject  string `json:"subject"`
	Model    string `json:"model"`
	Language string `json:"language"`
}

type OrderResponse struct {
	Uuid string `json:"uuid"`
}

type OrderStatusRequest struct {
	Uuid string `json:"uuid"`
}

type OrderStatusResponse struct {
	Status       string `json:"status"`
	WikibookUuid string `json:"wikibook_uuid"`
}

type CompleteRequest struct {
	Value    string `json:"value"`
	Language string `json:"language"`
}

type CompleteResponse struct {
	Titles []string `json:"titles"`
}

type GetWikibookRequest struct {
	Uuid string `json:"uuid"`
}

type GetWikibookResponse struct {
	Wikibook *Wikibook `json:"wikibook"`
}

type ListWikibookRequest struct {
	Page     int64  `json:"page"`
	Size     int64  `json:"size"`
	Language string `json:"language"`
}

type ListWikibookResponse struct {
	Wikibooks []*Wikibook `json:"wikibooks"`
}

type DownloadWikibookRequest struct {
	Uuid   string `json:"uuid"`
	Format string `json:"format"`
}
