package sqrl_test

import (
	"testing"

	"github.com/RaniSputnik/sqrl-go"
	"github.com/stretchr/testify/assert"
)

func TestTIFHas(t *testing.T) {
	testCases := []struct {
		Input  sqrl.TIF
		Test   sqrl.TIF
		Expect bool
	}{
		{
			Input:  sqrl.TIFIPMatch,
			Test:   sqrl.TIFIPMatch,
			Expect: true,
		},
		{
			Input:  sqrl.TIFClientFailure | sqrl.TIFCommandFailed,
			Test:   sqrl.TIFCommandFailed,
			Expect: true,
		},
		{
			Input:  sqrl.TIFBadIDAssociation | sqrl.TIFIPMatch,
			Test:   sqrl.TIFIPMatch,
			Expect: true,
		},
		{
			Input:  sqrl.TIFFunctionNotSupported,
			Test:   sqrl.TIFIPMatch,
			Expect: false,
		},
		{
			Input:  sqrl.TIF(0),
			Test:   sqrl.TIFSQRLDisabled,
			Expect: false,
		},
	}

	for _, testCase := range testCases {
		got := testCase.Input.Has(testCase.Test)
		assert.Equal(t, testCase.Expect, got)
	}
}
