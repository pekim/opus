package opus

import _ "embed"

// SampleStream is the data for a valid opus stream,
// that can be used for testing or demonstration purposes.
//
//go:embed sample1.opus
var SampleStream []byte
