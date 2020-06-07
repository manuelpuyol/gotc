CFLAGS = --ptxas-options=-v --compiler-options '-Icuda/include/ -fPIC'
CC = nvcc

all:
	$(CC) $(CFLAGS) -o libgpu.so --shared cuda/src/mine.cu
	go build

clean:
	rm libgpu.so
