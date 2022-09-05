from ctypes import c_uint32
import socket
from asyncio import IncompleteReadError

# const (
# 	AskWinner Intention = iota
# 	AskAmount
# 	Size
# )

# Serialization

def _append_lenght(data):
	return len(data).to_bytes(4, 'big') + data

def serialize_uint32(u):
	return _append_length(u.to_bytes(4, 'big'))

# Send / Recv
def Send(socket, msg):
	socket.serialize_uint32(msg)

def _recv_n_bytes(sock, num_bytes):
    buf = bytearray(num_bytes)
    pos = 0
    while pos < num_bytes:
        n = sock.recv_into(memoryview(buf)[pos:])
        if n == 0:
            raise IncompleteReadError(bytes(buf[:pos]), num_bytes)
        pos += n
    return bytes(buf)

def _recv_sized(socket):
	size = int.from_bytes(_recv_n_bytes(socket, 4), byteorder='big', signed=False)
	return _recv_n_bytes(socket, size)

def RecvIntention(socket):
	return int.from_bytes(_recv_n_bytes(socket, 4), byteorder='big', signed=False)

def RecvStr(socket):
	return str(_recv_sized(socket))

def RecvUint32(socket):
	return int.from_bytes(_recv_sized(socket), byteorder='big', signed=False)

def RecvPersonRecord(socket):
	personrecord = {
	'name' : RecvStr(socket),
	'surname' : RecvStr(socket),
	'dni' : RecvUint32(socket),
	'birthdate' : RecvStr(socket)
	}
	return personrecord


def Recv(socket):
	intention = RecvIntention(socket)
	return RecvPersonRecord(socket)
