package model

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
