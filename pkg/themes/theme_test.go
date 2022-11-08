package themes

import (
	"reflect"
	"testing"

	"github.com/matm/gocov-html/pkg/types"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		want  types.Beautifier
	}{
		{"empty string", "", nil},
		{"unknown", "bad", nil},
		{"default", "golang", defaultTheme{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.theme); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name  string
		after func([]types.Beautifier)
	}{
		{"all themes", func(ts []types.Beautifier) {
			if len(ts) == 0 {
				t.Error("no theme")
			}
			for _, p := range ts {
				if p.Name() == "" {
					t.Error("empty name")
				}
				if p.Description() == "" {
					t.Error("empty description")
				}
				if p.Template() == nil {
					t.Errorf("missing template for %q theme", p.Name())
				}
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := List()
			if tt.after != nil {
				tt.after(got)
			}
		})
	}
}
