package filter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ng-vu/githubx/backend/pkg/dot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIDsJSON(t *testing.T) {
	var nilSlice IDs
	tests := []struct {
		name       string
		input      string
		expectBare IDs
		expectPtr  *IDs
		error      string
	}{
		{
			"null",
			`null`,
			nil,
			nil,
			"",
		},
		{
			"empty",
			`""`,
			nil,
			&nilSlice, // pointer to nil slice
			"",
		},
		{
			"simple",
			`"1234"`,
			IDs{1234},
			&IDs{1234},
			"",
		},
		{
			"list",
			`"1234,5678,9012"`,
			IDs{1234, 5678, 9012},
			&IDs{1234, 5678, 9012},
			"",
		},
		{
			"empty array (invalid)",
			`[]`,
			nil,
			nil,
			"json: can not read `[]` as id",
		},
		{
			"array of string",
			`["1234","5678","9012"]`,
			IDs{1234, 5678, 9012},
			&IDs{1234, 5678, 9012},
			"",
		},
		{
			"array of int",
			`[1234,5678,9012]`,
			IDs{1234, 5678, 9012},
			&IDs{1234, 5678, 9012},
			"",
		},
		{
			"zero (invalid)",
			`0`,
			nil,
			nil,
			"json: can not read `0` as id",
		},
		{
			"zero as string (invalid)",
			`"0"`,
			nil,
			nil,
			"json: can not read `\"0\"` as id",
		},
		{
			"list with zero (invalid)",
			`"1234,0,5678"`,
			nil,
			nil,
			"json: can not read `\"1234,0,5678\"` as id",
		},
		{
			"list with empty (invalid)",
			`"1234,"`,
			nil,
			nil,
			"json: can not read `\"1234,\"` as id",
		},
		{
			"list with space",
			`"1234, 5678"`,
			IDs{1234, 5678},
			&IDs{1234, 5678},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{
				var bare IDs
				err := json.Unmarshal([]byte(tt.input), &bare)
				if tt.error != "" {
					assert.EqualError(t, err, tt.error)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tt.expectBare, bare)
			}
			{
				var ptr *IDs
				err := json.Unmarshal([]byte(tt.input), &ptr)
				if tt.error != "" {
					assert.EqualError(t, err, tt.error)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tt.expectPtr, ptr)
			}
		})
	}
}

func getDay(year int, month time.Month, day int) dot.Time {
	return dot.Time(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func TestDateJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect Date
		error  string
	}{
		{
			"null",
			`null`,
			Date{},
			"",
		},
		{
			"empty string",
			`""`,
			Date{},
			"",
		},
		{
			"empty object",
			`{}`,
			Date{},
			"",
		},
		{
			"empty object with empty string",
			`{"from":"","to":""}`,
			Date{},
			"",
		},
		{
			"empty object with null",
			`{"from":null,"to":null}`,
			Date{},
			"",
		},
		{
			"single date",
			`"2020-01-01T00:00:00Z"`,
			Date{
				From: getDay(2020, 1, 1),
				To:   dot.Time{},
			},
			"",
		},
		{
			"range",
			`"2020-01-01T00:00:00Z,2020-02-01T00:00:00Z"`,
			Date{
				From: getDay(2020, 1, 1),
				To:   getDay(2020, 2, 1),
			},
			"",
		},
		{
			"first",
			`"2020-01-01T00:00:00Z,"`,
			Date{
				From: getDay(2020, 1, 1),
				To:   dot.Time{},
			},
			"",
		},
		{
			"last",
			`",2020-01-01T00:00:00Z"`,
			Date{
				From: dot.Time{},
				To:   getDay(2020, 1, 1),
			},
			"",
		},
		{
			"no date (invalid)",
			`","`,
			Date{
				From: dot.Time{},
				To:   getDay(2020, 1, 1),
			},
			"json: can not read `\",\"` as date",
		},
		{
			"too many dates (invalid)",
			`"2020-01-01T00:00:00Z,2020-01-01T00:00:00Z,2020-01-01T00:00:00Z"`,
			Date{},
			"json: can not read `\"2020-01-01T00:00:00Z,2020-01-01T00:00:00Z,2020-01-01T00:00:00Z\"` as date",
		},
		{
			"date range is invalid (invalid)",
			`"2020-02-01T00:00:00Z,2020-01-01T00:00:00Z"`,
			Date{},
			"date range [2020-02-01T00:00:00Z,2020-01-01T00:00:00Z) is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output Date
			err := json.Unmarshal([]byte(tt.input), &output)
			if tt.error != "" {
				assert.EqualError(t, err, tt.error)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expect, output)
		})
	}
}
