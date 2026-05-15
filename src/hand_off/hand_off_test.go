package hand_off

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCreateHandoffGeneratesCardID(t *testing.T) {
	path, err := CreateHandoff(t.TempDir(), HandoffHeader{
		Sender:    "[CEO]",
		Recipient: "[DEV_BACKEND]",
		Intent:    "IMPLEMENT",
	}, map[string]string{"task": "build"})
	if err != nil {
		t.Fatalf("CreateHandoff: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var handoff HandoffSchema[map[string]string]
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatal(err)
	}
	if handoff.Header.CardID == "" {
		t.Fatal("card_id vazio")
	}
}

func TestSaveRawHandoffPreservesCardID(t *testing.T) {
	raw := []byte(`{"header":{"card_id":"card-1","sender":"[CEO]","recipient":"[DEV_BACKEND]","intent":"IMPLEMENT"},"payload":{"task":"build"}}`)
	path, err := SaveRawHandoff(t.TempDir(), raw)
	if err != nil {
		t.Fatalf("SaveRawHandoff: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var handoff HandoffSchema[map[string]string]
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatal(err)
	}
	if handoff.Header.CardID != "card-1" {
		t.Fatalf("card_id = %q, want card-1", handoff.Header.CardID)
	}
}

func TestSaveRawHandoffAcceptsUTF8BOM(t *testing.T) {
	raw := append([]byte{0xEF, 0xBB, 0xBF}, []byte(`{"header":{"card_id":"card-bom","sender":"[CEO]","recipient":"[DEV_BACKEND]","intent":"IMPLEMENT"},"payload":{"task":"build"}}`)...)
	path, err := SaveRawHandoff(t.TempDir(), raw)
	if err != nil {
		t.Fatalf("SaveRawHandoff: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var handoff HandoffSchema[map[string]string]
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatal(err)
	}
	if handoff.Header.CardID != "card-bom" {
		t.Fatalf("card_id = %q, want card-bom", handoff.Header.CardID)
	}
}
