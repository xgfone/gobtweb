package repository

import (
	"errors"
	"time"
)

//定义通用变量
var (
	ErrNotFound = errors.New("not found")
)

//Torrent 定义了资源详细信息
type Torrent struct {
	Infohash   string
	Name       string
	Length     int64
	Heat       int64
	FileCount  int64
	Files      []File
	CreateTime time.Time
}

//File 定义了资源中包含的文件
type File struct {
	Name   string
	Length int64
}

//Repository 定义了数据仓库需要实现的接口
type Repository interface {
	//Torrent
	GetTorrentByInfohash(infohash string) (torrent Torrent, err error)
	BatchGetTorrentByInfohash(infohash []string) (torrent []Torrent, err error)
	CreateTorrent(torrent Torrent) (err error)

	//Infohash
	BatchGetInfohash(limit int64) (infohash []string, err error)
	BatchDeleteInfohash(infohash []string) (err error)
	CreateInfohash(infohash string) (err error)

	//recommend
	GetRecommend() (name []string, err error)

	//history
	CreateHistory(keyword, source string) (err error)
}
