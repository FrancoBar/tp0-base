from ctypes import c_uint32
from asyncio import IncompleteReadError
from .utils import *
import socket



INTENTION_ASK_WINNER = 0
INTENTION_ASK_AMOUNT = 1

class InvalidIntentionError(Exception):
    pass


def _append_length(data):
	return len(data).to_bytes(4, 'big') + data

def serialize_uint32(u):
	return _append_length(u.to_bytes(4, 'big'))

def send(socket, msg):
	"""
	Serializes and sends uint32s through the provided socket
	"""
	m = serialize_uint32(msg)
	socket.sendall(m)


def _recv_n_bytes(socket, num_bytes):
	"""
	Receives exactly 'num_bytes' bytes through the provided socket.
	If no bytes are read from the socket IncompleteReadError is raised
	Source: https://stackoverflow.com/questions/55825905/how-can-i-reliably-read-exactly-n-bytes-from-a-tcp-socket
	"""
	buf = bytearray(num_bytes)
	pos = 0
	while pos < num_bytes:
		n = socket.recv_into(memoryview(buf)[pos:])
		if n == 0:
			raise IncompleteReadError(bytes(buf[:pos]), num_bytes)
		pos += n
	return bytes(buf)

def _recv_sized(socket):
	size = int.from_bytes(_recv_n_bytes(socket, 4), byteorder='big', signed=False)
	return _recv_n_bytes(socket, size)

def recv_intention(socket):
	return int.from_bytes(_recv_n_bytes(socket, 4), byteorder='big', signed=False)

def recv_str(socket):
	return _recv_sized(socket).decode('utf-8')

def recv_unsigned_number(socket):
	return int.from_bytes(_recv_sized(socket), byteorder='big', signed=False)

def recv_person_record(socket):
	return Contestant(
		recv_str(socket),
		recv_str(socket),
		recv_unsigned_number(socket),
		recv_str(socket)
	)

def recv_vector(socket, recv_type):
	vec_len = recv_unsigned_number(socket)
	vec = []
	while vec_len > 0:
		vec.append(recv_type(socket))
		vec_len -= 1
	return vec

def recv(socket):
	"""
	Listens to messages through the provided socket. Handles the intention field, even if it isn't currently used.
	"""
	intention = recv_intention(socket)
	if intention != INTENTION_ASK_WINNER:
		raise InvalidIntentionError
	return recv_vector(socket, recv_person_record)
