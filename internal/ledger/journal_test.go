package ledger

import (
	"context"
	"testing"
)

func TestRejectUnbalancedOrDuplicateAccountsBeforeDatabase(t *testing.T) {
	for _, entries := range [][]Entry{{{Account: "cash", Amount: 1}}, {{Account: "cash", Amount: 2}, {Account: "liability", Amount: -1}}, {{Account: "cash", Amount: 1}, {Account: "cash", Amount: -1}}} {
		if Post(context.Background(), nil, "j", "invalid", entries...) == nil {
			t.Fatal("invalid journal accepted")
		}
	}
}
