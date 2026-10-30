package ops

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestPackCommandRunsOnlyFleetWideAssets(t *testing.T) {
	asset := func(lines ...string) string {
		return "pack: p\nassets:\n  a:\n    " + strings.Join(lines, "\n    ") + "\n"
	}
	for _, tc := range []struct {
		name string
		pack string
		ok   bool
	}{
		{"run_id as a whole argument", asset(`run: [tool, "{run_id}"]`), true},
		{"no placeholder", asset("run: [tool, --all]"), true},
		{"repository partitions", asset("run: [tool]", "partitions: dynamic:repo"), false},
		{"other placeholder", asset(`run: [tool, "{repo}"]`), false},
		{"embedded placeholder", asset(`run: [tool, "run-{run_id}"]`), false},
		{"stray brace", asset(`run: [tool, "{"]`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pack, err := LoadPack([]byte(tc.pack))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := pack.command("a"); (err == nil) != tc.ok {
				t.Fatalf("command err=%v, want ok=%v", err, tc.ok)
			}
		})
	}
	pack, err := LoadPack([]byte(asset("run: [tool]")))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pack.command("missing"); err == nil {
		t.Fatal("an asset the pack doesn't define has a command")
	}
}

func TestCommandResolvesRunID(t *testing.T) {
	id := uuid.New()
	got := command{"tool", "{run_id}", "--flag"}.resolve(id)
	if !slices.Equal(got, []string{"tool", id.String(), "--flag"}) {
		t.Fatalf("resolved %v", got)
	}
}

func TestPackPinsOrdersAndListsItsAssets(t *testing.T) {
	raw := []byte("pack: p\nassets:\n  b:\n    run: [tool]\n  a:\n    run: [tool]\n")
	pack, err := LoadPack(raw)
	if err != nil {
		t.Fatal(err)
	}
	pin := pack.Pin()
	if pin.Name != "pack:p" || pin.Digest != pack.digest || len(pin.Digest) != 64 {
		t.Fatalf("pin %+v", pin)
	}
	if got := pack.Assets(); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("assets %v", got)
	}
}
