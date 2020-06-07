#ifndef MINE_CUH
#define MINE_CUH

#include<cuda.h>
#include<sha256.cuh>
#include<cmine.h>
#include<stdint.h>
#include<stdlib.h>
#include<inttypes.h>

__device__ volatile BYTE chars[11] = "0123456789";

__device__ BYTE uint32_to_byte(uint32_t val) {
  return chars[val];
}

__device__ BYTE *nonce_bytes(uint32_t nonce, int size) {
  BYTE *str = (BYTE*) malloc(size);

  for(int i = size - 1; i >= 0; i--) {
    str[i] = uint32_to_byte(nonce % 10);
    nonce /= 10;
  }

  return str;
}

__device__ int nonce_size(uint32_t nonce) {
    int size = 1;

    while (nonce /= 10)
      size++;

    return size;
}

__device__ BYTE *join_input_and_nonce(BYTE *input, uint32_t nonce, int size) {
  int nsize = nonce_size(nonce);
  int total_size = size + nsize;

  BYTE *nbytes = nonce_bytes(nonce, nsize);
  BYTE *bytes = (BYTE *) malloc(total_size);

  int i;

  for(i = 0; i < size; i++) {
    bytes[i] = input[i];
  }

  for(int j = 0; j < nsize; j++, i++) {
    bytes[i] = nbytes[j];
  }

  free(nbytes);

  return bytes;
}

__device__ bool verify(int id, BYTE *input, uint32_t nonce, int size, int difficulty, int *found) {
  BYTE *first = (BYTE *) malloc(SHA256_BLOCK_SIZE);
  BYTE *hash = (BYTE *) malloc(SHA256_BLOCK_SIZE);
  BYTE *test = join_input_and_nonce(input, nonce, size);

  csha256(test, first, size + nonce_size(nonce));
  csha256(get_sha_string(first), hash, SHA256_STRING_SIZE);

  free(first);
  free(test);

  int aux = difficulty;
  int blocks = (difficulty + 1) / 2;

  for(int i = 0; i < blocks; i++) {
    BYTE cmp;

    if(aux == 1) {
      cmp = 0x0F;
    } else {
      cmp = 0x00;
    }

    if(hash[i] > cmp) {
      free(hash);
      return false;
    }

    aux -= 2;
  }

  *found = FOUND;

  free(hash);
  return true;
}

#endif   // MINE_CUH
