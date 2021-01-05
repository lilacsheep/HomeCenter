package mallory

import (
	"errors"
	"math/rand"

	"github.com/gogf/gf/container/gmap"
)

var (
	ErrNoInstances = errors.New("no instance")
)

type Balance interface {
	DoBalance(instances *gmap.TreeMap, domain ...string) (interface{}, error)
}

type RandomBalance struct{}

func (self *RandomBalance) DoBalance(instances *gmap.TreeMap, domain ...string) (interface{}, error) {
	if instances.IsEmpty() {
		return nil, ErrNoInstances
	} else if instances.Size() == 1 {
		return instances.Values()[0], nil
	} else {
		index := rand.Intn(instances.Size())
		return instances.Values()[index], nil
	}
}

type RoundRobinBalance struct {
	curIndex int
}

func (self *RoundRobinBalance) DoBalance(instances *gmap.TreeMap, domain ...string) (interface{}, error) {
	if instances.IsEmpty() {
		return nil, ErrNoInstances
	} else if instances.Size() == 1 {
		return instances.Values()[0], nil
	} else {
		if self.curIndex >= instances.Size() {
			self.curIndex = 0
		}
		self.curIndex++
		return instances.Values()[self.curIndex], nil
	}
}

func NewRandomBalance() *RandomBalance {
	return &RandomBalance{}
}

func NewRoundRobinBalance() *RoundRobinBalance {
	return &RoundRobinBalance{
		curIndex: 0,
	}
}
