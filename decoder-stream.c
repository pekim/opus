#include "_cgo_export.h"
#include <opusfile.h>
#include <stdlib.h>

static int file_read(void *_stream, unsigned char *_ptr, int _nbytes)
{
  return goFileRead(_stream, _ptr, _nbytes);
}

static int file_seek(void *_stream, opus_int64 _offset, int _whence)
{
  return goFileSeek(_stream, _offset, _whence);
}

static opus_int64 file_tell(void *_stream)
{
  return goFileTell(_stream);
}

static int file_close(void *_stream)
{
  return goFileClose(_stream);
}

OpusFileCallbacks *create_file_callbacks()
{
  OpusFileCallbacks *callbacks = malloc(sizeof(OpusFileCallbacks));
  callbacks->read = file_read;
  callbacks->seek = file_seek;
  callbacks->tell = file_tell;
  callbacks->close = file_close;

  return callbacks;
}

void free_file_callbacks(OpusFileCallbacks *callbacks)
{
  free(callbacks);
}
