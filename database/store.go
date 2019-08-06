package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type IPInfo struct {
	Ip            string  `json:"ip"`
	Port          string  `json:"port"`
	LastCheckTime string  `json:"last_check_time"`
	Speed         float64 `json:"speed"`
	Result        []bool  `json:"result"`
	Rate          float64 `json:"rate"`
}

func NewIPInfo(ip string) *IPInfo {
	ipinfo := &IPInfo{}
	ipinfo.Result = []bool{}
	split := strings.Split(ip, ":")
	if len(split) == 2 {
		ipinfo.Ip = split[0]
		ipinfo.Port = split[1]
	}
	return ipinfo
}

func (ipinfo *IPInfo) Reliability() float64 {
	succ := 0
	fail := 0
	for _, v := range ipinfo.Result {
		if v {
			succ++
		} else {
			fail++
		}
	}
	if len(ipinfo.Result) == 0 {
		ipinfo.Rate = 0
	} else {
		ipinfo.Rate = float64(succ) / (float64(succ) + float64(fail))
	}
	return ipinfo.Rate
}

const maxResult = 10
const maxReliability = 0.3

func (ipinfo *IPInfo) Deletable() bool {
	return len(ipinfo.Result) >= maxResult && ipinfo.Rate < maxReliability
}

type Store struct {
	db *leveldb.DB
}

var store *Store
var once sync.Once

func NewStore() *Store {
	once.Do(func() {
		store = &Store{}
		var err error
		leveldbpath := `data\leveldb`
		os.MkdirAll(leveldbpath, os.ModePerm)
		if store.db, err = leveldb.OpenFile(leveldbpath, nil); err != nil {
			log.Fatal(err.Error())
		}
	})
	return store
}

// key -- ipinfo:ip:port
// value -- json

func (s *Store) Put(ipinfo *IPInfo) error {
	key := bytes.NewBufferString(strings.Join([]string{ipinfo.Ip, ipinfo.Port}, ":")).Bytes()

	value, err := s.db.Get(key, nil)
	if err == nil {
		tmp := NewIPInfo("")
		json.Unmarshal(value, tmp)
		ipinfo.Result = append(tmp.Result, ipinfo.Result...)
		length := len(ipinfo.Result)
		if length > maxResult {
			ipinfo.Result = ipinfo.Result[length-maxResult:]
		}
	}
	ipinfo.Reliability()
	if ipinfo.Deletable() {
		return s.db.Delete(key, nil)
	}

	value, err = json.Marshal(ipinfo)
	if err != nil {
		return err
	}
	return s.db.Put(key, value, nil)
}

func (s *Store) GetReliability(num int, limit float64) []*IPInfo {
	result := []*IPInfo{}
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		ip := NewIPInfo("")
		if err := json.Unmarshal(iter.Value(), ip); err != nil {
			continue
		}
		if ip.Rate >= limit {
			result = append(result, ip)
		}
	}
	iter.Release()
	return result
}

func (s *Store) GetAll() []*IPInfo {
	result := []*IPInfo{}
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		ip := NewIPInfo("")
		if err := json.Unmarshal(iter.Value(), ip); err != nil {
			continue
		}
		result = append(result, ip)
	}
	iter.Release()
	return result
}

func CheckStore(sysChan chan<- string) {
	fmt.Println("数据库中IP检测模块启动...")
	store := NewStore()
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			ips := store.GetAll()
			fmt.Printf("数据库中共有%d个IP需要进行检测\n", len(ips))
			for _, ip := range ips {
				time.Sleep(time.Duration(50) * time.Millisecond)
				sysChan <- strings.Join([]string{ip.Ip, ip.Port}, ":")
			}
		}()
		time.Sleep(time.Duration(120) * time.Second)
	}
}
