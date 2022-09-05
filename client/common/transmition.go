package common

import (
	"encoding/binary"
	"net"
	"io"
)

type Intention uint32

const (
	AskWinner Intention = iota
	AskAmount
	Size
)

// Serialization
func _append_lenght(data []byte) []byte{
	length := make([]byte, 4)
    binary.BigEndian.PutUint32(length, uint32(len(data)))
	return append(length, data...)
}

func _serialize_str(s string) []byte{
	data := []byte(s)
	return _append_lenght(data)
}

func _serialize_uint32(u uint32) []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, u)
	return _append_lenght(data)
}

func _serialize_bool(b bool) []byte{
	data := make([]byte, 1)
	if b {
		binary.BigEndian.PutUint32(data, 1)
	}else{
		binary.BigEndian.PutUint32(data, 0)
	}
	
	return _append_lenght(data)
}

func _serialize_intention(u Intention) []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(u))
	return data
}

func _serialize_person_record(p PersonRecord) []byte{
	name := _serialize_str(p.Name)
	surname := _serialize_str(p.Surname)
	dni := _serialize_uint32(p.Dni)
	birthdate := _serialize_str(p.Birthdate)
	data := append(name, surname...)
	data = append(data, dni...)
	data = append(data, birthdate...)
	return data
}

func _serialize_person_records(ps []PersonRecord) []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(len(ps)))
	for _, p := range ps {
		data = append(data, _serialize_person_record(p)...)
    // element is the element from someSlice for where we are
	}
	return data
}

func _serialize_empty() []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 0)
	return data
}

func Serialize(x interface{}) []byte {
    switch y := x.(type) {
        case string:
             return _serialize_str(y)
        case uint32:
             return _serialize_uint32(y)
        case Intention:
        	 return _serialize_intention(y)
        case PersonRecord:
        	 return _serialize_person_record(y)
        case []PersonRecord:
        	 return _serialize_person_records(y)
        default:
        	//TODO: Return error
        	return _serialize_empty()
    }
}

// Send / Recv
func Send(socket net.Conn, intention Intention, msg interface{}){
	if intention >= Size{
		//Err invalid intention
	}
	socket.Write(Serialize(intention))
	socket.Write(Serialize(msg))
}


func _recv_sized(reader io.Reader) []byte{
	buf := make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		//log.Fatal(err)
	}
	size := binary.BigEndian.Uint32(buf)
	buf = make([]byte, size)
	if _, err := io.ReadFull(reader, buf); err != nil {
		//log.Fatal(err)
	}
	return buf
}

func RecvUint32(reader io.Reader) uint32{
	return binary.BigEndian.Uint32(_recv_sized(reader))
}

func RecvStr(reader io.Reader) string{
	return string(_recv_sized(reader))
}

func RecvBool(reader io.Reader) bool{
	return (binary.BigEndian.Uint32(_recv_sized(reader)) != 0)
}

func RecvPersonRecord(reader io.Reader) PersonRecord{
	p := PersonRecord{
		Name: RecvStr(reader),
		Surname: RecvStr(reader),
		Dni: RecvUint32(reader),
		Birthdate: RecvStr(reader),
	}
	return p
}


func RecvPersonRecords(reader io.Reader) []PersonRecord{
	amount := RecvUint32(reader)
	var ps []PersonRecord
	for ; amount > 0; amount-- {
		ps = append(ps, RecvPersonRecord(reader))
	}
	return ps
}
