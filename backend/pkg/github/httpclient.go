package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ng-vu/githubx/backend/pkg/xerrors"
)

type Client struct {
	Client  *http.Client
	Token   string
	OrgRepo OrgRepo

	log *log.Logger
}

func NewClient(cfg OrgRepo, token string, logger io.Writer) *Client {
	return &Client{
		Client:  &http.Client{Timeout: 60 * time.Second},
		Token:   token,
		OrgRepo: cfg,
		log:     log.New(logger, "", log.Ldate|log.Ltime|log.Llongfile),
	}
}

type RequestBody struct {
	Query string `json:"query"`
}

type Response struct {
	Data   json.RawMessage `json:"data"`
	Errors []Error         `json:"errors"`
}

type ResponseX struct {
	Data   interface{} `json:"data"`
	Errors []Error     `json:"errors,omitempty"`
}

type Error struct {
	Path       []string        `json:"path"`
	Extensions *ErrorExtension `json:"extensions"`
	Location   []ErrorLocation `json:"location"`
	Message    string          `json:"message"`
}

type ErrorExtension struct {
	Code     string `json:"code"`
	NodeName string `json:"nodeName"`
	TypeName string `json:"typeName"`
}

type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (c *Client) WithToken(token string) *Client {
	client := *c
	client.Token = token
	return &client
}

func (c *Client) SendRequest(ctx context.Context, query string, v interface{}) (*http.Response, error) {
	query = strings.TrimSpace(query)
	reqBody := &RequestBody{Query: query}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, xerrors.Errorf(err, "%v", err)
	}

	var resp *http.Response
	var respBody []byte
	var wrapResp Response
	defer func() {
		c.log.Printf("POST /graphql\n%s\n", query)
		if err != nil {
			c.log.Printf("--> [error] %v\n", err)
			return
		}
		if resp.StatusCode != 200 {
			c.log.Printf("--> [error] code=%v %s\n", resp.StatusCode, respBody)
			return
		}

		resp0 := jsonReindent(respBody, "", "  ")
		resp1 := jsonMarshalIndent(ResponseX{
			Data:   v,
			Errors: wrapResp.Errors,
		})
		c.log.Printf("--> [raw]\n%s\n", resp0)
		c.log.Printf("--> [...]\n%s\n", resp1)
	}()

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, xerrors.Errorf(err, "%v", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf(err, "%v", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf(err, "%v", err)
	}
	if resp.StatusCode != 200 {
		return nil, xerrors.Errorf(nil, "http %v: %s", resp.Status, respBody)
	}
	if err = json.Unmarshal(respBody, &wrapResp); err != nil {
		return nil, xerrors.Errorf(err, "%v", err)
	}
	if wrapResp.Data == nil {
		return resp, xerrors.Errorf(err, "%v", wrapResp.Errors[0].Message)
	}
	if err = json.Unmarshal(wrapResp.Data, v); err != nil {
		return resp, xerrors.Errorf(err, "%v", err)
	}
	return resp, nil
}

func (c *Client) Ping(ctx context.Context) (*LoginResponse, error) {
	var resp LoginResponse
	_, err := c.SendRequest(ctx, LoginQuery(c.OrgRepo), &resp)
	return &resp, err
}

func (c *Client) Login() {

}
