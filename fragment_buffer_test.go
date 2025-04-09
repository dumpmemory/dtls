// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package dtls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFragmentBuffer(t *testing.T) {
	for _, test := range []struct {
		Name     string
		In       [][]byte
		Expected [][]byte
		Epoch    uint16
	}{
		{
			Name: "Single Fragment",
			In: [][]byte{
				{
					0x16, 0xfe, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xfe, 0xff, 0x00,
				},
			},
			Expected: [][]byte{
				{0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xfe, 0xff, 0x00},
			},
			Epoch: 0,
		},
		{
			Name: "Single Fragment Epoch 3",
			In: [][]byte{
				{
					0x16, 0xfe, 0xff, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x03,
					0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xfe, 0xff, 0x00,
				},
			},
			Expected: [][]byte{
				{0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xfe, 0xff, 0x00},
			},
			Epoch: 3,
		},
		{
			Name: "Multiple Fragments",
			In: [][]byte{
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x0b, 0x00,
					0x00, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04,
				},
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x0b, 0x00,
					0x00, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x05, 0x05, 0x06, 0x07, 0x08, 0x09,
				},
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x0b, 0x00,
					0x00, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x05, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E,
				},
			},
			Expected: [][]byte{
				{
					0x0b, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x01,
					0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e,
				},
			},
			Epoch: 0,
		},
		{
			Name: "Multiple Unordered Fragments",
			In: [][]byte{
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x0b, 0x00,
					0x00, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04,
				},
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x0b, 0x00,
					0x00, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x05, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E,
				},
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x81, 0x0b, 0x00,
					0x00, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x05, 0x05, 0x06, 0x07, 0x08, 0x09,
				},
			},
			Expected: [][]byte{
				{
					0x0b, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x01,
					0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e,
				},
			},
			Epoch: 0,
		},
		{
			Name: "Multiple Handshakes in Single Fragment",
			In: [][]byte{
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x30, /* record header */
					0x03, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xfe, 0xff, 0x01, 0x01, /*handshake msg 1*/
					0x03, 0x00, 0x00, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xfe, 0xff, 0x01, 0x01, /*handshake msg 2*/
					0x03, 0x00, 0x00, 0x04, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xfe, 0xff, 0x01, 0x01, /*handshake msg 3*/
				},
			},
			Expected: [][]byte{
				{0x03, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xfe, 0xff, 0x01, 0x01},
				{0x03, 0x00, 0x00, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xfe, 0xff, 0x01, 0x01},
				{0x03, 0x00, 0x00, 0x04, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xfe, 0xff, 0x01, 0x01},
			},
			Epoch: 0,
		},
		// Assert that a zero length fragment doesn't cause the fragmentBuffer to enter an infinite loop
		{
			Name: "Zero Length Fragment",
			In: [][]byte{
				{
					0x16, 0xfe, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0d, 0x00, 0x00,
					0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			Expected: [][]byte{
				{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00},
			},
			Epoch: 0,
		},
	} {
		fragmentBuffer := newFragmentBuffer()
		for _, frag := range test.In {
			status, _, err := fragmentBuffer.push(frag)
			assert.NoError(t, err)
			assert.Truef(t, status, "fragmentBuffer didn't accept fragments for '%s'", test.Name)
		}

		for _, expected := range test.Expected {
			out, epoch := fragmentBuffer.pop()
			assert.Equalf(t, expected, out, "fragmentBuffer '%s' pop should return expected output", test.Name)
			assert.Equalf(t, test.Epoch, epoch, "fragmentBuffer returend wrong epoch")
		}

		frag, _ := fragmentBuffer.pop()
		assert.Nilf(t, frag, "fragmentBuffer '%s' pop should return nil when no more fragments are available", test.Name)
	}
}

func TestFragmentBuffer_Overflow(t *testing.T) {
	fragmentBuffer := newFragmentBuffer()

	// Push a buffer that doesn't exceed size limits
	_, _, err := fragmentBuffer.push([]byte{
		0x16, 0xfe, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x03,
		0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xfe, 0xff, 0x00,
	})
	assert.NoError(t, err)

	// Allocate a buffer that exceeds cache size
	largeBuffer := make([]byte, fragmentBufferMaxSize)
	_, _, err = fragmentBuffer.push(largeBuffer)
	assert.ErrorIs(t, err, errFragmentBufferOverflow, "Pushing a large buffer should return an overflow error")
}
