// Copyright (c) 2014 Casey Marshall. See LICENSE file for details.

package basen_test

import (
	"crypto/rand"
	"testing"

	gc "launchpad.net/gocheck"

	"github.com/cmars/basen"
)

func Test(t *testing.T) { gc.TestingT(t) }

type Suite struct{}

var _ = gc.Suite(&Suite{})

func (s *Suite) TestRoundTrip62(c *gc.C) {
	testCases := []struct {
		val []byte
		b62 string
	}{
		{[]byte{1}, "1"},
		{[]byte{61}, "z"},
		{[]byte{62}, "10"},
	}

	for _, testCase := range testCases {
		b62 := basen.StdBase62.EncodeToString(testCase.val)
		c.Check(b62, gc.Equals, testCase.b62)

		val, err := basen.StdBase62.DecodeString(testCase.b62)
		c.Assert(err, gc.IsNil)
		c.Check(val, gc.DeepEquals, testCase.val, gc.Commentf("%s", testCase.b62))
	}
}

func (s *Suite) TestRand256(c *gc.C) {
	for i := 0; i < 100; i++ {
		v := basen.StdBase62.MustRandom(rand.Reader, 32)
		// Should be 43 chars or less because math.log(2**256, 62) == 42.994887413002736
		c.Assert(len(v) < 44, gc.Equals, true)
	}
}
