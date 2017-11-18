## logs [![wercker status](https://app.wercker.com/status/ed8de801ba4452aac5109cdd13ab55aa/s/ "wercker status")](https://app.wercker.com/project/bykey/ed8de801ba4452aac5109cdd13ab55aa) [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/gogits/logs)

logs is a Go logs manager. It can use many logs adapters. The repo is inspired by `database/sql` .

**This is a fork of [beego/logs](https://github.com/astaxie/beego/tree/master/logs).**


## How to install?

	go get github.com/gogits/logs


## What adapters are supported?

As of now this logs supports `console`, `file`, `smtp`, `conn` and `database`.


## How to use it?

First you must import it

	import (
		"github.com/gogits/logs"
	)

Then init a Log (example with console adapter)

	log := NewLogger(10000)
	log.SetLogger("console", "")	

> the first params stand for how many channel

Use it like this:	
	
	log.Trace("trace")
	log.Info("info")
	log.Warn("warning")
	log.Debug("debug")
	log.Critical("critical")


## File adapter

Configure file adapter like this:

	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test.log"}`)


## Conn adapter

Configure like this:

	log := NewLogger(1000)
	log.SetLogger("conn", `{"net":"tcp","addr":":7020"}`)
	log.Info("info")


## Smtp adapter

Configure like this:

	log := NewLogger(10000)
	log.SetLogger("smtp", `{"username":"beegotest@gmail.com","password":"xxxxxxxx","host":"smtp.gmail.com:587","sendTos":["xiemengjun@gmail.com"]}`)
	log.Critical("sendmail critical")
	time.Sleep(time.Second * 30)

## License

Gogs is under the MIT License. See the [LICENSE](https://github.com/gogits/logs/blob/master/LICENSE) file for the full license text.