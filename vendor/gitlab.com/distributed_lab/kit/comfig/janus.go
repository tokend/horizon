package comfig

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go-janus"
)

type Januser interface {
	Janus() *janus.Janus
}

func NewJanuser(getter kv.Getter) Januser {
	return &januser{
		getter: getter,
	}
}

type januser struct {
	getter kv.Getter
	once   Once
}

func (j *januser) Janus() *janus.Janus {
	return j.once.Do(func() interface{} {
		var config struct {
			URL     string `fig:"url,required"`
			Target  string `fig:"target,required"`
			Surname string `fig:"surname,required"`
		}
		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(j.getter, "janus")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out janus"))
		}

		js := janus.NewJanus(config.URL, config.Target, config.Surname)
		return js
	}).(*janus.Janus)
}
