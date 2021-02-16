package services

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"testing"
)

func Test_go_diff(t *testing.T) {
	const text = "whoami"
	const text2 = "whoami 333"

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text, text2, false)

	t.Log(dmp.DiffPrettyText(diffs))
}
