package conversation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//APIVersion of Watson conversation service.
const APIVersion = "2017-05-26"

//Conversation is the connection to Watson conversation service.
type Conversation struct {
	b            Broker
	cred         Credentials
	lastResponse *MessageResponse
}

//New returns new conversation connection using the provided credentials and
//the broker.
func New(c Credentials, b Broker) *Conversation {
	return &Conversation{
		b:    b,
		cred: c,
	}
}

//Message sends a HTTP POST message (via broker)
//it returns non nil error in case there was any trouble
//https://www.ibm.com/watson/developercloud/conversation/api/v1/#messages
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
	//save lastResponse
	c.lastResponse = &resp

	return resp, nil
}

func (c *Conversation) url() (string, error) {
	return fmt.Sprintf(
		"https://gateway.watsonplatform.net/conversation/api/v1/workspaces/%s/message?version=%s",
		c.cred.WorkspaceID, APIVersion), nil
}

//Continue will continue the conversation by preserving
//the Context from the last response if any.
func (c *Conversation) Continue(text string) MessageRequest {
	var req MessageRequest
	req.Input.Text = text
	if c.lastResponse != nil {
		req.Context = &c.lastResponse.Context
	}
	return req
}

//Credentials contain relevant workspace and login info.
type Credentials struct {
	User        string
	Password    string
	WorkspaceID string
}

//Broker is the interface that will be responsible for delivering
//the HTTP request to Watson's APIs
//http.Client from standard library implements this interface
//
//    var b Broker = &http.Client{}
//
type Broker interface {
	Do(*http.Request) (*http.Response, error)
}

//MessageRequest to send to Watson.
//Best to review official IBM documentation at https://www.ibm.com/watson/developercloud/conversation/api/v1/#send_message.
//Minimal example
//
//	mr := MessageRequest{
//		Input: InputData{"How are you today Watson?"},
//	}
//
type MessageRequest struct {
	Input            InputData       `json:"input"`
	AlternateIntents bool            `json:"alternate_intents,omitempty"`
	Context          *Context        `json:"context,omitempty"`
	Intents          []RuntimeIntent `json:"intents,omitempty"`
	Entities         []RuntimeEntity `json:"entities,omitempty"`
	Output           *OutputData     `json:"output,omitempty"`
}

//InputData is the actual text we want Watson to read.
type InputData struct {
	Text string `json:"text"`
}

//MessageResponse as receieved from Watson.
type MessageResponse struct {
	Input            MessageInput    `json:"input"`
	Intents          []RuntimeIntent `json:"intents"`
	Entities         []RuntimeEntity `json:"entities"`
	AlternateIntents bool            `json:"alternate_intents"`
	Context          Context         `json:"context"`
	Output           OutputData      `json:"output"`
}

//MessageInput is input that triggered this response.
type MessageInput struct {
	Text string `json:"text"`
}

//RuntimeIntent is the recognized intent.
type RuntimeIntent struct {
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`
}

//RuntimeEntity is the recognized entity.
type RuntimeEntity struct {
	Entity     string                 `json:"entity"`
	Location   []int                  `json:"location"`
	Value      string                 `json:"value"`
	Confidence float64                `json:"confidence"`
	Metadata   map[string]interface{} `json:"metadata"`
}

//Context stores the context for a conversation.
type Context struct {
	ConversationID string      `json:"conversation_id"`
	System         interface{} `json:"system"`
}

//OutputData provided by Watson.
type OutputData struct {
	LogMessages  []LogMessage `json:"log_messages"`
	Text         []string     `json:"text"`
	NodesVisited []string     `json:"nodes_visited"`
}

//LogMessage a message logged with the request.
type LogMessage struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}
