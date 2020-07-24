package model

type Wikibook struct {
	Uuid           string    `json:"uuid,omitempty"`
	Subject        string    `json:"subject,omitempty"`
	Model          string    `json:"model,omitempty"`
	Language       string    `json:"language,omitempty"`
	Title          string    `json:"title,omitempty"`
	Pages          int64     `json:"pages,omitempty"`
	Volumes        []*Volume `json:"volumes,omitempty"`
	GenerationDate string    `json:"generation_date,omitempty"`
}

type Volume struct {
	Title    string     `json:"title,omitempty"`
	Chapters []*Chapter `json:"chapters,omitempty"`
}

type Chapter struct {
	Title    string  `json:"title,omitempty"`
	Articles []*Page `json:"articles,omitempty"`
}

type Page struct {
	Id      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type StatusResponse struct {
	Status []string `json:"status,omitempty"`
}

type Void struct {
}

type OrderRequest struct {
	Subject  string `json:"subject,omitempty"`
	Model    string `json:"model,omitempty"`
	Language string `json:"language,omitempty"`
}

type OrderResponse struct {
	Uuid string `json:"uuid,omitempty"`
}

type OrderStatusRequest struct {
	Uuid string `json:"uuid,omitempty"`
}

type OrderStatusResponse struct {
	Status       string `json:"status,omitempty"`
	WikibookUuid string `json:"wikibook_uuid,omitempty"`
}

type CompleteRequest struct {
	Value    string `json:"value,omitempty"`
	Language string `json:"language,omitempty"`
}

type CompleteResponse struct {
	Titles []string `json:"titles,omitempty"`
}

type GetWikibookRequest struct {
	Uuid string `json:"uuid,omitempty"`
}

type GetWikibookResponse struct {
	Wikibook *Wikibook `json:"wikibook,omitempty"`
}

type ListWikibookRequest struct {
	Page     int64  `json:"page,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Language string `json:"language,omitempty"`
}

type ListWikibookResponse struct {
	Wikibooks []*Wikibook `json:"wikibooks,omitempty"`
}

type DownloadWikibookRequest struct {
	Uuid   string `json:"uuid,omitempty"`
	Format string `json:"format,omitempty"`
}

type GetAvailableDownloadFormatRequest struct {
	Uuid string `json:"uuid,omitempty"`
}

type GetAvailableDownloadFormatResponse struct {
	Epub bool `json:"epub,omitempty"`
	Pdf  bool `json:"pdf,omitempty"`
}
