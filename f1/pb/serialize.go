package pb

import (
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var order = binary.LittleEndian

// write msg size, then message
func Write(w io.Writer, msg protoreflect.ProtoMessage) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}

	size := uint32(len(b))
	sizeb := make([]byte, 4)
	order.PutUint32(sizeb, size)

	if _, err := w.Write(sizeb); err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

// read msg size, then message
func Read(r io.Reader, msg protoreflect.ProtoMessage) error {
	sizeb := make([]byte, 4)
	if _, err := io.ReadFull(r, sizeb); err != nil {
		return err
	}

	size := binary.LittleEndian.Uint32(sizeb)

	b := make([]byte, size)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}

	if err := proto.Unmarshal(b, msg); err != nil {
		return err
	}

	return nil
}
