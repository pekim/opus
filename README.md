# opus

[![PkgGoDev](https://pkg.go.dev/badge/github.com/pekim/opus)](https://pkg.go.dev/github.com/pekim/opus)

This Go package provides decoding of streams
that have been encoded with the
[opus coded](https://opus-codec.org/).
When provided with an opus stream the raw PCM data
can be extracted.

The xiph.org
[opus](https://github.com/xiph/opus),
[opusfile](https://github.com/xiph/opusfile), and
[ogg](https://github.com/xiph/ogg)
C libraries are used to do the decoding.

# decoding a stream

```go
readSeeker, err := os.Open("an-opus-file.opus")
if err != nil { ... }

decoder, err := opus.NewDecoder(readSeeker)
if err != nil { ... }

var pcm = make([]int16, 2_000)
for {
    n, err := decoder.Read(pcm)
    fmt.Printf("read %d samples (per channel)\n", n)
    samples := pcm[:n]
    // Do something with the samples...

    if err == io.EOF {
        break
    }
}

decoder.Destroy()
readSeeker.Close()
```

## building

Because of the use of cgo in the library,
the Go toolchain will need to be able to find a
C compiler (typically gcc or clang) and linker on the path.

When building an application that uses the opus package
(with `go build` or `go run`)
the first build may take several seconds.
This is because of the use of a C compiler to compile the
`opus`, `opusfile`, and `ogg` libraries.
Subsequent builds will be much quicker as the built package
will be cached by the Go toolchain.

## examples

There are a couple of examples in the `example` directory.

- `beep`
  - Uses the [beep](https://github.com/gopxl/beep) package to play an opus file.
  - `go run example\beep\*.go [your-file.opus]`
    - Without an argument, plays an embedded sample opus file.
    - With an argument that is a path to an opus file, will play the file.
- `simple`
  - Uses various methods of `opus.Decoder` to get and print various information
    about an opus stream.
  - `go run example\simple\*.go [your-file.opus]`
    - Without an argument, uses an embedded sample opus file.
    - With an argument that is a path to an opus file, uses the file.

## possible alternatives

### pion/opus

https://github.com/pion/opus is a pure Go implementation of the opus codec.

- It currently only supports
  the [SILK](https://en.wikipedia.org/wiki/SILK) codec,
  not the [CELT](https://celt-codec.org/) codec.
  See [issue #7](https://github.com/pion/opus/issues/7).

### hraban/opus

https://github.com/hraban/opus provides C bindings to the xiph.org C libraries.

- It is more comprehensive than this library,
  and supports encoding as well as decoding of opus streams.
- It requires libopus and libopusfile development packages to be installed on the system.

## development

### pre-commit hook

- install `goimports` if not already installed
  - https://pkg.go.dev/golang.org/x/tools/cmd/goimports
- install `golangci-lint` if not already installed
  - https://golangci-lint.run/usage/install/#local-installation
- install the `pre-commit` application if not already installed
  - https://pre-commit.com/index.html#install
- install pre-commit hook in this repo's workspace
  - `pre-commit install`
