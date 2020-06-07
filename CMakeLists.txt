project(BTC++)

cmake_minimum_required (VERSION 3.0.2)

################################################################################

set(CMAKE_CMAKE_BUILD_TYPE Debug)

find_package(CUDA REQUIRED)

################################################################################

add_subdirectory(${CMAKE_SOURCE_DIR}/cuda)
