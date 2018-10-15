package fluidpay

import (
	"testing"
)

func TestUrlBuilder(t *testing.T) {
	params := []string{"mycustom", "path"}
	actRes := urlBuilder(params)
	if actRes != "/mycustom/path" {
		t.Errorf("actRes should be '/mycustom/path', instead of %v", actRes)
	}
}
