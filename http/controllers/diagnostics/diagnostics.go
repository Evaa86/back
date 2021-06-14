package diagnostics

import (
	"fmt"
	"net/http"

	"github.com/Zetkolink/back/http/helpers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Controller struct {
	models *ModelSet
}

type ModelSet struct {
	Bot *tgbotapi.BotAPI
}

type diagnosticsRequests struct {
	OtherDevice string `json:"other_device"`
	OtherIssue  string `json:"other_issue"`
	Issue       int64  `json:"issue"`
	Phone       string `json:"phone"`
}

type diagnosticsResponse struct {
}

func NewController(models ModelSet) *Controller {
	return &Controller{
		models: &models,
	}
}

func (c *Controller) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", c.Create)

	return r
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	payload := &diagnosticsRequests{}
	err := render.Bind(r, payload)

	if err != nil {
		fmt.Println(err)
		helpers.BadRequest(w, r, err)
		return
	}

	msg := tgbotapi.NewMessage(245713737,
		fmt.Sprintf("Новая заявка!\n"+
			"\nНомер телефона: %s"+
			"\nУстройство: %s"+
			"\nПроблема: %s\n",
			payload.Phone, payload.OtherDevice, payload.OtherIssue))
	_, _ = c.models.Bot.Send(msg)

	w.WriteHeader(http.StatusCreated)
}

func (drq *diagnosticsRequests) Bind(_ *http.Request) error {
	return nil
}
