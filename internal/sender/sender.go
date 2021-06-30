package sender

import "github.com/alipniczkij/binder/pkg/models"

type Interface interface {
	Send(event *models.GitlabEvent) error
}
