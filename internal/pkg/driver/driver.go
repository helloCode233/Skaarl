package driver

import (
	"Skaarl/config"
	"Skaarl/internal/pkg/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
		err := d.db.AutoMigrate(&model.WireLog{}, &model.Cache{})
		if err != nil {
			return nil
		}
	}
	return d
}

func (d *Driver) SaveWireLogs(files map[string]string) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range files {
			wrieLog := &model.WireLog{Func: key, Import: value}
			err := tx.Create(wrieLog).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (d *Driver) Put(Key, Value string) error {
	d.db.Unscoped().Where("key = ?", Key).Delete(&model.Cache{})
	return d.db.Create(&model.Cache{Key: Key, Value: Value}).Error
}
func (d *Driver) Get(Key string) *model.Cache {
	var cache model.Cache
	err := d.db.First(&cache, "key = ?", Key).Error
	if err != nil {
		return nil
	}
	return &cache
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
	env := d.viperRead(filepath.Join(configPath, "env.yaml")).App.Env
	d.config = d.viperRead(filepath.Join(configPath, env+".yaml"))
	return nil
}
func (p *Driver) SelectWireFiles() map[string]string {
	result := make(map[string]string)
	err := filepath.Walk(p.Get("ProjectName").Value, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// 匹配带 @wire 标记的函数定义（更精确的版本）
		regexWireFunc := `(?s)//\s*@wire:\w+\s*\nfunc.*?\s*\{`
		// 正则表达式匹配带 @wire 标记的函数
		wireFuncRe := regexp.MustCompile(regexWireFunc)
		matches := wireFuncRe.FindAllStringSubmatch(string(data), -1)
		value := strings.Replace(filepath.Dir(path), "\\", "/", -1)
		for _, match := range matches {
			key := match[0]
			result[key[strings.Index(key, "func"):len(key)-2]] = value
		}
		return err
	})
	if err != nil {
		return nil
	}
	return result
}

func (p *Driver) CheckWireFiles() bool {
	selectWireFiles := p.SelectWireFiles()
	var checkWireFiles []*model.WireLog
	err := p.db.Find(&checkWireFiles).Error
	if err != nil || len(checkWireFiles) != len(selectWireFiles) {
		return false
	}
	return true
}
