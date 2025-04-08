package flow

import (
	"log/slog"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type UUID [16]byte

func NewUUID() UUID {
	var u = UUID(uuid.Must(uuid.NewRandom()))
	return u
}
func MustUUID(s string) UUID {
	var u, err = ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}
func ParseUUID(s string) (u UUID, err error) {
	var x uuid.UUID
	if x, err = uuid.Parse(s); err != nil {
		return u, err
	}
	u = UUID(x)
	return u, nil
}
func (it UUID) GoString() string {
	return "\"" + it.String() + "\""
}
func (it UUID) LogValue() slog.Value {
	return slog.StringValue(it.String())
}
func (it UUID) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(it.String())
}
func (it UUID) String() string {
	return uuid.UUID(it).String()
}
