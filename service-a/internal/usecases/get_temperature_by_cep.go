package usecases

import (
	"context"

	"github.com/arfurlaneto/goexpert-challenge-opentelemetry/service-a/internal/services"
)

type TemperatureByCepInput struct {
	Cep string
}

type TemperatureByCepOutput struct {
	StatusCode int
	Data       string
}

type GetTemperatureByCepUseCase interface {
	Execute(ctx context.Context, input *TemperatureByCepInput) (*TemperatureByCepOutput, error)
}

type GetTemperatureByCepUseCaseImpl struct {
	serviceBService services.ServiceBService
}

func NewGetTemperatureByCepUseCase(serviceBService services.ServiceBService) GetTemperatureByCepUseCase {
	return &GetTemperatureByCepUseCaseImpl{
		serviceBService: serviceBService,
	}
}

func (u *GetTemperatureByCepUseCaseImpl) Execute(ctx context.Context, input *TemperatureByCepInput) (*TemperatureByCepOutput, error) {
	serviceBResponse, err := u.serviceBService.QueryCep(ctx, input.Cep)
	if err != nil {
		return nil, err
	}

	return &TemperatureByCepOutput{
		StatusCode: serviceBResponse.StatusCode,
		Data:       serviceBResponse.Data,
	}, nil
}
