package keep

import (
	"iter"
	"slices"

	"github.com/typomaker/flow"
)

type Node struct{}

func (it *Node) When(flow.When) When {
	return When{}
}
func (it *Node) Fill(v []flow.Node) (err error) {
	sortNode(v)

	return nil
}
func (it *Node) Save(v []flow.Node) (err error) {
	sortNode(v)

	return nil
}
func (it *Node) Drop(v []flow.Node) (err error) {
	sortNode(v)

	return nil
}
func (it *Node) Case(v []flow.Case) (err error) {

	return nil
}

type When struct{}

func (it When) Read() Read
func (it When) Drop() (err error)

type Read struct{}

func (it Read) UUID() iter.Seq[flow.UUID]
func (it Read) Full() iter.Seq[flow.Node]

func sortNode(v []flow.Node) {
	slices.SortFunc(v, func(a, b flow.Node) int {
		var au = a.UUID.Get()
		var bu = b.UUID.Get()
		for i, v := range au {
			if v < bu[i] {
				return -1
			}
			if v > bu[i] {
				return 1
			}
		}
		return 0
	})
}
