#include<mine.cuh>

__global__ void mine(BYTE *in, int *found, uint32_t *nonce, int size, int difficulty) {
  int id = (blockIdx.x * blockDim.x) + threadIdx.x;

  uint32_t test = uint32_t(id) * BUCKET;
  uint32_t end = id == TOTAL - 1
    ? MAX_NONCE
    : uint32_t(id + 1) * BUCKET;

  while(test < end && *found != FOUND) {
    if(verify(id, in, test, size, difficulty, found)) {
      *nonce = test;
    }

    test++;
  }
}

extern "C" {
  uint32_t cmine(const char *str, int difficulty) {
    // host
    BYTE *buff = (BYTE *) str;
    int size = strlen(str);
    int res = NOT_FOUND;
    uint32_t n;

    // device
    BYTE *in;
    int *found;
    uint32_t *nonce;

    cudaMalloc((void **)&in, size);
    cudaMalloc((void **)&found, sizeof(int));
    cudaMalloc((void **)&nonce, sizeof(uint32_t));

    cudaMemcpy(in, buff, size * sizeof(BYTE), cudaMemcpyHostToDevice);
    cudaMemcpy(found, &res, sizeof(int), cudaMemcpyHostToDevice);

    pre_sha256();
    mine<<< BLOCKS, THREADS >>>(in, found, nonce, size, difficulty);

    cudaDeviceSynchronize();

    cudaMemcpy(&res, found, sizeof(int), cudaMemcpyDeviceToHost);
    cudaMemcpy(&n, nonce, sizeof(uint32_t), cudaMemcpyDeviceToHost);

    cudaFree(in);
    cudaFree(found);
    cudaFree(nonce);

    return n;
  }
}
