package types

import (
	"fmt"
	"container/list"
	"net/textproto"
)

func NewArray(el *list.List) *BasicArr {
	a := new(BasicArr)
	ea := make([]BasicType, el.Len())
	i := 0
	for e := el.Front(); e != nil; e = e.Next() {
		ea[i] = e.Value.(BasicType)
		i++
	}
	a.elements = ea
	return a
}

func NewInt(v int64) *BasicInt {
	i := new(BasicInt)
	i.i = v
	return i
}

func NewErr(v string) *BasicErr {
	e := new(BasicErr)
	e.e = v
	return e
}

func NewStr(v string) *BasicStr {
	s := new(BasicStr)
	s.s = v
	return s
}

func NewBulkStr(v []byte) *BasicBulkStr {
	s := new(BasicBulkStr)
	s.s = v
	return s
}


type BasicType interface{
	ToString() string
	Write(to *textproto.Conn)
}

type BasicStr struct {
	s string
}
func (this *BasicStr) ToString() string {
	return fmt.Sprintf("S(\"%s\")", this.s)
}
func (this *BasicStr) Write(to *textproto.Conn) {
	to.Cmd("+%s", this.s)
}

type BasicErr struct {
	e string
}
func (this *BasicErr) ToString() string {
	return fmt.Sprintf("E(\"%s\")", this.e)
}
func (this *BasicErr) Write(to *textproto.Conn) {
	to.Cmd("-%s", this.e)
}

type BasicInt struct {
	i int64
}
func (this *BasicInt) ToString() string {
	return fmt.Sprintf("\"I(%d)\"", this.i)
}
func (this *BasicInt) Write(to *textproto.Conn) {
	to.Cmd(":%d", this.i)
}

type BasicBulkStr struct {
	s []byte
}
func (this *BasicBulkStr) ToString() string {
	return fmt.Sprintf("B(\"%s\")", this.s)
}
func (this *BasicBulkStr) Write(to *textproto.Conn) {
	to.Cmd("$%d", len(this.s))
	to.W.Write(this.s)
	to.Cmd("")
}

type BasicArr struct {
	elements []BasicType
}
func (this *BasicArr) ToString() string {
	acc := "A("
	for i := 0; i < len(this.elements); i++ {
		acc += this.elements[i].ToString() + ", "
	}
	acc += ")"
	return acc
}

func (this *BasicArr) Write(to *textproto.Conn) {
	to.Cmd("*%d", len(this.elements))
	for i := 0; i < len(this.elements); i++ {
		this.elements[i].Write(to)
	}
}


