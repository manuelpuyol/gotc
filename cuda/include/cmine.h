#ifndef CMINE_H
#define CMINE_H

#define NOT_FOUND -1
#define FOUND 1
#define BLOCKS 512
#define THREADS 512
#define TOTAL BLOCKS * THREADS
#define MAX_NONCE UINT32_MAX
#define BUCKET (MAX_NONCE / uint32_t(TOTAL)) + 1

extern "C" {
  uint32_t cmine(const char *str, int difficulty);
}

#endif