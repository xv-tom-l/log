// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logs

import (
	"encoding/json"

	"github.com/go-xorm/xorm"
)

type Log struct {
	Id    int64
	Level int
	Msg   string `xorm:"TEXT"`
}

// DatabaseWriter implements LoggerInterface and is used to log into database.
type DatabaseWriter struct {
	Driver string `json:"driver"`
	Conn   string `json:"conn"`
	Level  int    `json:"level"`
	x      *xorm.Engine
}

func NewDatabase() LoggerInterface {
	return &DatabaseWriter{Level: LevelTrace}
}

// init database writer with json config.
// config like:
//	{
//		"driver": "mysql"
//		"conn":"root:root@tcp(127.0.0.1:3306)/gogs?charset=utf8",
//		"level": 0
//	}
// connection string is based on xorm.
func (d *DatabaseWriter) Init(jsonconfig string) (err error) {
	if err = json.Unmarshal([]byte(jsonconfig), d); err != nil {
		return err
	}
	d.x, err = xorm.NewEngine(d.Driver, d.Conn)
	if err != nil {
		return err
	}
	return d.x.Sync(new(Log))
}

// write message in database writer.
func (d *DatabaseWriter) WriteMsg(msg string, level int) error {
	if level < d.Level {
		return nil
	}

	_, err := d.x.Insert(&Log{Level: level, Msg: msg})
	return err
}

// implementing method. empty.
func (d *DatabaseWriter) Destroy() {

}

// implementing method.
func (d *DatabaseWriter) Flush() {

}

func init() {
	Register("database", NewDatabase)
}
