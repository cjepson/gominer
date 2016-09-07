/**
 * BLAKE256 14-round kernel
 */

#define ROTR(v,n) rotate(v,(uint)(32U-n))

// #ifdef _AMD_OPENCL
// #define SWAP(v)   rotate(v, 16U)
// #define ROTR8(v)  rotate(v, 24U)
// #else
// 
#define SWAP(v)  as_uint(as_uchar4(v).zwxy)
#define ROTR8(v) as_uint(as_uchar4(v).yzwx)
//
// from cuda
// #define SWAP(v)   as_uint(as_uchar4(v).xwzy)
// #define ROTR8(v)  as_uint(as_uchar4(v).wzyx)
//
// also might be
// #define SWAP(v)  as_uint(as_uchar4(v).yzwx)
// #define ROTR8(v) as_uint(as_uchar4(v).zwxy)
// #endif

#define pxorGS(a,b,c,d) { \
	v[a]+= xorLUT[i++] + v[b]; \
	v[d] = SWAP(v[d] ^ v[a]); \
	v[c]+= v[d]; \
	v[b] = ROTR(v[b] ^ v[c], 12); \
	v[a]+= xorLUT[i++] + v[b]; \
	v[d] = ROTR8(v[d] ^ v[a]); \
	v[c]+= v[d]; \
	v[b] = ROTR(v[b] ^ v[c], 7); \
}

#define pxorGS2(a,b,c,d,a1,b1,c1,d1) {\
	v[ a]+= xorLUT[i++] + v[ b]; \          
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = SWAP(v[ d] ^ v[ a]); \          
	v[d1] = SWAP(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 12); \
	v[b1] = ROTR(v[b1] ^ v[c1], 12); \
	v[ a]+= xorLUT[i++] + v[ b]; \           
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = ROTR8(v[ d] ^ v[ a]); \           
	v[d1] = ROTR8(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 7); \      
	v[b1] = ROTR(v[b1] ^ v[c1], 7); \
}

#define pxory1GS2(a,b,c,d,a1,b1,c1,d1) { \
	v[ a]+= xorLUT[i++] + v[ b]; \           
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = SWAP(v[ d] ^ v[ a]); \          
	v[d1] = SWAP(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 12); \     
	v[b1] = ROTR(v[b1] ^ v[c1], 12); \
	v[ a]+= xorLUT[i++] + v[ b]; \           
	v[a1]+= (xorLUT[i++]^nonce) + v[b1]; \
	v[ d] = ROTR8(v[ d] ^ v[ a]); \           
	v[d1] = ROTR8(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 7); \      
	v[b1] = ROTR(v[b1] ^ v[c1], 7); \
}

#define pxory0GS2(a,b,c,d,a1,b1,c1,d1) { \
	v[ a]+= xorLUT[i++] + v[ b]; \           
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = SWAP(v[ d] ^ v[ a]); \          
	v[d1] = SWAP(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 12); \     
	v[b1] = ROTR(v[b1] ^ v[c1], 12); \
	v[ a]+= (xorLUT[i++]^nonce) + v[ b]; \   
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = ROTR8(v[ d] ^ v[ a]); \           
	v[d1] = ROTR8(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 7); \      
	v[b1] = ROTR(v[b1] ^ v[c1], 7); \
}

#define pxorx1GS2(a,b,c,d,a1,b1,c1,d1) { \
	v[ a]+= xorLUT[i++] + v[ b]; \           
	v[a1]+= (xorLUT[i++]^nonce) + v[b1]; \
	v[ d] = SWAP(v[ d] ^ v[ a]); \          
	v[d1] = SWAP(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 12); \     
	v[b1] = ROTR(v[b1] ^ v[c1], 12); \
	v[ a]+= xorLUT[i++] + v[ b]; \           
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = ROTR8(v[ d] ^ v[ a]); \          
	v[d1] = ROTR8(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 7); \      
	v[b1] = ROTR(v[b1] ^ v[c1], 7); \
}

#define pxorx0GS2(a,b,c,d,a1,b1,c1,d1) { \
	v[ a]+= (xorLUT[i++]^nonce) + v[ b]; \   
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = SWAP(v[ d] ^ v[ a]); \	        
	v[d1] = SWAP(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 12); \     
	v[b1] = ROTR(v[b1] ^ v[c1], 12); \
	v[ a]+= xorLUT[i++] + v[ b]; \			
	v[a1]+= xorLUT[i++] + v[b1]; \
	v[ d] = ROTR8(v[ d] ^ v[ a]); \	        
	v[d1] = ROTR8(v[d1] ^ v[a1]); \
	v[ c]+= v[ d]; \                         
	v[c1]+= v[d1]; \
	v[ b] = ROTR(v[ b] ^ v[ c], 7); \		
	v[b1] = ROTR(v[b1] ^ v[c1], 7); \
}

__attribute__((reqd_work_group_size(WORKSIZE, 1, 1)))
__kernel void search(
	volatile __global uint * restrict output,
	
	// Precomputation state of v in the BLAKE256 block hash.
	const uint v0,
	const uint v1,
	const uint v2,
	const uint v3,
	const uint v4,
	const uint v5,
	const uint v6,
	const uint v7,
	const uint v8,
	const uint v9,
	const uint vA,
	const uint vB,
	const uint vC,
	const uint vD,
	const uint vE,
	const uint vF,
	
	// h[7], the last uint32 of the original midstate.
	const uint pre7,
	
	// Precomputed LUT of pre-XORed values.
	__constant uint *xorLUTGlobal
)
{	
	// Load the block state.
	__local uint v[16];
	v[ 0] = v0;
	v[ 1] = v1;
	v[ 2] = v2;
	v[ 3] = v3;
	v[ 4] = v4;
	v[ 5] = v5;
	v[ 6] = v6;
	v[ 7] = v7;

	v[ 8] = v8;
	v[ 9] = v9;
	v[10] = vA;
	v[11] = vB;
	v[12] = vC;
	v[13] = vD;
	v[14] = vE;
	v[15] = vF;

	const uint nonce = get_global_id(0);
	
	__local uint xorLUT[215];
	#pragma unroll
	for (uint i = 0; i<215; i++) {
		xorLUT[i] = xorLUTGlobal[i];
	}

	// 14 rounds.
	//
	// Modified first round accounting from precomputation excluding 
	// the nonce.
	v[ 1]+= (nonce ^ 0x13198A2E);
	v[13] = ROTR8(v[13] ^ v[1]);
	v[ 9]+= v[13];
	v[ 5] = ROTR(v[5] ^ v[9], 7);

	uint i = 0;
	v[ 1]+= xorLUT[i++]; // + v[ 6];
	v[ 0]+= v[5];
	v[12] = SWAP(v[12] ^ v[ 1]);         
	v[15] = SWAP(v[15] ^ v[ 0]);
	v[11]+= v[12];                        
	v[10]+= v[15];
	v[ 6] = ROTR(v[ 6] ^ v[11], 12);    
	v[ 5] = ROTR(v[5] ^ v[10], 12);
	v[ 1]+= xorLUT[i++] + v[ 6];          
	v[ 0]+= xorLUT[i++] + v[ 5];
	v[12] = ROTR8(v[12] ^ v[ 1]);          
	v[15] = ROTR8(v[15] ^ v[ 0]);
	v[11]+= v[12];                        
	v[10]+= v[15];
	v[ 6] = ROTR(v[ 6] ^ v[11], 7);     
	v[ 5] = ROTR(v[ 5] ^ v[10], 7);
	
	// Remaining 13 rounds.
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxory1GS2( 2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorx1GS2( 0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorx1GS2( 0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorx1GS2( 2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxory1GS2( 2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxory1GS2( 0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorx1GS2( 2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxory0GS2( 2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorx0GS2( 2, 7, 8, 13, 3, 4, 9, 14);
	pxory1GS2( 0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxory1GS2( 2, 7, 8, 13, 3, 4, 9, 14);
	pxorGS2(   0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorx1GS2( 0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS2(   2, 7, 8, 13, 3, 4, 9, 14);
	pxorx1GS2( 0, 4, 8, 12, 1, 5, 9, 13); 
	pxorGS2(   2, 6, 10, 14, 3, 7, 11, 15); 
	pxorGS2(   0, 5, 10, 15, 1, 6, 11, 12); 
	pxorGS(    2, 7, 8, 13);
		
	/* The final chunks of the hash
	 * are calculated as:
	 * h0 = h0 ^ V0 ^ V8;
	 * h1 = h1 ^ V1 ^ V9;
	 * h2 = h2 ^ V2 ^ VA;
	 * h3 = h3 ^ V3 ^ VB;
	 * h4 = h4 ^ V4 ^ VC;
	 * h5 = h5 ^ V5 ^ VD;
	 * h6 = h6 ^ V6 ^ VE;
	 * h7 = h7 ^ V7 ^ VF;
	 *
	 * We just check if the last byte
	 * is zeroed and if it is, we tell
	 * cgminer that we've found a
	 * and to check it against the
	 * target.
	*/

	/* Debug code to help you assess the correctness
	 * of your hashing function in case someone decides
	 * to try to optimize.
	if (!((pre7 ^ v[7] ^ v[15]) & 0xFFFF0000)) {
		printf("hash on gpu %x %x %x %x %x %x %x %x\n",
			h0 ^ V0 ^ V8,
			h1 ^ V1 ^ V9,
			h2 ^ V2 ^ VA,
			h3 ^ V3 ^ VB,
			h4 ^ V4 ^ VC,
			h5 ^ V5 ^ VD,
			h6 ^ V6 ^ VE,
			h7 ^ V7 ^ VF);
		printf("nonce for hash on gpu %x\n",
			nonce);
	}
	*/

	// There's something in the last 4 bytes of the hash, 
	// so this share is invalid.
	if ((pre7 ^ v[15]) != v[7]) return;
        
	// The share is valid. Push this share to our output 
	// buffer.
	output[++output[0]] = nonce;
}
