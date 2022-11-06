package themes

import (
	"reflect"
	"testing"

	"github.com/matm/gocov-html/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert := assert.New(t)
	require := require.New(t)
	tests := []struct {
		name  string
		after func([]types.Beautifier)
	}{
		{"all themes", func(ts []types.Beautifier) {
			assert.NotEmpty(ts)
			for _, p := range ts {
				assert.NotEmpty(p.Name(), "empty name")
				assert.NotEmpty(p.Description(), "empty description")
				z, err := p.Template()
				require.NoError(err)
				assert.NotNil(z, "missing template for %q theme", p.Name())
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
