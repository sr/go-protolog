package protolog_benchmark_marshal

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"go.pedge.io/protolog"
	"go.pedge.io/protolog/testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkDelimitedMarshaller(b *testing.B) {
	benchmarkMarshaller(b, protolog.DelimitedMarshaller)
}

func BenchmarkDefaultTextMarshaller(b *testing.B) {
	benchmarkMarshaller(b, protolog.NewTextMarshaller(protolog.MarshallerOptions{}))
}

func benchmarkMarshaller(b *testing.B, marshaller protolog.Marshaller) {
	b.StopTimer()
	goEntry := getBenchGoEntry()
	_, err := marshaller.Marshal(goEntry)
	require.NoError(b, err)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = marshaller.Marshal(goEntry)
	}
}

func getBenchGoEntry() *protolog.GoEntry {
	foo := &protolog_testing.Foo{
		StringField: "one",
		Int32Field:  2,
	}
	bar := &protolog_testing.Bar{
		StringField: "one",
		Int32Field:  2,
	}
	baz := &protolog_testing.Baz{
		Bat: &protolog_testing.Baz_Bat{
			Ban: &protolog_testing.Baz_Bat_Ban{
				StringField: "one",
				Int32Field:  2,
			},
		},
	}
	goEntry := &protolog.GoEntry{
		ID:    "123",
		Level: protolog.Level_LEVEL_INFO,
		Time:  time.Now().UTC(),
		Contexts: []proto.Message{
			foo,
			bar,
		},
		Event: baz,
	}
	return goEntry
}
