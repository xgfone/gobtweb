package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v3"
)

type TorrentSearch struct {
	Infohash   string
	Name       string
	Length     int64
	Heat       int64
	CreateTime time.Time
}

// func GetTorrentByInfohashFromDB(client repository.Repository, infohash string) (*repository.Torrent, error) {
// 	if t, err := client.GetTorrentByInfohash(infohash); err != nil || t.Infohash != infohash {
// 		return nil, err
// 	} else {
// 		return *t, nil
// 	}
// }

func GetTorrentByInfohashFromSE(client *elastic.Client, infohash string) (*TorrentSearch, error) {
	indexType := strings.ToLower(string(infohash[0]))
	result, err := client.Get().Index("torrent").Type(indexType).Id(infohash).Do()
	if err != nil || result == nil || !result.Found || result.Source == nil {
		return nil, fmt.Errorf("can't get the result")
	}

	var data TorrentSearch
	err = json.Unmarshal(*result.Source, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func SearchKeyword(client *elastic.Client, key string, from, size int) ([]TorrentSearch, int, error) {
	// _keys := make([]interface{}, 0, len(keys))
	// for _, k := range keys {
	// 	_keys = append(_keys, k)
	// }
	//query := elastic.NewTermsQuery("Name", _keys...)
	query := elastic.NewQueryStringQuery(key)

	result, err := client.Search().Index("torrent").Query(query).From(from).Size(size).Do()
	if err != nil {
		return nil, 0, err
	}

	torrents := make([]TorrentSearch, 0)
	if result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {
			var torrent TorrentSearch
			if err := json.Unmarshal(*hit.Source, &torrent); err == nil {
				torrent.Infohash = hit.Id
				torrents = append(torrents, torrent)
			}
		}
	}

	return torrents, int(result.Hits.TotalHits), nil
}
