package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

// TestGJSONQuery tests the gjson query syntax, which assumes a hacky conversion of each line into a json array.
// https://github.com/tidwall/gjson/blob/master/gjson_test.go#L1648
func TestGJSONQuery(t *testing.T) {
	type queryCase struct {
		line   string
		query  string
		exists bool
	}
	lines := []string{
		`{"type":"Feature","id":1,"geometry":{"type":"Point","coordinates":[-93.25535583496094,44.98938751220703]},"properties":{"Accuracy":16.58573341369629,"Activity":"Stationary","Elevation":266.4858703613281,"Heading":349.74163818359375,"Name":"Rye13","Pressure":94.55844116210938,"Speed":0,"Time":"2023-12-08T10:04:09.017Z","UUID":"05C63745-BFA3-4DE3-AF2F-CDE2173C0E11","UnixTime":1702029849,"Version":"V.customizableCatTrackHat"}}
`,
		`{"type":"Feature","id":1,"geometry":{"type":"Point","coordinates":[-93.25735583496094,44.98638751220703]},"properties":{"Accuracy":101.7849,"Activity":"Running","Elevation":256.4858703613281,"Heading":347.74163818359375,"Name":"Rye13","Pressure":95.55844116210938,"Speed":0,"Time":"2023-12-08T10:04:10.017Z","UUID":"05C63745-BFA3-4DE3-AF2F-CDE2173C0E11","UnixTime":1702029850,"Version":"V.customizableCatTrackHat"}
`,
	}
	queryCases := []queryCase{
		{lines[0], `#[properties.Accuracy<100]`, true},
		{lines[1], `#[properties.Accuracy<100]`, false},

		{lines[0], `#[properties.Activity="Running"]`, false},
		{lines[1], `#[properties.Activity="Running"]`, true},

		/*
			Please note that prior to v1.3.0, queries used the #[...] brackets. This was changed in v1.3.0 as to avoid confusion with the new multipath syntax. For backwards compatibility, #[...] will continue to work until the next major release.
			https://github.com/tidwall/gjson/blob/master/SYNTAX.md#queries
		*/
		{lines[0], `#(properties.Accuracy<100)`, true},
		{lines[1], `#(properties.Accuracy<100)`, false},

		{lines[0], `#(properties.Activity="Running")`, false},
		{lines[1], `#(properties.Activity="Running")`, true},
	}

	for _, c := range queryCases {
		result := gjson.Get(fmt.Sprintf("[%s]", c.line), c.query)
		t.Logf("%s: %q (exists=%v)", c.query, result, result.Exists())
		if result.Exists() != c.exists {
			t.Errorf("expected %v, got %v", c.exists, result.Exists())
		}
	}
}

func TestFilter(t *testing.T) {
	lines := []string{
		`{"type":"Feature","id":1,"geometry":{"type":"Point","coordinates":[-93.25535583496094,44.98938751220703]},"properties":{"Accuracy":16.58573341369629,"Activity":"Stationary","Elevation":266.4858703613281,"Heading":349.74163818359375,"Name":"Rye13","Pressure":94.55844116210938,"Speed":0,"Time":"2023-12-08T10:04:09.017Z","UUID":"05C63745-BFA3-4DE3-AF2F-CDE2173C0E11","UnixTime":1702029849,"Version":"V.customizableCatTrackHat"}}
`,
		`{"type":"Feature","id":1,"geometry":{"type":"Point","coordinates":[-93.25735583496094,44.98638751220703]},"properties":{"Accuracy":101.7849,"Activity":"Running","Elevation":256.4858703613281,"Heading":347.74163818359375,"Name":"Rye13","Pressure":95.55844116210938,"Speed":0,"Time":"2023-12-08T10:04:10.017Z","UUID":"05C63745-BFA3-4DE3-AF2F-CDE2173C0E11","UnixTime":1702029850,"Version":"V.customizableCatTrackHat"}
`,
	}
	type queryCase struct {
		line                          string
		matchAll, matchAny, matchNone []string
		err                           error
	}
	queryCases := []queryCase{
		{line: lines[0], matchAll: []string{`#(properties.Accuracy<100)`}, err: nil},
		{line: lines[1], matchAll: []string{`#(properties.Accuracy<100)`}, err: errInvalidMatchAll},
	}
	for _, c := range queryCases {
		err := filter([]byte(c.line), c.matchAll, c.matchAny, c.matchNone)
		if !errors.Is(err, c.err) {
			t.Errorf("expected %v, got %v", c.err, err)
		}
	}
}
