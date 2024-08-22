package save

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	"service_on_go/internal/lib/api/response"
	"service_on_go/internal/lib/logger/sl"
	"service_on_go/internal/lib/random"
	"service_on_go/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"` // error || ok
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

// TODO then move to config
const aliasLength = 6

type URLSaver interface {
	SaveUrl(url string, alias string) (uint, error)
}

func New(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.save.New"

		requestID, exists := c.Get("RequestID")
		if !exists {
			c.IndentedJSON(500, response.Error("RequestID not found"))
			return
		}

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", fmt.Sprintf("%s", requestID.(string))),
		)

		var req Request
		err := c.ShouldBindJSON(&req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			c.IndentedJSON(400, response.Error("failed to decode request")) // todo if it not work try to with gin.H{}

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			c.IndentedJSON(400, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exist", slog.String("url", req.URL))

			c.IndentedJSON(http.StatusConflict, response.Error("url already exist"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			c.IndentedJSON(http.StatusBadRequest, response.Error("failed to add url"))
		}

		log.Info("url added", slog.Uint64("id", uint64(id)))

		c.IndentedJSON(http.StatusOK, gin.H{
			"Status": http.StatusOK,
			"Alias":  alias,
		})
	}
}

//curl http://localhost:8082/url \
//--include \
//--header "Content-Type: application/json" \
//--request "POST" \
//--data '{"url": "htts://google.com","alias": "my_google"}'
