package repository

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	//gorm
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//NewMysqlRepository 新建mysql数据仓库
func NewMysqlRepository(conn string, maxIdleConns int, maxOpenConns int) (repo *MysqlRepo, err error) {
	db, err := gorm.Open("mysql", conn)
	if err != nil {
		return
	}
	db.LogMode(false)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	return &MysqlRepo{db: db}, err
}

//MysqlRepo 实现基于mysql的数据仓库
type MysqlRepo struct {
	db *gorm.DB
}

//GetTorrentByInfohash 通过infohash查询资源
func (p *MysqlRepo) GetTorrentByInfohash(infohash string) (t Torrent, err error) {
	if len(infohash) == 0 {
		return
	}
	table := "torrent" + string(strings.ToLower(infohash)[0])
	var trt torrent
	if p.db.Table(table).Where("infohash=?", infohash).Scan(&trt).Error != nil {
		return
	}
	if len(trt.Data) > 0 {
		json.Unmarshal([]byte(trt.Data), &t)
	}

	return
}

//BatchGetTorrentByInfohash 通过infohash批量查询资源
func (p *MysqlRepo) BatchGetTorrentByInfohash(infohash []string) (ts []Torrent, err error) {
	if len(infohash) == 0 {
		return
	}
	m := map[string][]string{}
	for _, v := range infohash {
		table := "torrent" + string(strings.ToLower(v)[0])
		m[table] = append(m[table], v)
	}
	for table, value := range m {
		var trt []torrent
		err = p.db.Table(table).Where("infohash in (?)", value).Find(&trt).Error
		if err != nil {
			return
		}
		for _, item := range trt {
			var t Torrent
			info := json.Unmarshal([]byte(item.Data), &t)
			if info == nil {
				ts = append(ts, t)
			}
		}
	}
	return
}

//CreateTorrent 增加资源
func (p *MysqlRepo) CreateTorrent(t Torrent) (err error) {
	if len(t.Infohash) == 0 {
		return
	}
	table := "torrent" + string(strings.ToLower(t.Infohash)[0])

	var trt torrent
	trt.CreateTime = t.CreateTime
	trt.Infohash = t.Infohash
	b, err := json.Marshal(t)
	if err != nil {
		return
	}
	trt.Data = string(b)
	err = p.db.Table(table).Create(&trt).Error
	return
}

//BatchGetInfohash 批量获取infohash
func (p *MysqlRepo) BatchGetInfohash(limit int64) (is []string, err error) {
	if limit <= 0 {
		return
	}
	var data []infohash
	err = p.db.Table("infohash").Limit(limit).Find(&data).Error
	if err != nil {
		return
	}
	for _, v := range data {
		is = append(is, v.Infohash)
	}
	return
}

//BatchDeleteInfohash 批量删除infohash
func (p *MysqlRepo) BatchDeleteInfohash(infohash []string) (err error) {
	if len(infohash) == 0 {
		return
	}
	return p.db.Table("infohash").Where("infohash in (?)", infohash).Delete(nil).Error
}

//CreateInfohash 增加infohash
func (p *MysqlRepo) CreateInfohash(hash string) (err error) {
	return p.db.Table("infohash").Create(&infohash{CreateTime: time.Now(), Infohash: hash}).Error
}

//GetRecommend 获取推荐列表
func (p *MysqlRepo) GetRecommend() (name []string, err error) {
	var rs []recommend
	err = p.db.Table("recommend").Order("id").Find(&rs).Error
	if err != nil {
		return
	}
	for _, v := range rs {
		name = append(name, v.Name)
	}
	return
}

//CreateHistory 增加搜索历史
func (p *MysqlRepo) CreateHistory(keyword, source string) (err error) {
	return p.db.Table("history").Create(&history{Keyword: keyword, Source: source, CreateTime: time.Now()}).Error
}

type torrent struct {
	ID         int64     `gorm:"column:id"`
	Infohash   string    `gorm:"column:infohash"`
	Data       string    `gorm:"column:data"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type infohash struct {
	ID         int64     `gorm:"column:id"`
	Infohash   string    `gorm:"column:infohash"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type history struct {
	ID         int64     `gorm:"column:id"`
	Keyword    string    `gorm:"column:keyword"`
	Source     string    `gorm:"column:source"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type recommend struct {
	ID   int64  `gorm:"column:id"`
	Name string `gorm:"column:name"`
}
