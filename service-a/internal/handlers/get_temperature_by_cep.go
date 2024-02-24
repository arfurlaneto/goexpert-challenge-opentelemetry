package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/arfurlaneto/goexpert-challenge-opentelemetry/service-a/internal/usecases"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type GetTemperatureByCepRequest struct {
	Cep string `json:"cep"`
}

type GetTemperatureByCepHandler struct {
	getTemperatureByCepUseCase usecases.GetTemperatureByCepUseCase
}

func NewGetTemperatureByCepHandler(getTemperatureByCepUseCase usecases.GetTemperatureByCepUseCase) *GetTemperatureByCepHandler {
	return &GetTemperatureByCepHandler{
		getTemperatureByCepUseCase: getTemperatureByCepUseCase,
	}
}

func (h *GetTemperatureByCepHandler) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, "GetTemperatureByCepHandler")
	defer span.End()

	cep, ok := h.getCepFromBody(r)
	if !ok {
		w.WriteHeader(422)
		w.Write([]byte("invalid zipcode"))
		return
	}

	input := &usecases.TemperatureByCepInput{Cep: cep}
	output, err := h.getTemperatureByCepUseCase.Execute(ctx, input)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		w.WriteHeader(500)
		w.Write([]byte("internal server error"))
		return
	}

	w.WriteHeader(output.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(output.Data))
}

func (h *GetTemperatureByCepHandler) getCepFromBody(r *http.Request) (string, bool) {
	var requestData GetTemperatureByCepRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return "", false
	}

	cep := requestData.Cep

	if cep == "" {
		return "", false
	}

	cepRegex := regexp.MustCompile(`^\d{5}-{0,1}\d{3}$`)
	if !cepRegex.Match([]byte(cep)) {
		return "", false
	}

	return cep, true
}
