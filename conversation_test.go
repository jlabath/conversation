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
		t.Fatal("expected to have exatly 1 intent match")
	}
	if response.Intents[0].Intent != "complaint" {
		t.Fatalf(
			"expected complaint but got %s",
			response.Intents[0].Intent)
	}
}
