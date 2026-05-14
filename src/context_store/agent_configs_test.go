package context_store

import "testing"

func TestAgentConfigLifecycle(t *testing.T) {
	store := openTestStore(t, t.TempDir())
	defer store.Close()

	config, err := store.GetAgentConfig("CEO")
	if err != nil {
		t.Fatalf("GetAgentConfig seeded CEO: %v", err)
	}
	if config.Model == "" || config.Tier != 3 {
		t.Fatalf("seeded config = %#v", config)
	}

	updated, err := store.UpsertAgentConfig(AgentConfig{
		AgentName: "QA",
		Model:     "qwen2.5-coder",
		Tier:      2,
		Active:    true,
		UpdatedBy: "test",
	})
	if err != nil {
		t.Fatalf("UpsertAgentConfig: %v", err)
	}
	if updated.Model != "qwen2.5-coder" || updated.Tier != 2 {
		t.Fatalf("updated = %#v", updated)
	}

	configs, err := store.ListAgentConfigs()
	if err != nil {
		t.Fatalf("ListAgentConfigs: %v", err)
	}
	if len(configs) == 0 {
		t.Fatal("expected seeded configs")
	}
}
