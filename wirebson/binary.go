// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wirebson

import (
	"encoding/binary"
	"fmt"
)

// BinarySubtype represents BSON Binary's subtype.
type BinarySubtype byte

const (
	// BinaryGeneric represents a BSON Binary generic subtype.
	BinaryGeneric = BinarySubtype(0x00) // generic

	// BinaryFunction represents a BSON Binary function subtype.
	BinaryFunction = BinarySubtype(0x01) // function

	// BinaryGenericOld represents a BSON Binary generic-old subtype.
	BinaryGenericOld = BinarySubtype(0x02) // generic-old

	// BinaryUUIDOld represents a BSON Binary UUID old subtype.
	BinaryUUIDOld = BinarySubtype(0x03) // uuid-old

	// BinaryUUID represents a BSON Binary UUID subtype.
	BinaryUUID = BinarySubtype(0x04) // uuid

	// BinaryMD5 represents a BSON Binary MD5 subtype.
	BinaryMD5 = BinarySubtype(0x05) // md5

	// BinaryEncrypted represents a BSON Binary encrypted subtype.
	BinaryEncrypted = BinarySubtype(0x06) // encrypted

	// BinaryUser represents a BSON Binary user-defined subtype.
	BinaryUser = BinarySubtype(0x80) // user
)

// Binary represents BSON scalar type binary.
type Binary struct {
	B       []byte
	Subtype BinarySubtype
}

// sizeBinary returns the size of the encoding of v [Binary] in bytes.
func sizeBinary(v Binary) int {
	return len(v.B) + 5
}

// encodeBinary encodes [Binary] value v into b.
//
// b must be at least len(v.B)+5 ([sizeBinary]) bytes long; otherwise, encodeBinary will panic.
// Only b[0:len(v.B)+5] bytes are modified.
func encodeBinary(b []byte, v Binary) {
	i := len(v.B)

	binary.LittleEndian.PutUint32(b, uint32(i))
	b[4] = byte(v.Subtype)
	copy(b[5:5+i], v.B)
}

// decodeBinary decodes [Binary] value from b.
//
// If there is not enough bytes, decodeBinary will return a wrapped [ErrDecodeShortInput].
// If the input is otherwise invalid, a wrapped [ErrDecodeInvalidInput] is returned.
func decodeBinary(b []byte) (Binary, error) {
	var res Binary

	if len(b) < 5 {
		return res, fmt.Errorf("DecodeBinary: expected at least 5 bytes, got %d: %w", len(b), ErrDecodeShortInput)
	}

	i := int(binary.LittleEndian.Uint32(b))
	if e := 5 + i; len(b) < e {
		return res, fmt.Errorf("DecodeBinary: expected at least %d bytes, got %d: %w", e, len(b), ErrDecodeShortInput)
	}

	res.Subtype = BinarySubtype(b[4])

	if i > 0 {
		res.B = make([]byte, i)
		copy(res.B, b[5:5+i])
	}

	return res, nil
}
