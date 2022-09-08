package common

const (
	AskWinner uint32 = iota
	AskAmount
)

type PersonRecord struct {
	FirstName string
	LastName string
	Document uint64
	Birthdate string
}

func SerializePersonRecord(p PersonRecord) []byte{
	first_name :=	SerializeStr(p.FirstName)
	last_name :=	SerializeStr(p.LastName)
	document :=		SerializeUint64(p.Document)
	birthdate :=	SerializeStr(p.Birthdate)
	data := append(first_name, last_name...)
	data = append(data, document...)
	data = append(data, birthdate...)
	return data
}

func SerializePersonRecordArray(vec []PersonRecord) []byte{
	data := SerializeUint32(uint32(len(vec)))
	for _, elem := range vec {
		data = append(data, SerializePersonRecord(elem)...)
	}
	return data
}
