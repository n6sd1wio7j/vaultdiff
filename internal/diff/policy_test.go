package diff

import (
	"strings"
	"testing"
)

var samplePolicyEntries = []Entry{
	{Key: "db_password", OldValue: "", NewValue: "s3cr3t", Status: StatusAdded},
	{Key: "api_key", OldValue: "old", NewValue: "", Status: StatusRemoved},
	{Key: "auth_token", OldValue: "v1", NewValue: "v2", Status: StatusModified},
	{Key: "app_name", OldValue: "myapp", NewValue: "myapp", Status: StatusUnchanged},
}

func TestEnforcePolicy_DetectsAddedPassword(t *testing.T) {
	rules := DefaultPolicyRules()
	vs := EnforcePolicy(samplePolicyEntries, rules)
	for _, v := range vs {
		if v.Entry.Key == "db_password" && v.Rule.Name == "warn-added-password" {
			return
		}
	}
	t.Error("expected warn-added-password violation for db_password")
}

func TestEnforcePolicy_DetectsRemovedKey(t *testing.T) {
	rules := DefaultPolicyRules()
	vs := EnforcePolicy(samplePolicyEntries, rules)
	for _, v := range vs {
		if v.Entry.Key == "api_key" && v.Rule.Name == "no-removed-keys" {
			return
		}
	}
	t.Error("expected no-removed-keys violation for api_key")
}

func TestEnforcePolicy_DetectsModifiedToken(t *testing.T) {
	rules := DefaultPolicyRules()
	vs := EnforcePolicy(samplePolicyEntries, rules)
	for _, v := range vs {
		if v.Entry.Key == "auth_token" && v.Rule.Name == "warn-modified-token" {
			return
		}
	}
	t.Error("expected warn-modified-token violation for auth_token")
}

func TestEnforcePolicy_UnchangedNoViolation(t *testing.T) {
	rules := DefaultPolicyRules()
	vs := EnforcePolicy([]Entry{
		{Key: "app_name", OldValue: "x", NewValue: "x", Status: StatusUnchanged},
	}, rules)
	if len(vs) != 0 {
		t.Errorf("expected no violations, got %d", len(vs))
	}
}

func TestFormatViolations_NoViolations(t *testing.T) {
	out := FormatViolations(nil)
	if !strings.Contains(out, "passed") {
		t.Error("expected 'passed' in output")
	}
}

func TestFormatViolations_ListsMessages(t *testing.T) {
	vs := EnforcePolicy(samplePolicyEntries, DefaultPolicyRules())
	out := FormatViolations(vs)
	if !strings.Contains(out, "violations") {
		t.Error("expected 'violations' header in output")
	}
	if !strings.Contains(out, "api_key") {
		t.Error("expected api_key in output")
	}
}
