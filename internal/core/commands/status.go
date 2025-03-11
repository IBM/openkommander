package commands

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

type Failure struct {
	Err      error
	HttpCode int
}

func NewFailure(err string, httpCode int) *Failure {
	return &Failure{
		Err:      errors.New(err),
		HttpCode: httpCode,
	}
}

func GetAdminClient() (sarama.ClusterAdmin, *Failure) {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		return nil, NewFailure("No active session found", http.StatusUnauthorized)
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		return nil, NewFailure(fmt.Sprintf("Error connecting to cluster: %v", err), http.StatusInternalServerError)
	}

	return client, nil
}

func GetClient() (sarama.Client, *Failure) {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		return nil, NewFailure("No active session found", http.StatusUnauthorized)
	}

	client, err := currentSession.GetClient()
	if err != nil {
		return nil, NewFailure(fmt.Sprintf("Error connecting to cluster: %v", err), http.StatusInternalServerError)
	}

	return client, nil
}
