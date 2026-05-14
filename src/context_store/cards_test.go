package context_store

import "testing"

func TestCardLifecycle(t *testing.T) {
	store := openTestStore(t, t.TempDir())
	defer store.Close()

	card, err := store.SaveCardHandoff(CardHandoff{
		CardID:    "card-1",
		Title:     "Implement endpoint",
		TaskRef:   "TASK-1",
		Sender:    "[CEO]",
		Recipient: "[DEV_BACKEND]",
		Intent:    "IMPLEMENT",
		RawJSON:   `{"ok":true}`,
	})
	if err != nil {
		t.Fatalf("SaveCardHandoff: %v", err)
	}
	if card.Status != "open" {
		t.Fatalf("status = %q, want open", card.Status)
	}

	if err := store.UpdateCardStatus("card-1", "in_progress"); err != nil {
		t.Fatalf("UpdateCardStatus: %v", err)
	}
	if err := store.AddCardComment("card-1", "DEV_BACKEND", "response", "done"); err != nil {
		t.Fatalf("AddCardComment: %v", err)
	}
	if err := store.CancelCard("card-1"); err != nil {
		t.Fatalf("CancelCard: %v", err)
	}

	thread, err := store.GetCardThread("card-1")
	if err != nil {
		t.Fatalf("GetCardThread: %v", err)
	}
	if thread.Card.Status != "canceled" {
		t.Fatalf("status = %q, want canceled", thread.Card.Status)
	}
	if len(thread.Comments) != 2 {
		t.Fatalf("comments = %d, want 2", len(thread.Comments))
	}

	cards, err := store.ListCards(CardFilter{Status: "canceled", TaskRef: "TASK-1"})
	if err != nil {
		t.Fatalf("ListCards: %v", err)
	}
	if len(cards) != 1 || cards[0].ID != "card-1" {
		t.Fatalf("cards = %#v", cards)
	}
}

func TestHumanCardStartsBlocked(t *testing.T) {
	store := openTestStore(t, t.TempDir())
	defer store.Close()

	card, err := store.SaveCardHandoff(CardHandoff{
		CardID:    "card-human",
		Sender:    "[DEV_BACKEND]",
		Recipient: "[HUMAN]",
		Intent:    "HUMAN_CLARIFICATION_REQUEST",
		RawJSON:   `{}`,
	})
	if err != nil {
		t.Fatalf("SaveCardHandoff: %v", err)
	}
	if card.Status != "blocked" {
		t.Fatalf("status = %q, want blocked", card.Status)
	}
}

func TestRecordCardFailureRetriesUntilLimit(t *testing.T) {
	store := openTestStore(t, t.TempDir())
	defer store.Close()
	if _, err := store.SaveCardHandoff(CardHandoff{
		CardID:    "card-retry",
		Title:     "Retry me",
		Sender:    "[CEO]",
		Recipient: "[DEV_BACKEND]",
		Intent:    "IMPLEMENT",
		RawJSON:   `{}`,
	}); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		retry, err := store.RecordCardFailure("card-retry", "temporary error")
		if err != nil {
			t.Fatalf("RecordCardFailure %d: %v", i, err)
		}
		if !retry {
			t.Fatalf("retry %d = false, want true", i)
		}
	}
	retry, err := store.RecordCardFailure("card-retry", "final error")
	if err != nil {
		t.Fatalf("RecordCardFailure final: %v", err)
	}
	if retry {
		t.Fatal("retry after limit = true, want false")
	}
	card, err := store.GetCard("card-retry")
	if err != nil {
		t.Fatal(err)
	}
	if card.Status != "failed" || card.RetryCount != 3 {
		t.Fatalf("card = %#v", card)
	}
}
