// precalculation.go continues precalculation optimizations for BLAKE256 as found
// originally in decred.cu by Alexis Provos/Tanguy Pruvot in ccminer.
package main

import "github.com/decred/gominer/util"

func rotateRight(v uint32, c uint32) uint32 {
	return (v >> c) | (v << (32 - c))
}

var z = [16]uint32{
	0x243F6A88, 0x85A308D3, 0x13198A2E, 0x03707344,
	0xA4093822, 0x299F31D0, 0x082EFA98, 0xEC4E6C89,
	0x452821E6, 0x38D01377, 0xBE5466CF, 0x34E90C6C,
	0xC0AC29B7, 0xC97C50DD, 0x3F84D5B5, 0xB5470917,
}

func precalcXors(m *[16]uint32, xorLUT *[215]uint32, i *int, x, y int) {
	xorLUT[*i] = (m[x] ^ z[y])
	*i++
	xorLUT[*i] = (m[y] ^ z[x])
	*i++
}

func precalcXors2(m *[16]uint32, xorLUT *[215]uint32, i *int, x, y, x1, y1 int) {
	xorLUT[*i] = (m[x] ^ z[y])
	*i++
	xorLUT[*i] = (m[x1] ^ z[y1])
	*i++
	xorLUT[*i] = (m[y] ^ z[x])
	*i++
	xorLUT[*i] = (m[y1] ^ z[x1])
	*i++
}

func precalculateStatesAndLUT(midstate [8]uint32, work [45]uint32) (*[2]uint32, *[16]uint32, *[215]uint32) {
	m := new([16]uint32)
	v := new([16]uint32)
	h := new([2]uint32)
	xorLUT := new([215]uint32)

	// Preconfigure the compression of the block.
	v[0] = midstate[0]
	v[1] = midstate[1]
	v[2] = midstate[2]
	v[3] = midstate[3]
	v[4] = midstate[4]
	v[5] = midstate[5]
	v[8] = midstate[6]

	v[12] = work[35]
	v[13] = midstate[7]

	util.Uint32EndiannessSwap(midstate[7])

	// Load the message.
	m[0] = work[32]
	m[1] = work[33]
	m[2] = work[34]
	m[3] = 0
	m[4] = work[36]
	m[5] = work[37]
	m[6] = work[38]
	m[7] = work[39]
	m[8] = work[40]
	m[9] = work[41]
	m[10] = work[42]
	m[11] = work[43]
	m[12] = work[44]
	m[13] = 0x80000001
	m[14] = 0
	m[15] = 0x000005a0

	// Save these for later use.
	h[0] = v[8]
	h[1] = v[13]

	// Start the first round and precalculate as much as possible.
	v[0] += (m[0] ^ z[1]) + v[4]
	v[12] = rotateRight(z[4]^0x5A0^v[0], 16)

	v[8] = z[0] + v[12]
	v[4] = rotateRight(v[4]^v[8], 12)
	v[0] += (m[1] ^ z[0]) + v[4]
	v[12] = rotateRight(v[12]^v[0], 8)
	v[8] += v[12]
	v[4] = rotateRight(v[4]^v[8], 7)

	v[1] += (m[2] ^ z[3]) + v[5]
	v[13] = rotateRight((z[5]^0x5A0)^v[1], 16)
	v[9] = z[1] + v[13]
	v[5] = rotateRight(v[5]^v[9], 12)
	v[1] += v[5] //+nonce ^ ...

	v[2] += (m[4] ^ z[5]) + h[0]
	v[14] = rotateRight(z[6]^v[2], 16)
	v[10] = z[2] + v[14]
	v[6] = rotateRight(h[0]^v[10], 12)
	v[2] += (m[5] ^ z[4]) + v[6]
	v[14] = rotateRight(v[14]^v[2], 8)
	v[10] += v[14]
	v[6] = rotateRight(v[6]^v[10], 7)

	v[3] += (m[6] ^ z[7]) + h[1]
	v[15] = rotateRight(z[7]^v[3], 16)
	v[11] = z[3] + v[15]
	v[7] = rotateRight(h[1]^v[11], 12)
	v[3] += (m[7] ^ z[6]) + v[7]
	v[15] = rotateRight(v[15]^v[3], 8)
	v[11] += v[15]
	v[7] = rotateRight(v[11]^v[7], 7)
	v[0] += m[8] ^ z[9]

	// Generate the look up table for the XOR values.
	i := new(int)
	precalcXors(m, xorLUT, i, 10, 11)
	xorLUT[0] += v[6]
	xorLUT[*i] = (m[9] ^ z[8])
	*i++

	precalcXors2(m, xorLUT, i, 12, 13, 14, 15)
	precalcXors2(m, xorLUT, i, 14, 10, 4, 8)
	precalcXors2(m, xorLUT, i, 9, 15, 13, 6)
	precalcXors2(m, xorLUT, i, 1, 12, 0, 2)
	precalcXors2(m, xorLUT, i, 11, 7, 5, 3)
	precalcXors2(m, xorLUT, i, 11, 8, 12, 0)
	precalcXors2(m, xorLUT, i, 5, 2, 15, 13)
	precalcXors2(m, xorLUT, i, 10, 14, 3, 6)
	precalcXors2(m, xorLUT, i, 7, 1, 9, 4)
	precalcXors2(m, xorLUT, i, 7, 9, 3, 1)
	precalcXors2(m, xorLUT, i, 13, 12, 11, 14)
	precalcXors2(m, xorLUT, i, 2, 6, 5, 10)
	precalcXors2(m, xorLUT, i, 4, 0, 15, 8)
	precalcXors2(m, xorLUT, i, 9, 0, 5, 7)
	precalcXors2(m, xorLUT, i, 2, 4, 10, 15)
	precalcXors2(m, xorLUT, i, 14, 1, 11, 12)
	precalcXors2(m, xorLUT, i, 6, 8, 3, 13)
	precalcXors2(m, xorLUT, i, 2, 12, 6, 10)
	precalcXors2(m, xorLUT, i, 0, 11, 8, 3)
	precalcXors2(m, xorLUT, i, 4, 13, 7, 5)
	precalcXors2(m, xorLUT, i, 15, 14, 1, 9)
	precalcXors2(m, xorLUT, i, 12, 5, 1, 15)
	precalcXors2(m, xorLUT, i, 14, 13, 4, 10)
	precalcXors2(m, xorLUT, i, 0, 7, 6, 3)
	precalcXors2(m, xorLUT, i, 9, 2, 8, 11)
	precalcXors2(m, xorLUT, i, 13, 11, 7, 14)
	precalcXors2(m, xorLUT, i, 12, 1, 3, 9)
	precalcXors2(m, xorLUT, i, 5, 0, 15, 4)
	precalcXors2(m, xorLUT, i, 8, 6, 2, 10)
	precalcXors2(m, xorLUT, i, 6, 15, 14, 9)
	precalcXors2(m, xorLUT, i, 11, 3, 0, 8)
	precalcXors2(m, xorLUT, i, 12, 2, 13, 7)
	precalcXors2(m, xorLUT, i, 1, 4, 10, 5)
	precalcXors2(m, xorLUT, i, 10, 2, 8, 4)
	precalcXors2(m, xorLUT, i, 7, 6, 1, 5)
	precalcXors2(m, xorLUT, i, 15, 11, 9, 14)
	precalcXors2(m, xorLUT, i, 3, 12, 13, 0)
	precalcXors2(m, xorLUT, i, 0, 1, 2, 3)
	precalcXors2(m, xorLUT, i, 4, 5, 6, 7)
	precalcXors2(m, xorLUT, i, 8, 9, 10, 11)
	precalcXors2(m, xorLUT, i, 12, 13, 14, 15)
	precalcXors2(m, xorLUT, i, 14, 10, 4, 8)
	precalcXors2(m, xorLUT, i, 9, 15, 13, 6)
	precalcXors2(m, xorLUT, i, 1, 12, 0, 2)
	precalcXors2(m, xorLUT, i, 11, 7, 5, 3)
	precalcXors2(m, xorLUT, i, 11, 8, 12, 0)
	precalcXors2(m, xorLUT, i, 5, 2, 15, 13)
	precalcXors2(m, xorLUT, i, 10, 14, 3, 6)
	precalcXors2(m, xorLUT, i, 7, 1, 9, 4)
	precalcXors2(m, xorLUT, i, 7, 9, 3, 1)
	precalcXors2(m, xorLUT, i, 13, 12, 11, 14)
	precalcXors2(m, xorLUT, i, 2, 6, 5, 10)

	precalcXors(m, xorLUT, i, 4, 0)
	precalcXors(m, xorLUT, i, 15, 8)

	return h, v, xorLUT
}
