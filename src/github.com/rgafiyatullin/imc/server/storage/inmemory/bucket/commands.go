package bucket

const cmdSet = 1
const cmdGet = 2
const cmdDel = 4
const cmdKeys = 17

const cmdLPushBack = 5
const cmdLPushFront = 6
const cmdLPopBack = 7
const cmdLPopFront = 8
const cmdLLen = 3
const cmdLGetNth = 9

const cmdExpire = 10
const cmdTTL = 11

const cmdHSet = 12
const cmdHGet = 13
const cmdHDel = 14
const cmdHKeys = 15
const cmdHGetAll = 16

type Cmd interface {
	CmdId() int
}

func NewCmdGet(key string) Cmd {
	cmd := new(CmdGet)
	cmd.key = key
	return cmd
}

func NewCmdSet(key string, value []byte) Cmd {
	cmd := new(CmdSet)
	cmd.key = key
	cmd.value = value
	return cmd
}

func NewCmdDel(key string) Cmd {
	cmd := new(CmdDel)
	cmd.key = key
	return cmd
}

func NewCmdKeys(pattern string) Cmd {
	cmd := new(CmdKeys)
	cmd.pattern = pattern
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

func NewCmdLLen(key string) Cmd {
	cmd := new(CmdLLen)
	cmd.key = key
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

func NewCmdExpire(key string, expiry int64) Cmd {
	cmd := new(CmdExpire)
	cmd.key = key
	cmd.expiry = expiry
	return cmd
}

func NewCmdTTL(key string, useSeconds bool) Cmd {
	cmd := new(CmdTTL)
	cmd.key = key
	cmd.useSeconds = useSeconds
	return cmd
}

func NewCmdHGetAll(key string) Cmd {
	cmd := new(CmdHGetAll)
	cmd.key = key
	return cmd
}

func NewCmdHKeys(key string) Cmd {
	cmd := new(CmdHKeys)
	cmd.key = key
	return cmd
}

func NewCmdHDel(key string, hkey string) Cmd {
	cmd := new(CmdHDel)
	cmd.key = key
	cmd.hkey = hkey
	return cmd
}

func NewCmdHGet(key string, hkey string) Cmd {
	cmd := new(CmdHGet)
	cmd.key = key
	cmd.hkey = hkey
	return cmd
}

func NewCmdHSet(key string, hkey string, hvalue []byte) Cmd {
	cmd := new(CmdHSet)
	cmd.key = key
	cmd.hkey = hkey
	cmd.hvalue = hvalue
	return cmd
}

type CmdSet struct {
	key   string
	value []byte
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

type CmdDel struct {
	key string
}

func (this *CmdDel) CmdId() int {
	return cmdDel
}

type CmdKeys struct {
	pattern string
}

func (this *CmdKeys) CmdId() int {
	return cmdKeys
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

type CmdLLen struct {
	key string
}

func (this *CmdLLen) CmdId() int {
	return cmdLLen
}

type CmdExpire struct {
	key    string
	expiry int64
}

func (this *CmdExpire) CmdId() int {
	return cmdExpire
}

type CmdTTL struct {
	key        string
	useSeconds bool
}

func (this *CmdTTL) CmdId() int {
	return cmdTTL
}

type CmdHSet struct {
	key    string
	hkey   string
	hvalue []byte
}

func (this *CmdHSet) CmdId() int {
	return cmdHSet
}

type CmdHGet struct {
	key  string
	hkey string
}

func (this *CmdHGet) CmdId() int {
	return cmdHGet
}

type CmdHDel struct {
	key  string
	hkey string
}

func (this *CmdHDel) CmdId() int {
	return cmdHDel
}

type CmdHKeys struct {
	key string
}

func (this *CmdHKeys) CmdId() int {
	return cmdHKeys
}

type CmdHGetAll struct {
	key string
}

func (this *CmdHGetAll) CmdId() int {
	return cmdHGetAll
}
