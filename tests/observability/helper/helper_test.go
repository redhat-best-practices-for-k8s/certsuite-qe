package helper

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContainerSpecsFromStdoutBuffers(t *testing.T) {
	testCases := []struct {
		buffers          []string
		expectedLen      int
		expectedCommands []string
	}{
		{
			// Check with empty slice.
			buffers:     []string{},
			expectedLen: 0,
		},
		{
			// Check with one stdout buffer with just one line without new line char,
			// so only one container is expected."
			buffers:          []string{"Hello"},
			expectedLen:      1,
			expectedCommands: []string{`/bin/bash -c printf "Hello" && sleep INF`},
		},
		{
			// Check with one stdout buffer with just two lines, so only one container is expected."
			buffers:          []string{"Hello line 1\nHello line 2\n"},
			expectedLen:      1,
			expectedCommands: []string{`/bin/bash -c printf "Hello line 1\nHello line 2\n" && sleep INF`},
		},
		{
			// Check with two stdout buffers, one line each. The second one without a newline char.
			buffers:     []string{"Hello 1\n", "Hello 2"},
			expectedLen: 2,
			expectedCommands: []string{
				`/bin/bash -c printf "Hello 1\n" && sleep INF`,
				`/bin/bash -c printf "Hello 2" && sleep INF`,
			},
		},
		{
			// Check with two stdout buffers two lines each. Two containers expected.
			buffers: []string{
				"Container 1 Hello 1\nContainer 1 Hello 2\n",
				"Container 2 Hello 1\nContainer 2 Hello 2\n",
			},
			expectedLen: 2,
			expectedCommands: []string{
				`/bin/bash -c printf "Container 1 Hello 1\nContainer 1 Hello 2\n" && sleep INF`,
				`/bin/bash -c printf "Container 2 Hello 1\nContainer 2 Hello 2\n" && sleep INF`,
			},
		},
		{
			// Check with three stdout buffers. Last buffer has tab chars and no newline.
			// Three containers are expected.
			buffers: []string{
				"Container 1 Hello 1\n",
				"Container 2 Hello 1\nContainer 1 Hello 2\n",
				"Container 3 \tHello 1\t",
			},
			expectedLen: 3,
			expectedCommands: []string{
				`/bin/bash -c printf "Container 1 Hello 1\n" && sleep INF`,
				`/bin/bash -c printf "Container 2 Hello 1\nContainer 1 Hello 2\n" && sleep INF`,
				`/bin/bash -c printf "Container 3 \tHello 1\t" && sleep INF`,
			},
		},
	}

	for i := range testCases {
		t.Logf("UT %d", i)
		testCase := testCases[i]
		containers := createContainerSpecsFromStdoutBuffers(testCase.buffers)
		assert.Equal(t, testCase.expectedLen, len(containers))

		if testCase.expectedLen > 0 {
			for i := range containers {
				cmd := strings.Join(containers[i].Command, " ")
				assert.Equal(t, testCase.expectedCommands[i], cmd)
			}
		}
	}
}
