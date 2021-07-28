package tool

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

var DbEngine *Orm

type Orm struct {
	*xorm.Engine
}

func OrmEngine(cfg *Config) (*Orm, error) {
	database := cfg.Database
	conn := database.User + ":" + database.Password + "@tcp(" + database.Host + ":" + database.Port + ")/" + database.DbName + "?charset=" + database.Charset
	engine, err := xorm.NewEngine(database.Driver, conn)

	if err != nil {
		return nil, err
	}
	engine.SetMapper(core.SameMapper{})
	//engine.SetMapper(names.GonicMapper{})
	engine.ShowSQL(database.ShowSql)
	//err = engine.Sync2(new(model.IOT_DeviceTopic))
	err = engine.Sync2()
	if err != nil {
		return nil, err
	}
	orm := new(Orm)
	orm.Engine = engine
	DbEngine = orm
	return orm, nil
}
