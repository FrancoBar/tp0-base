UINT32_SIZE = 4
UINT64_SIZE = 8
BOOL_SIZE = 1

def serialize_uint32(u):
	return u.to_bytes(UINT32_SIZE, 'big')

def serialize_bool(u):
	return int(u).to_bytes(BOOL_SIZE, 'big')

def deserialize_unsigned_number(b):
	return int.from_bytes(b, byteorder='big', signed=False)

def deserialize_str(b):
	return b.decode('utf-8')