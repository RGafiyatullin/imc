package actor

import "github.com/rgafiyatullin/imc/server/actor/logging"

type Ctx interface {
	Path() string
	Log() logging.Ctx
	NewChild(name string) Ctx
}

type impl struct {
	path_ string
	log_  logging.Ctx
}

func (this *impl) Log() logging.Ctx {
	return this.log_
}
func (this *impl) Path() string {
	return this.path_
}

func (this *impl) NewChild(name string) Ctx {
	child := new(impl)
	childName := this.path_ + "/" + name
	child.path_ = childName
	child.log_ = this.log_.CloneWithName(childName)
	return child
}

func NewCtx() Ctx {
	ctx := new(impl)
	ctx.log_ = logging.New()
	ctx.path_ = ""
	return ctx
}
