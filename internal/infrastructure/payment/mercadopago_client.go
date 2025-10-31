package payment

import (
	"context"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type MercadoPagoClient struct {
	client preference.Client
}

func NewMercadoPagoClient(accessToken string) *MercadoPagoClient {
	cfg, _ := config.New(accessToken)
	
	return &MercadoPagoClient{
		client: preference.NewClient(cfg),
	}
}

type PreferenceRequest struct {
	Title        string
	Description  string
	Quantity     int
	UnitPrice    float64
	ExternalRef  string
	NotifyURL    string
	BackURL      string
}

type PreferenceResponse struct {
	ID           string
	InitPointURL string
}

func (c *MercadoPagoClient) CreatePreference(ctx context.Context, req *PreferenceRequest) (*PreferenceResponse, error) {
	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:       req.Title,
				Description: req.Description,
				Quantity:    req.Quantity,
				UnitPrice:   req.UnitPrice,
			},
		},
		ExternalReference: req.ExternalRef,
		NotificationURL:   req.NotifyURL,
		BackURLs: &preference.BackURLsRequest{
			Success: req.BackURL + "/success",
			Failure: req.BackURL + "/failure",
			Pending: req.BackURL + "/pending",
		},
		AutoReturn: "approved",
	}

	result, err := c.client.Create(ctx, request)
	if err != nil {
		return nil, err
	}

	return &PreferenceResponse{
		ID:           result.ID,
		InitPointURL: result.InitPoint,
	}, nil
}

