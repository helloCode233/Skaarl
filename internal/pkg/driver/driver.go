package driver

import (
	"Skaarl/config"
	"Skaarl/internal/pkg/helper"
	"Skaarl/internal/pkg/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Driver struct {
	DbPath string
	db     *gorm.DB
	Config *config.Configuration
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

func (d *Driver) InitSqLiteGorm(dbPath string) *Driver {
	if d.db != nil {
		panic("failed to db saved")
	}
	db, err := gorm.Open(sqlite.Open(dbPath), getGormConfig())
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

func (d *Driver) InitLog(ProjectName string) (*Driver, error) {
	project := NewDriver(filepath.Join(".", ProjectName, "skaarl-lock.log")).InitSqLiteGorm(d.DbPath).InitProject()
	project.Put("ProjectName", ProjectName)
	abs, _ := filepath.Abs(filepath.Join(".", ProjectName))
	project.Put("ProjectPath", strings.Replace(abs, "\\", "/", -1))
	return project, project.SaveWireLogs(project.SelectWireFiles(ProjectName))
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
func (d *Driver) GetWireLog() []string {
	var imports []string
	err := d.db.Model(&model.WireLog{}).
		Select("import").
		Group("import").
		Pluck("import", &imports).
		Error
	if err != nil {
		return nil
	}
	return imports
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

func (d *Driver) InitConfig(configPath string) (*Driver, error) {
	if !filepath.IsAbs(configPath) {
		abs, err := filepath.Abs(configPath)
		if err != nil {
			return d, err
		}
		configPath = abs
	}
	env := d.viperRead(filepath.Join(configPath, "env.yaml")).App.Env
	d.Config = d.viperRead(filepath.Join(configPath, env+".yaml"))
	return d, nil
}
func (p *Driver) SelectWireFiles(Path string) map[string]string {
	result := make(map[string]string)
	err := filepath.Walk(Path, func(path string, info os.FileInfo, err error) error {
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
			println(strings.Replace(value, Path+"/", "", -1))
			result[key[strings.Index(key, "func")+4:strings.Index(key, "(")]] = strings.Replace(value, Path+"/", "", -1)
		}
		return err
	})
	if err != nil {
		return nil
	}
	return result
}

func (p *Driver) CheckWireFiles() (bool, map[string]string) {
	selectWireFiles := p.SelectWireFiles(p.Get("ProjectPath").Value)
	var checkWireFiles []*model.WireLog
	err := p.db.Find(&checkWireFiles).Error
	if err != nil || len(checkWireFiles) != len(selectWireFiles) {
		return false, selectWireFiles
	}
	for _, file := range checkWireFiles {
		_, flag := selectWireFiles[file.Func]
		if !flag {
			return false, selectWireFiles
		}
	}
	err = p.SaveWireLogs(selectWireFiles)
	if err != nil {
		return false, selectWireFiles
	}
	return true, selectWireFiles
}
func (p *Driver) InitMySqlGorm(conf *config.Configuration) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Database.UserName,
		conf.Database.Password,
		conf.Database.Host,
		strconv.Itoa(conf.Database.Port),
		conf.Database.Database,
		conf.Database.Charset,
	)
	db, err := gorm.Open(mysql.Open(dsn), getGormConfig())
	if err != nil {
		panic("failed to connect database")
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
		err := db.AutoMigrate()
		if err != nil {
			panic("failed to connect database")
		}
		return db
	}

}
func (d *Driver) GetRemoteDb() *gorm.DB {
	var db *gorm.DB
	switch d.Config.Database.Driver {

	case "mysql":
		db = d.InitMySqlGorm(d.Config)

	case "sqlite":
		db = d.InitSqLiteGorm(d.DbPath).db

	default:
		panic("unknown database driver: " + d.Config.Database.Driver)
	}
	return db
}

type TableInfo struct {
	TableName    string `gorm:"column:TABLE_NAME"`
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

func (d *Driver) GetRemoteDbTables() []*TableInfo {
	db := d.GetRemoteDb()
	var tables []*TableInfo
	// 执行 SQL 查询以获取所有表名和注释
	err := db.Raw("SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?", d.Config.Database.Database).Scan(&tables).Error
	if err != nil {
		fmt.Println("Failed to get tables:", err)
		return nil
	}
	return tables
}

func (p *Driver) GenStart(models string, Dal bool) {
	db := p.GetRemoteDb()
	//dsn := fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/test_gen?charset=utf8mb4&parseTime=True&loc=Local")
	outPath := filepath.Join("internal")
	helper.GenFile(filepath.Join(outPath, "model"), "time", "time.gen", "create", nil)
	abs, _ := filepath.Abs(filepath.Join(outPath, "dal"))
	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath: abs,
		Mode:    gen.WithoutContext,
	})

	var dataMap = map[string]func(gorm.ColumnType) (dataType string){
		// int mapping
		"datetime": func(columnType gorm.ColumnType) (dataType string) {
			return "LocalTime"
		},

		// bool mapping
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.HasPrefix(ct, "tinyint(1)") {
				return "bool"
			}
			return "byte"
		},
	}

	g.WithDataTypeMap(dataMap)
	g.UseDB(db)
	if models == "all" {
		g.ApplyBasic(g.GenerateAllTable()...)
	} else {
		g.ApplyBasic(g.GenerateModel(models))
	}
	g.Execute()
}
