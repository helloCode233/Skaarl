package driver

import (
	"Skaarl/config"
	"Skaarl/internal/pkg/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path/filepath"
)

type Driver struct {
	DbPath string
	db     *gorm.DB
	config *config.Configuration
}

func NewDriver(Path string) *Driver {
	return &Driver{
		DbPath: Path,
		db:     nil,
	}
}

func getGormConfig() *gorm.Config {
	gorm_conf := &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	}

	return gorm_conf
}

func (d *Driver) InitSqLiteGorm() *Driver {
	if d.db != nil {
		panic("failed to db saved")
	}
	db, err := gorm.Open(sqlite.Open(d.DbPath), getGormConfig())
	if err != nil {
		panic("failed to connect database")
	} else {
		//err := db.AutoMigrate(&model.Conn{})
		//if err != nil {
		//	print(err.Error())
		//}
	}
	d.db = db
	return d
}

func (d *Driver) InitProject() *Driver {
	if d.db != nil {
		err := d.db.AutoMigrate(&model.WireLog{})
		if err != nil {
			return nil
		}
	}
	return d
}

func (d *Driver) SaveWireLogs(files map[string][]string) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for key, values := range files {
			for _, value := range values {
				println(key + "|" + value)
				wrieLog := &model.WireLog{Func: value, Import: key}
				err := tx.Create(wrieLog).Error
				if err != nil {
					return err
				}
			}
		}

		// 返回 nil 提交事务
		return nil
	})
	return err
}

func (d *Driver) viperRead(configPath string) *config.Configuration {
	var conf *config.Configuration
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}
	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}
	return conf
}

func (d *Driver) InitConfig(configPath string) error {

	if !filepath.IsAbs(configPath) {
		abs, err := filepath.Abs(configPath)
		if err != nil {
			return err
		}
		configPath = abs
	}
	//fmt.Println("load config:" + configPath)
	env := d.viperRead(filepath.Join(configPath, "env.yaml")).App.Env
	d.config = d.viperRead(filepath.Join(configPath, env+".yaml"))
	return nil
}
