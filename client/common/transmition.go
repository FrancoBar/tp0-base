package common

import (
	"encoding/binary"
	"fmt"
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

func _serialize_uint64(u uint64) []byte{
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, u)
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
	first_name := _serialize_str(p.FirstName)
	last_name := _serialize_str(p.LastName)
	document :=	_serialize_uint64(p.Document)
	birthdate := _serialize_str(p.Birthdate)
	data := append(first_name, last_name...)
	data = append(data, document...)
	data = append(data, birthdate...)
	return data
}

func _serialize_person_records(ps []PersonRecord) []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(len(ps)))
	for _, p := range ps {
		data = append(data, _serialize_person_record(p)...)
	}
	return data
}

func _serialize_empty() []byte{
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 0)
	return data
}

func Serialize(x interface{}) ([]byte, error) {
    switch y := x.(type) {
        case string:
             return _serialize_str(y), nil
        case uint32:
             return _serialize_uint32(y), nil
        case uint64:
             return _serialize_uint64(y), nil
        case Intention:
        	 return _serialize_intention(y), nil
        case PersonRecord:
        	 return _serialize_person_record(y), nil
        case []PersonRecord:
        	 return _serialize_person_records(y), nil
        default:
        	return nil, fmt.Errorf("Error: %T type is not currently supported", x)
    }
}

// Send / Recv
func _serialize_and_write(socket net.Conn, msg interface{}) error{
	m, err := Serialize(msg)
	if err != nil {
		return err
	}

	nwritten_acum := len(m)
	for ;nwritten_acum > 0; {
		nwritten, err := socket.Write(m)
		if err != nil {
			return err
		}
		nwritten_acum -= nwritten
    	
	}
	return nil
}

func Send(socket net.Conn, intention Intention, msg interface{}) error{
	if intention >= Size{
		return fmt.Errorf("Error: Invalid intention", intention)
	}
	err := _serialize_and_write(socket, intention)
	if err != nil {
		return err
	}
	return _serialize_and_write(socket, msg)
}


func _recv_sized(reader io.Reader) ([]byte, error){
	buf := make([]byte, 4)
	nread, err := io.ReadFull(reader, buf)
	if err != nil || nread < 4{
		return nil, io.ErrUnexpectedEOF
	}

	size := binary.BigEndian.Uint32(buf)
	buf = make([]byte, size)
	nread, err = io.ReadFull(reader, buf)
	if err != nil || uint32(nread) < size{
		return nil, io.ErrUnexpectedEOF
	}
	return buf, nil
}

func RecvUint32(reader io.Reader) (uint32, error){
	u, err := _recv_sized(reader)
	if err != nil{
		return 0, err
	}
	return binary.BigEndian.Uint32(u), nil
}

func RecvBool(reader io.Reader) (bool, error){
	u, err := RecvUint32(reader)
	if err != nil{
		return false, err
	}
	return  (u != 0), nil
}
