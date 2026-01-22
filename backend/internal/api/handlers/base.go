package api

import (
	config "github.com/chungweeeei/Temporal-robot-project/internal/config/api"
	"github.com/chungweeeei/Temporal-robot-project/internal/repository/dao"
)

type Handler struct {
	App         *config.AppConfig
	WorkflowDAO *dao.WorkflowDAO
}

func NewHandler(app *config.AppConfig, wfDAO *dao.WorkflowDAO) *Handler {
	return &Handler{
		App:         app,
		WorkflowDAO: wfDAO,
	}
}
