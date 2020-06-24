package sutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ng-vu/githubx/backend/pkg/xerrors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Requester struct{ *http.Request }

func Request(r *http.Request) Requester { return Requester{r} }

func (r Requester) Decode(v interface{}) error {
	if r.Method != "POST" {
		return xerrors.Errorf(nil, "method not supported")
	}
	body, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return err
	}
	log.Printf("--> %s\n%s", r.URL.Path, body)
	if err = json.Unmarshal(body, v); err != nil {
		return err
	}
	decoded, _ := json.MarshalIndent(v, "", "  ")
	log.Printf("--> %s [decoded]\n%s", r.URL.Path, decoded)
	return nil
}

type Responder struct{ http.ResponseWriter }

func Respond(w http.ResponseWriter) Responder { return Responder{w} }

func (r Responder) Error(code int, err error) {
	resp := &ErrorResponse{Error: fmt.Sprintf("%v", err)}
	respData, _ := json.Marshal(resp)
	r.ResponseWriter.WriteHeader(code)
	_, _ = r.ResponseWriter.Write(respData)
	log.Printf("--> %s", respData)
}

func (r Responder) Resp(v interface{}, err error) {
	if err != nil {
		r.Error(400, err)
		return
	}
	r.Header().Set("Content-Type", "application/json")
	respData, err := json.Marshal(v)
	must(err)
	_, _ = r.Write(respData)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
