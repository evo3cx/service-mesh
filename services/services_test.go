package services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HeaderTracer(t *testing.T) {
	testCase := []struct {
		in     http.Header
		output http.Header
	}{
		{
			in: http.Header{
				"X-Client-Trace-Id": []string{"a"},
				"X-B3-TraceId":      []string{"a"},
			},
			output: http.Header{
				"X-Client-Trace-Id": []string{"a"},
				"X-B3-TraceId":      []string{"a"},
			},
		},

		{
			in: http.Header{},
			output: http.Header{
				"X-Client-Trace-Id": []string{"traceid"},
				"X-B3-TraceId":      []string{"b4"},
			},
		},
	}

	for i, tc := range testCase {
		act := headerTracer(tc.in, "test")
		tc.in.Set("X-B3-TraceId", "a")
		require.Equal(t, tc.output.Get("X-Client-Trace-Id"), act.Get("X-Client-Trace-Id"), fmt.Sprintf("Client index of %d", i))
		require.Equal(t, tc.output.Get("X-B3-TraceId"), act.Get("X-B3-TraceId"), fmt.Sprintf(" TraceId index of %d", i))
	}
}
