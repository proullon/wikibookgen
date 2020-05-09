package wikibookgen

import (
	"context"

	. "github.com/proullon/wikibookgen/api/model"
)

func WG(ctx context.Context) *WikibookGen {
	return ctx.Value("wg").(*WikibookGen)
}

func statusHandler(ctx context.Context, v *Void) (*StatusResponse, error) {
	return &StatusResponse{}, nil
}

func orderHandler(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	uuid, err := WG(ctx).QueueGenerationJob(req.Subject, req.Model)
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
	list, err := WG(ctx).List()
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
	return nil, nil
}