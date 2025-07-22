package flow

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestNextif_match_all(t *testing.T) {
	target := []Node{
		{},
		{},
		{},
	}
	next := Next(func(target []Node) error {
		target[0].UUID = option.Some(MustUUID("14725bfb-6562-4f14-8841-df255fa9082a"))
		target[1].UUID = option.Some(MustUUID("24725bfb-6562-4f14-8841-df255fa9082a"))
		target[2].UUID = option.Some(MustUUID("34725bfb-6562-4f14-8841-df255fa9082a"))
		return nil
	})
	predicate := func(n Node) bool {
		return true
	}
	err := nextIf(target, next, predicate)
	require.NoError(t, err)
	require.Equal(t, "14725bfb-6562-4f14-8841-df255fa9082a", target[0].UUID.Get().String())
	require.Equal(t, "24725bfb-6562-4f14-8841-df255fa9082a", target[1].UUID.Get().String())
	require.Equal(t, "34725bfb-6562-4f14-8841-df255fa9082a", target[2].UUID.Get().String())
}
func TestNextif_match_none(t *testing.T) {
	target := []Node{
		{UUID: option.Some(MustUUID("14725bfb-6562-4f14-8841-df255fa9082a"))},
		{UUID: option.Some(MustUUID("24725bfb-6562-4f14-8841-df255fa9082a"))},
		{UUID: option.Some(MustUUID("34725bfb-6562-4f14-8841-df255fa9082a"))},
	}
	next := Next(func(target []Node) error {
		target[0].UUID = option.Some(MustUUID("64725bfb-6562-4f14-8841-df255fa9082a"))
		target[1].UUID = option.Some(MustUUID("64725bfb-6562-4f14-8841-df255fa9082a"))
		target[2].UUID = option.Some(MustUUID("64725bfb-6562-4f14-8841-df255fa9082a"))
		return nil
	})
	predicate := func(n Node) bool {
		return false
	}
	err := nextIf(target, next, predicate)
	require.NoError(t, err)
	require.Equal(t, "14725bfb-6562-4f14-8841-df255fa9082a", target[0].UUID.Get().String())
	require.Equal(t, "24725bfb-6562-4f14-8841-df255fa9082a", target[1].UUID.Get().String())
	require.Equal(t, "34725bfb-6562-4f14-8841-df255fa9082a", target[2].UUID.Get().String())
}
func TestNextif_match_first(t *testing.T) {
	target := []Node{
		{UUID: option.Some(MustUUID("64725bfb-6562-4f14-8841-df255fa9082a"))},
		{},
		{},
	}
	next := Next(func(target []Node) error {
		target[0].UUID = option.Some(MustUUID("14725bfb-6562-4f14-8841-df255fa9082a"))
		return nil
	})
	predicate := func(n Node) bool {
		return n.UUID.GetOrZero() == MustUUID("64725bfb-6562-4f14-8841-df255fa9082a")
	}
	err := nextIf(target, next, predicate)
	require.NoError(t, err)
	require.Equal(t, "14725bfb-6562-4f14-8841-df255fa9082a", target[0].UUID.Get().String())
	require.Zero(t, target[1].UUID)
	require.Zero(t, target[2].UUID)
}
func TestNextif_match_middle(t *testing.T) {
	target := []Node{
		{},
		{UUID: option.Some(MustUUID("64725bfb-6562-4f14-8841-df255fa9082a"))},
		{},
	}
	next := Next(func(target []Node) error {
		target[0].UUID = option.Some(MustUUID("14725bfb-6562-4f14-8841-df255fa9082a"))
		return nil
	})
	predicate := func(n Node) bool {
		return n.UUID.GetOrZero() == MustUUID("64725bfb-6562-4f14-8841-df255fa9082a")
	}
	err := nextIf(target, next, predicate)
	require.NoError(t, err)
	require.Zero(t, target[0].UUID)
	require.Equal(t, "14725bfb-6562-4f14-8841-df255fa9082a", target[1].UUID.Get().String())
	require.Zero(t, target[2].UUID)
}
func TestNextif_match_last(t *testing.T) {
	target := []Node{
		{},
		{},
		{UUID: option.Some(MustUUID("64725bfb-6562-4f14-8841-df255fa9082a"))},
	}
	next := Next(func(target []Node) error {
		target[0].UUID = option.Some(MustUUID("14725bfb-6562-4f14-8841-df255fa9082a"))
		return nil
	})
	predicate := func(n Node) bool {
		return n.UUID.GetOrZero() == MustUUID("64725bfb-6562-4f14-8841-df255fa9082a")
	}
	err := nextIf(target, next, predicate)
	require.NoError(t, err)
	require.Zero(t, target[0].UUID)
	require.Zero(t, target[1].UUID)
	require.Equal(t, "14725bfb-6562-4f14-8841-df255fa9082a", target[2].UUID.Get().String())
}
