package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ServiceBResponse struct {
	StatusCode int
	Data       string
}

type ServiceBService interface {
	QueryCep(ctx context.Context, cep string) (*ServiceBResponse, error)
}

type ServiceBServiceImpl struct {
	client *http.Client
}

func NewServiceBService() ServiceBService {
	return &ServiceBServiceImpl{
		client: &http.Client{},
	}
}

func (s *ServiceBServiceImpl) QueryCep(ctx context.Context, cep string) (*ServiceBResponse, error) {
	url := fmt.Sprintf("%s?cep=%s", os.Getenv("SERVICE_B_URL"), cep)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))
	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Service B error")
	}

	return &ServiceBResponse{
		StatusCode: response.StatusCode,
		Data:       string(body),
	}, nil
}
