package repository

import (
	"testing"
	"time"
)

var (
	mysqlRepo Repository
)

func initMysqlRepo(t *testing.T) {
	conn := "root:@tcp(127.0.0.1:3306)/torrent?charset=utf8&parseTime=True&loc=Local"
	repo, err := NewMysqlRepository(conn, 128, 128)
	if err != nil {
		t.Fatal(err)
	}
	mysqlRepo = repo
}

func TestGetTorrentByInfohash(t *testing.T) {
	initMysqlRepo(t)
	trt, err := mysqlRepo.GetTorrentByInfohash("11")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(trt)
}

func TestBatchGetTorrentByInfohash(t *testing.T) {
	initMysqlRepo(t)
	trt, err := mysqlRepo.BatchGetTorrentByInfohash([]string{"11", "01"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(trt)
}

func TestCreateTorrent(t *testing.T) {
	initMysqlRepo(t)
	var trt Torrent
	trt.CreateTime = time.Now()
	trt.FileCount = 12
	trt.Files = append(trt.Files, File{Name: "test", Length: 1024})
	trt.Infohash = "21"
	trt.Length = 2048
	trt.Name = "just test"
	err := mysqlRepo.CreateTorrent(trt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBatchGetInfohash(t *testing.T) {
	initMysqlRepo(t)
	is, err := mysqlRepo.BatchGetInfohash(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(is)
}

func TestBatchDeleteInfohash(t *testing.T) {
	initMysqlRepo(t)
	err := mysqlRepo.BatchDeleteInfohash([]string{"01", "11"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateInfohash(t *testing.T) {
	initMysqlRepo(t)
	err := mysqlRepo.CreateInfohash("11")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRecommend(t *testing.T) {
	initMysqlRepo(t)
	name, err := mysqlRepo.GetRecommend()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(name)
}

func TestCreateHistory(t *testing.T) {
	initMysqlRepo(t)
	err := mysqlRepo.CreateHistory("x", "s")
	if err != nil {
		t.Fatal(err)
	}
}
