package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/sender"
	"github.com/alipniczkij/binder/pkg/models"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	sender sender.Interface
	cmds   map[string]commands.Handler
}

func New(sender sender.Interface, cmds map[string]commands.Handler) *Handler {
	return &Handler{sender: sender, cmds: cmds}
}

func (h *Handler) RegisterPublicHTTP() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/gitlab", h.getGitlabEvent()).Methods(http.MethodPost)
	r.HandleFunc("/slack", h.getSlackCommand()).Methods(http.MethodPost)
	return r
}

func (h *Handler) getGitlabEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &models.GitlabEvent{}
		err := unmarshalRequestBody(r, req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		err = h.sender.Send(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
	}
}

func (h *Handler) getSlackCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cmdHandler, found := h.cmds[s.Command]
		if !found {
			sendMsg(commandNotFound(), w)
		}
		sendMsg(cmdHandler.Execute(s), w)
	}
}

func unmarshalRequestBody(r *http.Request, output interface{}) error {
	var err error
	defer r.Body.Close()
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	_, err = io.Copy(buffer, r.Body)
	if err != nil {
		return err
	}
	b := buffer.Bytes()
	marshaller, ok := output.(json.Unmarshaler)
	if ok {
		err = marshaller.UnmarshalJSON(b)
	} else {
		err = json.NewDecoder(bytes.NewBuffer(b)).Decode(output)
	}
	return err
}

func sendMsg(msg slack.Msg, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	responseData, _ := json.Marshal(msg)
	w.Write(responseData)
}

func commandNotFound() slack.Msg {
	return slack.Msg{
		ResponseType: slack.ResponseTypeInChannel,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.SectionBlock{
					Type: "section",
					Text: slack.NewTextBlockObject(
						slack.MarkdownType,
						"Command not found",
						false,
						false,
					),
				},
			},
		},
	}
}
