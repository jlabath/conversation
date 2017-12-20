### Conversation

[![GoDoc](https://godoc.org/github.com/jlabath/conversation?status.svg)](https://godoc.org/github.com/jlabath/conversation)

Conversation is a minimal implementation of IBM Watson conversation [API](https://www.ibm.com/watson/developercloud/conversation/api/v1/).

Minimal code may look like this

```golang

creds := conversation.Credentials{
     User: "Sherlock",
     Password: "Holmes",
     WorkspaceID: "1234",
}

c := &http.Client{}
conv := conversation.New(creds, c)
req := conversation.Continue("Hi Watson!")
response, err := conv.Message(req)

```
