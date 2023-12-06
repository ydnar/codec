# codec

[![build status](https://img.shields.io/github/actions/workflow/status/ydnar/codec/go.yaml?branch=main)](https://github.com/ydnar/codec/actions)
[![pkg.go.dev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/ydnar/codec)

A [Go](https://go.dev/) experiment for encoding and decoding serialization formats like JSON or XML without [reflect](https://pkg.go.dev/reflect).

Extracted from [wasm-tools-go](https://github.com/ydnar/wasm-tools-go).

## TODO

- [ ] Remove implicit dependencies on `reflect`
  - [ ] Replace `encoding/json` with [`jsontext`](https://pkg.go.dev/github.com/go-json-experiment/json/jsontext)

## License

This project is licensed under the Apache 2.0 license with the LLVM exception. See [LICENSE](LICENSE) for more details.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in this project by you, as defined in the Apache-2.0 license, shall be licensed as above, without any additional terms or conditions.
