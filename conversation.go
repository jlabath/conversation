package conversation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const APIVersion = "2017-05-26"

type Conversation struct {
	b    Broker
	cred Credentials
}

func New(c Credentials, b Broker) *Conversation {
	return &Conversation{
		b:    b,
		cred: c,
	}
}

func (c *Conversation) Message(r MessageRequest) (MessageResponse, error) {
	var (
		resp MessageResponse
		body bytes.Buffer
	)
	err := json.NewEncoder(&body).Encode(r)
	if err != nil {
		return resp, err
	}
	//craft ourselves new http request
	url, err := c.url()
	if err != nil {
		return resp, err
	}
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.cred.User, c.cred.Password)

	//call broker to actually perform the request
	httpResp, err := c.b.Do(req)
	if err != nil {
		return resp, err
	}
	//ok process the response
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("HTTP %d", httpResp.StatusCode)
	}
	//ok all good try to decode
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Conversation) url() (string, error) {
	return fmt.Sprintf(
		"https://gateway.watsonplatform.net/conversation/api/v1/workspaces/%s/message?version=%s",
		c.cred.WorkspaceID, APIVersion), nil
}

type Credentials struct {
	User        string
	Password    string
	WorkspaceID string
}

//request

type MessageRequest struct {
	Input InputData `json:"input"`
}

type InputData struct {
	Text string `json:"text"`
}

//response

type MessageResponse struct {
	Input   MessageInput    `json:"input"`
	Intents []RuntimeIntent `json:"intents"`
}

type MessageInput struct {
	Text string `json:"text"`
}

type RuntimeIntent struct {
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`
}

type Broker interface {
	Do(*http.Request) (*http.Response, error)
}
