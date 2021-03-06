package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zheng-ji/goSnowFlake"
)

type Model struct {
	ID          int64      `gorm:"primary_key"`
	CreatedAt   *time.Time `gorm:"type:datetime"  json:",omitempty"`
	UpdatedAt   *time.Time `gorm:"type:datetime"  json:",omitempty"`
	DeletedAt   *time.Time `gorm:"type:datetime" sql:"index" json:",omitempty"`
	Description string     `json:",omitempty"`
}

func JsonToObj(jsonStr string) map[string]string {
	var jsonMap map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &jsonMap); err == nil {
		return jsonMap
	} else {
		beego.Error(err)
		return nil
	}
}

func DBModelToJson(v interface{}) string {
	if bytes, err := json.Marshal(v); err == nil {
		r := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\[+-]?\d{2}:\d{2})`)
		jsonstr := r.ReplaceAllStringFunc(string(bytes), func(s string) string {
			if t, err := time.Parse(time.RFC3339Nano, s); err != nil {
				panic(err)
			} else {
				return t.Format("2006-01-02 15:04:05")
			}
		})
		return jsonstr
	} else {
		return err.Error()
	}
}

var connstr string

func DB() (*gorm.DB, error) {
	if db, err := gorm.Open("mysql", connstr); err != nil {
		return nil, err
	} else {
		if err = db.DB().Ping(); err != nil {
			return nil, err
		}
		if beego.AppConfig.String("runmode") == "dev" {
			//db.LogMode(true)
		}
		return db, nil
	}
}

var Worker *goSnowFlake.IdWorker

func init() {

	beego.Info("Initializing IdWorker...")
	var err error
	Worker, err = goSnowFlake.NewIdWorker(1)
	if err != nil {
		fmt.Println(err)
	}
	beego.Info("Initializing Database...")
	connstr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		beego.AppConfig.String("mysqluser"),
		beego.AppConfig.String("mysqlpass"),
		beego.AppConfig.String("mysqladdr"),
		beego.AppConfig.String("mysqlport"),
		beego.AppConfig.String("mysqldb"),
	)
	if db, err := DB(); err != nil {
		beego.Error("Initialize Database Error!")
		beego.Error(err)
		return
	} else {
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&User{}, &UserGroup{},
			&Permission{}, &PermissionGroup{},
			&Role{}, &RoleGroup{},
		)
		beego.Info("Database Initializing Success！")
	}
	//InitTestData()
}

func testdb() {

}

func InitTestData() {
	ids := GetMany(100)

	role10 := Role{Model: Model{ID: ids[3]}, Platform: "平台1", Name: "管理员"}

	role11 := Role{Model: Model{ID: ids[2]}, Platform: "平台1", Name: "角色1"}

	role12 := Role{Model: Model{ID: ids[1]}, Platform: "平台1", Name: "角色2"}

	SaveRole(&role10)
	SaveRole(&role11)
	SaveRole(&role12)
	role20 := Role{Platform: "平台2", Name: "管理员"}
	SaveRole(&role20)
	role30 := Role{Platform: "平台3", Name: "管理员"}
	role31 := Role{Platform: "平台3", Name: "角色1"}
	role32 := Role{Platform: "平台3", Name: "角色2"}
	SaveRole(&role30)
	SaveRole(&role31)
	SaveRole(&role32)

	rg0 := RoleGroup{Name: "分组1", ParentID: 0, RoleGroups: []RoleGroup{
		RoleGroup{Name: "2级分组"}, RoleGroup{Name: "分组12"}, RoleGroup{Name: "分组13"},
	},
	}
	rg1 := RoleGroup{Name: "分组2", ParentID: 0, RoleGroups: []RoleGroup{
		RoleGroup{Name: "2级分组", Roles: []Role{role10, role11, role12}}, RoleGroup{Name: "分组22"}, RoleGroup{Name: "分组23"},
	}}
	rg2 := RoleGroup{Name: "分组3", ParentID: 0}
	rg3 := RoleGroup{Name: "分组4", ParentID: 0}
	rg4 := RoleGroup{Name: "分组5", ParentID: 0}

	SaveRoleGroup(&rg0)
	SaveRoleGroup(&rg1)
	SaveRoleGroup(&rg2)
	SaveRoleGroup(&rg3)
	SaveRoleGroup(&rg4)

	user01 := User{Model: Model{ID: ids[10]}, Account: "admin", Password: "e10adc3949ba59abbe56e057f20f883e", Name: "管理员", QQ: "2342342342"}
	user02 := User{Model: Model{ID: ids[11]}, Account: "admin1", Password: "e10adc3949ba59abbe56e057f20f883e", Name: "张津杰", QQ: "982372873"}
	user03 := User{Model: Model{ID: ids[13]}, Account: "admin2", Password: "e10adc3949ba59abbe56e057f20f883e", Name: "于中玮", QQ: "610750125"}
	user04 := User{Model: Model{ID: ids[12]}, Account: "admin3", Password: "e10adc3949ba59abbe56e057f20f883e", Name: "赵毅", QQ: "982639692"}
	user05 := User{Model: Model{ID: ids[20]}, Account: "admin4", Password: "e10adc3949ba59abbe56e057f20f883e", Name: "程业俊", QQ: "2987329723"}
	SaveUser(&user01)
	SaveUser(&user02)
	SaveUser(&user03)
	SaveUser(&user04)
	SaveUser(&user05)

	p11 := Permission{Platform: "平台1", Name: "线路查看", Path: "/Line/View"}
	p12 := Permission{Platform: "平台1", Name: "线路编辑", Path: "/Line/Edit"}
	p13 := Permission{Platform: "平台1", Name: "站点修改", Path: "/Stop/Modify"}
	p21 := Permission{Platform: "平台1", Name: "站点删除", Path: "/Stop/Update"}
	p22 := Permission{Platform: "平台1", Name: "新增站点", Path: "/Stop/Create"}
	SavePermission(&p11)
	SavePermission(&p12)
	SavePermission(&p13)
	SavePermission(&p21)
	SavePermission(&p22)
}
