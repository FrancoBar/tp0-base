package common

import (
	"io"
	"encoding/binary"
)

func appendLenght(data []byte) []byte{
	length := make([]byte, 4)
    binary.BigEndian.PutUint32(length, uint32(len(data)))
	return append(length, data...)
}

func SerializeStr(s string) []byte{
	data := []byte(s)
	return appendLenght(data)
}

func SerializeUint32(u uint32) []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, u)
	return data
}

func SerializeUint64(u uint64) []byte{
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, u)
	return data
}

func SerializeBool(b bool) []byte{
	data := make([]byte, 1)
	if b {
		data[0] = 1
	}else{
		data[0] = 0
	}
	return data
}

func recvSized(reader io.Reader, size uint32) ([]byte, error){
	buf := make([]byte, size)
	nread, err := io.ReadFull(reader, buf)
	if err != nil || uint32(nread) < size{
		return nil, io.ErrUnexpectedEOF
	}
	return buf, nil
}

func DeserializeUint32(reader io.Reader) (uint32, error){
	u, err := recvSized(reader, 4)
	if err != nil{
		return 0, err
	}
	return binary.BigEndian.Uint32(u), nil
}

func DeserializeBool(reader io.Reader) (bool, error){
	u, err := recvSized(reader, 1)
	if err != nil{
		return false, err
	}
	return  (u[0] != 0), nil
}
