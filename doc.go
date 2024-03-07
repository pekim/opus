/*
Package opus provides decoding of opus files,
using the xiph.org C libraries
[opus], [opusfile], and [ogg].

Rather than use cgo to utilise system installed instances of the libraries
(that may or may not be installed on a given system),
the source for the libraries is included with this package.
The libraries are compiled, with the help of cgo, when the package is compiled.
This means that a C compiler (such as gcc or clang) and a linker need to be
available on the path for cgo to use when building the package.

It may take several seconds to build this package the first time,
until the Go tools cache the result of the build.

[opus]: https://github.com/xiph/opus
[opusfile]: https://github.com/xiph/opusfile
[ogg]: https://github.com/xiph/ogg
*/
package opus
