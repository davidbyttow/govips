package vips

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveKeep(t *testing.T) {
	allButGainmap := int(ForeignKeepAll &^ ForeignKeepGainmap)

	tests := []struct {
		name         string
		keep         ForeignKeep
		major, minor int
		wantFlags    int
		wantSet      bool
		wantErr      error // nil, errKeepUnsupported, or errAny
	}{
		{"unset on 8.14", 0, 8, 14, 0, false, nil},
		{"unset on 8.18", 0, 8, 18, 0, false, nil},
		{"set on 8.14", ForeignKeepIcc, 8, 14, 0, false, errKeepUnsupported},
		{"set on 7.x", ForeignKeepIcc, 7, 99, 0, false, errKeepUnsupported},
		{"none maps to zero", ForeignKeepNone, 8, 15, 0, true, nil},
		{"none ored with icc", ForeignKeepNone | ForeignKeepIcc, 8, 15, int(ForeignKeepIcc), true, nil},
		{"icc", ForeignKeepIcc, 8, 15, int(ForeignKeepIcc), true, nil},
		{"unknown bit", ForeignKeep(1 << 6), 8, 18, 0, false, errAny},
		{"all on 8.17 drops gainmap", ForeignKeepAll, 8, 17, allButGainmap, true, nil},
		{"all on 8.18", ForeignKeepAll, 8, 18, int(ForeignKeepAll), true, nil},
		{"all on 9.0", ForeignKeepAll, 9, 0, int(ForeignKeepAll), true, nil},
		{"gainmap only on 8.17", ForeignKeepGainmap, 8, 17, 0, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, set, err := resolveKeep(tt.keep, tt.major, tt.minor)
			switch {
			case tt.wantErr == errAny:
				require.Error(t, err)
				assert.False(t, errors.Is(err, errKeepUnsupported))
			case tt.wantErr != nil:
				// testify v1.6.1 (go.mod) has no ErrorIs; repo tests use errors.Is.
				require.True(t, errors.Is(err, tt.wantErr), "want %v, got: %v", tt.wantErr, err)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.wantFlags, flags)
				assert.Equal(t, tt.wantSet, set)
			}
		})
	}
}

// errAny marks table rows that expect some error other than errKeepUnsupported.
var errAny = errors.New("any error")
