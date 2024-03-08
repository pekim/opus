#pragma once

#include <opusfile.h>

OpusFileCallbacks *create_file_callbacks();
void free_file_callbacks(OpusFileCallbacks *callbacks);
