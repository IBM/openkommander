package commands

import (
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

type Success[T any] struct {
	Body T
}

func (Success[T]) IsSuccess() bool {
	return true
}

type Failure struct {
	Err      error
	HttpCode int
}

func (Failure) IsSuccess() bool {
	return false
}

func NewSuccess[T any](body T) *Success[T] {
	return &Success[T]{
		Body: body,
	}
}

func NewFailure(err string, httpCode int) *Failure {
	return &Failure{
		Err:      fmt.Errorf(err),
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
