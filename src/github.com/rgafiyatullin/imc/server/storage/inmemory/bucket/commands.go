package bucket

const cmdSet = 1
const cmdGet = 2
const cmdExists = 3
const cmdDel = 4

const cmdLPushBack = 5
const cmdLPushFront = 6
const cmdLPopBack = 7
const cmdLPopFront = 8
const cmdLGetNth = 9

type Cmd interface {
	CmdId() int
}

func NewCmdGet(key string) Cmd {
	cmd := new(CmdGet)
	cmd.key = key
	return cmd
}

func NewCmdSet(key string, value []byte, expiry uint64) Cmd {
	cmd := new(CmdSet)
	cmd.key = key
	cmd.value = value
	cmd.expiry = expiry
	return cmd
}

func NewCmdDel(key string) Cmd {
	cmd := new(CmdDel)
	cmd.key = key
	return cmd
}

func NewCmdLPopBack(key string) Cmd {
	cmd := new(CmdLPopBack)
	cmd.key = key
	return cmd
}

func NewCmdLPopFront(key string) Cmd {
	cmd := new(CmdLPopFront)
	cmd.key = key
	return cmd
}

func NewCmdLGetNth(key string, idx int) Cmd {
	cmd := new(CmdLGetNth)
	cmd.key = key
	cmd.idx = idx
	return cmd
}

func NewCmdLPushBack(key string, value []byte) Cmd {
	cmd := new(CmdLPushBack)
	cmd.key = key
	cmd.value = value
	return cmd
}

func NewCmdLPushFront(key string, value []byte) Cmd {
	cmd := new(CmdLPushFront)
	cmd.key = key
	cmd.value = value
	return cmd
}

type CmdSet struct {
	key    string
	value  []byte
	expiry uint64
}

func (this *CmdSet) CmdId() int {
	return cmdSet
}

type CmdGet struct {
	key string
}

func (this *CmdGet) CmdId() int {
	return cmdGet
}

type CmdExists struct {
	key string
}

func (this *CmdExists) CmdId() int {
	return cmdExists
}

type CmdDel struct {
	key string
}

func (this *CmdDel) CmdId() int {
	return cmdDel
}

type CmdLPushBack struct {
	key   string
	value []byte
}

func (this *CmdLPushBack) CmdId() int {
	return cmdLPushBack
}

type CmdLPushFront struct {
	key   string
	value []byte
}

func (this *CmdLPushFront) CmdId() int {
	return cmdLPushFront
}

type CmdLPopBack struct {
	key string
}

func (this *CmdLPopBack) CmdId() int {
	return cmdLPopBack
}

type CmdLPopFront struct {
	key string
}

func (this *CmdLPopFront) CmdId() int {
	return cmdLPopFront
}

type CmdLGetNth struct {
	key string
	idx int
}

func (this *CmdLGetNth) CmdId() int {
	return cmdLGetNth
}
