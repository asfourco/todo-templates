package api

import (
	"net/http"

	"github.com/asfourco/templates/backend/models"
)

func getHealth(r *http.Request) (interface{}, error) {
	return models.HealthResponse{Healthy: true}, nil
}
