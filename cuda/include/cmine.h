#ifndef CMINE_H
#define CMINE_H

#define NOT_FOUND -1
#define FOUND 1
#define BLOCKS 512
#define THREADS 512
#define TOTAL BLOCKS * THREADS

extern "C" {
  uint32_t cmine(const char *str, int difficulty, uint32_t max_nonce);
}

#endif