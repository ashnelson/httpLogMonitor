package stats

import (
	"reflect"
	"testing"
)

type logLineTest struct {
	logLine  string
	expected []string
}

type reqStrTest struct {
	request  string
	expected string
}

func TestSplitLogLine(t *testing.T) {
	tests := []logLineTest{
		logLineTest{
			logLine:  `127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123`,
			expected: []string{"127.0.0.1", "-", "james", "[09/May/2018:16:00:39 +0000]", "GET /report HTTP/1.0", "200", "123"},
		},
		logLineTest{
			logLine:  `127.0.0.1 - jill [09/May/2018:16:00:41 +0000] "GET /api/user HTTP/1.0" 200 234`,
			expected: []string{"127.0.0.1", "-", "jill", "[09/May/2018:16:00:41 +0000]", "GET /api/user HTTP/1.0", "200", "234"},
		},
		logLineTest{
			logLine:  `127.0.0.1 - frank [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 200 34`,
			expected: []string{"127.0.0.1", "-", "frank", "[09/May/2018:16:00:42 +0000]", "POST /api/user HTTP/1.0", "200", "34"},
		},
		logLineTest{
			logLine:  `127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12`,
			expected: []string{"127.0.0.1", "-", "mary", "[09/May/2018:16:00:42 +0000]", "POST /api/user HTTP/1.0", "503", "12"},
		},
	}

	var actual []string
	for i, tst := range tests {
		actual = splitLogLine(tst.logLine)

		if !reflect.DeepEqual(actual, tst.expected) {
			t.Fatalf("Test %d failed; invalid split log line\n\tExpected: %v\n\tActual:   %v", i, tst.expected, actual)
		}
	}
}

func TestGetSectionName(t *testing.T) {
	tests := []reqStrTest{
		reqStrTest{
			request:  "GET /report HTTP/1.0",
			expected: "/report",
		},
		reqStrTest{
			request:  "GET /api/user HTTP/1.0",
			expected: "/api",
		},
		reqStrTest{
			request:  "POST /api/user HTTP/1.0",
			expected: "/api",
		},
		reqStrTest{
			request:  "POST /api/user HTTP/1.0",
			expected: "/api",
		},
	}

	var actual string
	for i, tst := range tests {
		actual = getSectionName(tst.request)

		if actual != tst.expected {
			t.Fatalf("Test %d failed; invalid section name\n\tExpected: %v\n\tActual:   %v", i, tst.expected, actual)
		}
	}
}

func TestGetMethod(t *testing.T) {
	tests := []reqStrTest{
		reqStrTest{
			request:  "GET /report HTTP/1.0",
			expected: "GET",
		},
		reqStrTest{
			request:  "DELETE /api/user HTTP/1.0",
			expected: "DELETE",
		},
		reqStrTest{
			request:  "POST /api/user HTTP/1.0",
			expected: "POST",
		},
		reqStrTest{
			request:  "PATCH /api/user HTTP/1.0",
			expected: "PATCH",
		},
	}

	var actual string
	for i, tst := range tests {
		actual = getRequestMethod(tst.request)

		if actual != tst.expected {
			t.Fatalf("Test %d failed; invalid request method\n\tExpected: %v\n\tActual:   %v", i, tst.expected, actual)
		}
	}
}
