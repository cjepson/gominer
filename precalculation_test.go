// precalculation_test.go
package main

import (
	"testing"

	"github.com/decred/gominer/blake256"
	"github.com/decred/gominer/util"
)

// TestPreCalc
func TestPreCalc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header [45]uint32
		m      [16]uint32
		h      [2]uint32
		v      [16]uint32
		xorLUT [215]uint32
	}{
		{
			name: "test",
			header: [45]uint32{
				0x1, 0x2c4b84ed, 0xc0cc8bc0, 0xa37efc5c, 0x2e39633f, 0x5ad3cc9, 0x8ccd16b7, 0x1b1, 0x0, 0xdd0f8c7, 0x60da8d5, 0x41a16b0d, 0x5ff32a17, 0xb4e4c924, 0xb83f1486, 0xbd22b3ba, 0x73970d02, 0xf0f694c4, 0x17b9117d, 0x3930e215, 0x6575a084, 0x362f92fb, 0x43e63088, 0x9ba94afb, 0xe9693f2f, 0x53e20001, 0xc3b8ca8a, 0x5, 0xa0fe, 0x1a11fa87, 0x21149ca7, 0x1, 0xeeb0, 0xd56, 0x57cf10a5, 0x0, 0x2, 0xec8f1800, 0x45bd01e3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				/*
					0x1, 0x2c4b84ed, 0xc0cc8bc0, 0xa37efc5c, 0x2e39633f, 0x5ad3cc9,
					0x8ccd16b7, 0x1b1, 0x0, 0xdd0f8c7, 0x60da8d5, 0x41a16b0d, 0x5ff32a17,
					0xb4e4c924, 0xb83f1486, 0xbd22b3ba, 0x73970d02, 0xf0f694c4,
					0x17b9117d, 0x3930e215, 0x6575a084, 0x362f92fb, 0x43e63088,
					0x9ba94afb, 0xe9693f2f, 0x53e20001, 0xc3b8ca8a, 0x5, 0xa0fe,
					0x1a11fa87, 0x21149ca7, 0x1, 0xeeb0, 0xd56, 0x57cf10a5, 0x0, 0x2,
					0xec8f1800, 0x45bd01e3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				*/
			},
			m: [16]uint32{0xb0ee0000, 0x560d0000, 0xa510cf57, 0x0, 0x2000000,
				0x188fec, 0xe301bd45, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80000001,
				0x0, 0x5a0},
			h: [2]uint32{0x45223e33, 0x355bca51},
			v: [16]uint32{0xdfa9114b, 0x4a16cbc4, 0x68c162fb, 0x3ab9670f, 0x3d4106bb,
				0x67b6667e, 0x43e24956, 0x4a074d7c, 0xba1ec8cc, 0xd6d9d4d2,
				0x53f47069, 0xc09f1b46, 0xd566a043, 0x5136cbff, 0xb5dbc4b,
				0x119124e4},
			xorLUT: [215]uint32{0x78cb55c2, 0xbe5466cf, 0x452821e6, 0xc97c50dd,
				0xb5470917, 0x40ac29b6, 0x3f84d015, 0xbe5466cf, 0x472821e6, 0x3f84d5b5,
				0xa4093822, 0xb5470917, 0x882efa99, 0x38d016d7, 0x2a7ded98, 0x96a129b7,
				0xa3f78a2e, 0x85a308d3, 0x812fa5df, 0xec4e6c89, 0x368fca8, 0x34e90c6c,
				0x299f31d0, 0x452821e6, 0x243f6a88, 0x34e90c6c, 0x704229b7, 0x130105c2,
				0xc97c557d, 0x8c8ffe87, 0x35470916, 0x3f84d5b5, 0x82efa98, 0xbe5466cf,
				0xe071ce01, 0x85a308d3, 0xa4093822, 0xba436c89, 0x3ad01377, 0x38d01377,
				0x85a308d3, 0xec4e6c89, 0x557d7344, 0x40ac29b6, 0x3f84d5b5, 0xc97c50dd,
				0x34e90c6c, 0xad3e35cf, 0xbe4ce923, 0xf018376b, 0x299f31d0, 0x263f6a88,
				0x45282446, 0x14e73822, 0xb5470917, 0x243f6a88, 0xec56e365, 0x883e1377,
				0x299f31d0, 0x119f775, 0xb5470917, 0x11198a2e, 0xbe54636f, 0x85a308d3,
				0xc0ac29b7, 0x6989d5b5, 0x34e90c6c, 0xa6299ca3, 0xc97c50dd, 0x82efa98,
				0x83707345, 0x65bce6e0, 0x5d55db8a, 0x13198a2e, 0x82efa98, 0x84070c6c,
				0x3707344, 0x243f6a88, 0x452821e6, 0xcb7c50dd, 0x299f31d0, 0x24093823,
				0xec56e365, 0x3f84d015, 0x6edd1377, 0xb5470917, 0x85a308d3, 0x299f31d0,
				0xe34a0917, 0xc0b4a65b, 0x85a30d73, 0xc97c50dd, 0xbc5466cf, 0xbf84d5b4,
				0xa4093822, 0x5ca06c89, 0xe071ce01, 0x243f6a88, 0x82efa98, 0x13198a2e,
				0x34e90c6c, 0x9dc0dc20, 0x452821e6, 0xb4e90c6d, 0x3f84d5b5, 0xc97c50dd,
				0xec4e6c89, 0x85a308d3, 0x38d01377, 0x96a129b7, 0x3707344, 0x2427e564,
				0xa4093d82, 0x997131d0, 0xb7470917, 0x82efa98, 0x1b44a998, 0xa6299ca3,
				0x13198a2e, 0x5646b452, 0x38d01377, 0x82eff38, 0x3f84d5b5, 0x3707344,
				0xf5c621e6, 0x34e90c6c, 0x243f6a88, 0x13198a2e, 0x6c4e6c88, 0x65bce6e0,
				0xc97c50dd, 0xf2043822, 0x299f31d0, 0x87a308d3, 0xbe4ce923, 0x13198a2e,
				0xa4093822, 0x1b44a998, 0x472821e6, 0x82efa98, 0x7f9231d0, 0xf4fd1cc,
				0x85bb873f, 0x34e909cc, 0x3f84d5b5, 0xb5470917, 0x38d01377, 0xc0ac29b7,
				0xa43f6a89, 0x3707344, 0x799250dd, 0x354d08d3, 0xa660bc13, 0x72326a88,
				0x13198a2e, 0x2b9f31d0, 0xf4fd1cc, 0xa411b7ce, 0x82efa98, 0x38d01377,
				0x34e90c6c, 0x452821e6, 0xbe5466cf, 0xc97c50dd, 0xb5470917, 0x40ac29b6,
				0x3f84d015, 0xbe5466cf, 0x472821e6, 0x3f84d5b5, 0xa4093822, 0xb5470917,
				0x882efa99, 0x38d016d7, 0x2a7ded98, 0x96a129b7, 0xa3f78a2e, 0x85a308d3,
				0x812fa5df, 0xec4e6c89, 0x368fca8, 0x34e90c6c, 0x299f31d0, 0x452821e6,
				0x243f6a88, 0x34e90c6c, 0x704229b7, 0x130105c2, 0xc97c557d, 0x8c8ffe87,
				0x35470916, 0x3f84d5b5, 0x82efa98, 0xbe5466cf, 0xe071ce01, 0x85a308d3,
				0xa4093822, 0xba436c89, 0x3ad01377, 0x38d01377, 0x85a308d3, 0xec4e6c89,
				0x557d7344, 0x40ac29b6, 0x3f84d5b5, 0xc97c50dd, 0x34e90c6c, 0xad3e35cf,
				0xbe4ce923, 0xf018376b, 0x299f31d0, 0x263f6a88, 0x14e73822, 0x45282446,
				0xb5470917},
		},
	}

	for _, test := range tests {
		var midstate [8]uint32
		copy(midstate[:], blake256.IV256[:])

		// Hash the two first blocks.
		headerAsBytes := util.ConvertUint32SliceHeaderToByteSlice(test.header)
		blake256.Block(midstate[:], headerAsBytes[0:64], 512)
		blake256.Block(midstate[:], headerAsBytes[64:128], 1024)

		h, v, xorLUT := precalculateStatesAndLUT(midstate, test.header)
		if *h != test.h {
			t.Errorf("got %x, want %x", h, test.h)
		}
		if *v != test.v {
			t.Errorf("got %x, want %x", v, test.v)
		}
		if *xorLUT != test.xorLUT {
			t.Errorf("got %x, want %x", xorLUT, test.xorLUT)
		}
	}
}
