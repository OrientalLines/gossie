package utils_test

import (
	"reflect"
	"testing"

	"github.com/orientallines/gossie/pkg/utils"
)

func TestParseTagsSimple(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    []utils.Tag
		wantErr bool
	}{
		{
			name:  "valid tags",
			input: "gossie:foo=bar,baz=qux",
			want:  []utils.Tag{{Name: "foo", Value: "bar"}, {Name: "baz", Value: "qux"}},
		},
		{
			name:    "invalid format",
			input:   "invalid:foo=bar",
			wantErr: true,
		},
		{
			name:    "missing value",
			input:   "gossie:foo",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := utils.ParseTags(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseTags() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseTags() = %v, want %v", got, tc.want)
			}
		})
	}
}
