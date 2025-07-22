package flow

import (
	"context"
	"log/slog"
	"slices"

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
func (it UUID) In(u ...UUID) Handler {
	u = append(u, it)
	slices.SortFunc(u, UUID.Compare)
	var predicat = func(n Node) bool {
		if !n.UUID.IsSome() {
			return false
		}
		var ok bool
		_, ok = slices.BinarySearchFunc(u, n.UUID.Get(), UUID.Compare)
		return ok
	}
	return func(ctx context.Context, target []Node, next Next) (err error) {
		return nextIf(target, next, predicat)
	}
}
func (it UUID) Compare(t UUID) int {
	for i := range it {
		switch {
		case it[i] < t[i]:
			return -1
		case it[i] > t[i]:
			return 1
		}
	}
	return 0
}
