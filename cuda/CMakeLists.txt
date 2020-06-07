project(gpu LANGUAGES CUDA)

################################################################################

file(GLOB sources "${CMAKE_SOURCE_DIR}/cuda/src/*.cu")
include_directories("${CMAKE_SOURCE_DIR}/cuda/include/" "/usr/local/cuda/include/")

################################################################################

set(CMAKE_CUDA_STANDARD 14)
set(CMAKE_CUDA_STANDARD_REQUIRED ON)

add_library(gpu SHARED ${sources})
set_target_properties(gpu PROPERTIES CUDA_SEPARABLE_COMPILATION ON)
