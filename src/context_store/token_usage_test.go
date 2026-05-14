package context_store

import "testing"

func TestTokenUsageLifecycle(t *testing.T) {
	store := openTestStore(t, t.TempDir())
	defer store.Close()
	if _, err := store.SaveCardHandoff(CardHandoff{
		CardID:    "card-token",
		Title:     "Token card",
		Sender:    "[CEO]",
		Recipient: "[DEV_BACKEND]",
		Intent:    "IMPLEMENT",
		RawJSON:   `{}`,
	}); err != nil {
		t.Fatal(err)
	}

	if err := store.SaveTokenUsage(TokenUsage{
		CardID:           "card-token",
		AgentName:        "DEV_BACKEND",
		Model:            "qwen2.5-coder",
		Tier:             2,
		PromptTokens:     10,
		CompletionTokens: 5,
		LatencyMS:        123,
	}); err != nil {
		t.Fatalf("SaveTokenUsage: %v", err)
	}

	sum, err := store.SumTokenUsage(TokenUsageFilter{AgentName: "DEV_BACKEND"})
	if err != nil {
		t.Fatalf("SumTokenUsage: %v", err)
	}
	if sum.TotalTokens != 15 {
		t.Fatalf("sum = %#v", sum)
	}

	usages, err := store.ListTokenUsageByCard("card-token")
	if err != nil {
		t.Fatalf("ListTokenUsageByCard: %v", err)
	}
	if len(usages) != 1 || usages[0].TotalTokens != 15 {
		t.Fatalf("usages = %#v", usages)
	}

	thread, err := store.GetCardThread("card-token")
	if err != nil {
		t.Fatal(err)
	}
	if len(thread.Comments) != 2 || thread.Comments[1].CommentType != "token_usage" {
		t.Fatalf("thread = %#v", thread)
	}
}
