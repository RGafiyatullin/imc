package server

import (
	"container/list"
	"github.com/rgafiyatullin/imc/protocol/resp/types/request"
	"fmt"
)

type Context interface {
	AddChunk(chunk []byte)

	HasRequest() bool
	FetchRequest() request.Request
}
type context struct {
	trailingCR  bool
	rawRequests *list.List // list.List<[]byte>
	chunks      *list.List // list.List<[]byte>
}

func New() Context {
	ctx := new(context)
	ctx.trailingCR = false
	ctx.chunks = list.New()
	ctx.rawRequests = list.New()
	return ctx
}

func (this *context) AddChunk(chunk []byte) {
	if len(chunk) == 0 {
		fmt.Println("!!! zero-lengthed chunk")
		return
	}

	if this.trailingCR && chunk[0] == '\n' {
		fmt.Println("!!! trailingCR and starts with LF")
		lastChunkElt := this.chunks.Back()
		withLF := lastChunkElt.Value.([]byte)
		withNoLF := withLF[:len(withLF)-1]
		lastChunkElt.Value = withNoLF // is this okay?

		this.finalizeRawRequest()

		nextChunk := chunk[1:]
		this.AddChunk(nextChunk)
	} else {
		fmt.Println("!!! general case")
		head, tail, ok := splitWithCRLF(chunk)
		if ok {
			this.chunks.PushBack(head)
			this.finalizeRawRequest()
			this.AddChunk(tail)
		} else {
			this.chunks.PushBack(chunk)
			this.trailingCR = chunk[len(chunk)-1] == '\r'
		}
	}
}

func (this *context) HasRequest() bool {
	return this.rawRequests.Len() > 0
}

func (this *context) FetchRequest() request.Request {
	if !this.HasRequest() {
		return nil
	}

	elt := this.rawRequests.Front()
	this.rawRequests.Remove(elt)
	return request.FromBytes(elt.Value.([]byte))
}

func (this *context) finalizeRawRequest() {
	rawReqSize := 0
	for element := this.chunks.Front(); element != nil; element = element.Next() {
		rawReqSize += len(element.Value.([]byte))
	}
	copyToPos := 0
	rawReq := make([]byte, rawReqSize)
	for element := this.chunks.Front(); element != nil; element = element.Next() {
		chunk := element.Value.([]byte)
		copyTo := rawReq[copyToPos : copyToPos+len(chunk)]
		copy(copyTo, chunk)
	}
	this.rawRequests.PushBack(rawReq)
	this.chunks = list.New()

	fmt.Printf("!!! finalizeRawRequest %v\n", rawReq)
}

func splitWithCRLF(chunk []byte) ([]byte, []byte, bool) {
	for i := 0; i < len(chunk)-1; i++ {
		if chunk[i] == '\r' && chunk[i+1] == '\n' {
			return chunk[:i], chunk[i+2:], true
		}
	}
	return nil, nil, false
}
