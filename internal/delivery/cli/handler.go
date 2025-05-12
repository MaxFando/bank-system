package cli

import (
	"context"
	"github.com/go-co-op/gocron"
	"time"
)

type Handler struct {
	Scheduler *gocron.Scheduler
}

func New() *Handler {
	return &Handler{
		Scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (h *Handler) Start() {
	_, err := h.Scheduler.Every(1).Hour().Do(h.CheckCredits)
	if err != nil {
		panic(err)
	}

	h.Scheduler.StartAsync()
}

func (h *Handler) CheckCredits(ctx context.Context) {

}
