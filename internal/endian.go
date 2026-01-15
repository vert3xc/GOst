package internal

import "encoding/binary"


func BEToUint32(b []byte) uint32 {
    return binary.BigEndian.Uint32(b)
}

func LEToUint32(b []byte) uint32 {
    return binary.LittleEndian.Uint32(b)
}

func Uint32ToBE(dst []byte, v uint32) {
    binary.BigEndian.PutUint32(dst, v)
}

func Uint32ToLE(dst []byte, v uint32) {
    binary.LittleEndian.PutUint32(dst, v)
}

func BEToUint64(b []byte) uint64 {
    return binary.BigEndian.Uint64(b)
}

func Uint64ToBE(dst []byte, v uint64) {
    binary.BigEndian.PutUint64(dst, v)
}

func Uint64ToLE(dst []byte, v uint64) {
    binary.LittleEndian.PutUint64(dst, v)
}