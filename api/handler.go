package wikibookgen

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	. "github.com/proullon/wikibookgen/api/model"
)

func WG(ctx context.Context) *WikibookGen {
	return ctx.Value("wg").(*WikibookGen)
}

func Writer(ctx context.Context) http.ResponseWriter {
	return ctx.Value("w").(http.ResponseWriter)
}

func statusHandler(ctx context.Context, v *Void) (*StatusResponse, error) {
	return &StatusResponse{}, nil
}

func completeHandler(ctx context.Context, req *CompleteRequest) (*CompleteResponse, error) {
	if len(req.Value) < 4 {
		return &CompleteResponse{}, nil
	}

	titles, err := WG(ctx).Complete(req.Value, req.Language)
	if err != nil {
		return nil, err
	}

	return &CompleteResponse{
		Titles: titles,
	}, nil
}

func orderHandler(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	uuid, err := WG(ctx).QueueGenerationJob(req.Subject, req.Model, req.Language)
	if err != nil {
		return nil, err
	}

	return &OrderResponse{
		Uuid: uuid,
	}, nil
}

func orderStatusHandler(ctx context.Context, req *OrderStatusRequest) (*OrderStatusResponse, error) {
	status, uuid, err := WG(ctx).LoadOrder(req.Uuid)
	if err != nil {
		return nil, err
	}

	return &OrderStatusResponse{
		Status:       status,
		WikibookUuid: uuid,
	}, nil
}

func listWikibookHandler(ctx context.Context, req *ListWikibookRequest) (*ListWikibookResponse, error) {
	list, err := WG(ctx).List(req.Page, req.Size, req.Language)
	if err != nil {
		return nil, err
	}

	return &ListWikibookResponse{
		Wikibooks: list,
	}, err
}

func getWikibookHandler(ctx context.Context, req *GetWikibookRequest) (*GetWikibookResponse, error) {
	book, err := WG(ctx).Load(req.Uuid)
	if err != nil {
		return nil, err
	}

	return &GetWikibookResponse{
		Wikibook: book,
	}, nil
}

func downloadWikibookHandler(ctx context.Context, req *DownloadWikibookRequest) (*Void, error) {
	return nil, WG(ctx).Download(req.Uuid, req.Format, Writer(ctx))
}

func getAvailableDownloadFormatHandler(ctx context.Context, req *GetAvailableDownloadFormatRequest) (*GetAvailableDownloadFormatResponse, error) {

	epub, pdf, err := WG(ctx).AvailableFormat(req.Uuid)
	if err != nil {
		return nil, err
	}

	return &GetAvailableDownloadFormatResponse{
		Epub: epub,
		Pdf:  pdf,
	}, nil
}

func printWikibookHandler(ctx context.Context, req *PrintWikibookRequest) (*Void, error) {
	err := WG(ctx).PrintWikibook(req.Uuid, req.Format)
	if err != nil {
		log.Errorf("Cannot print %+v: %s", req, err)
		return nil, err
	}

	return nil, nil
}
