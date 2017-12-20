package conversation

import (
	"net/http"
	"os"
	"testing"
)

var cred Credentials

func init() {
	cred.User = os.Getenv("USER")
	cred.Password = os.Getenv("PASS")
	cred.WorkspaceID = os.Getenv("WORKSPACE")
}

func TestConv(t *testing.T) {
	c := &http.Client{}
	conv := New(cred, c)
	req := MessageRequest{
		Input: InputData{"I have a huge problem"},
	}
	response, err := conv.Message(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(response.Intents) != 1 {
		t.Logf("%+v", response.Intents)
		t.Fatal("expected to have exatly 1 intent match")
	}
	if response.Intents[0].Intent != "complaint" {
		t.Fatalf(
			"expected complaint but got %s",
			response.Intents[0].Intent)
	}
	if response.Intents[0].Confidence < 0.5 {
		t.Fatalf(
			"expected hight confidende but got %.2f",
			response.Intents[0].Confidence)
	}
	if len(response.Entities) > 0 {
		t.Fatalf(
			"unexpected entities in response %+v",
			response.Entities)
	}
	if response.Context.ConversationID == "" {
		t.Fatalf("expected to have conv ID in context")
	}
	oldCtxID := response.Context.ConversationID

	if len(response.Output.Text) < 1 {
		t.Fatalf("expected to see some output texts in response")
	}

	t.Log(response.Input.Text)

	for _, txt := range response.Output.Text {
		t.Log(txt)
	}
	if len(response.Output.NodesVisited) < 1 {
		t.Fatalf("expected to see some visited nodes in response")
	}

	//continue conversation
	response, err = conv.Message(conv.Continue("I'd like to return this book"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.Context.ConversationID != oldCtxID {
		t.Fatalf(
			"Expected ConversationID: %s but got: %s",
			oldCtxID,
			response.Context.ConversationID)
	}
	if len(response.Entities) != 1 {
		t.Fatalf("expected entities with one result but got %+v", response.Entities)
	}

	ent := response.Entities[0]
	if ent.Entity != "returnItems" {
		t.Fatalf("expected returnItems but got %s", ent.Entity)
	}
	if ent.Value != "book" {
		t.Fatalf("expected book but got %s", ent.Value)
	}

}
