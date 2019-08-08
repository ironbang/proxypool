package database

import (
	"fmt"
	"github.com/ironbang/proxypool/common/function"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"sync"
	"time"
)

const maxResult = 10
const maxReliability = 0.3

type ProxyIPInfo struct {
	ID     uint   `json:"-"`
	IpPort string `json:"ip"`
	// Port	int	`gorm:"not null"`
	Http          bool
	Https         bool
	Socks         bool
	LastCheckTime string  `json:"last_check_time"`
	Speed         float64 `json:"speed"`
	Results       string  `json:"-"`
	Rate          float64 `json:"rate"`
	Created       string  `json:"-"`
	Updated       string  `json:"-"`
}

func (ipinfo *ProxyIPInfo) CalcRate() {
	total := len(ipinfo.Results)
	if total == 0 {
		ipinfo.Rate = 0
	} else {
		succ := 0
		for _, v := range ipinfo.Results {
			if v == '1' {
				succ++
			}
		}
		ipinfo.Rate = float64(succ) / float64(total)
	}
}

func (ipinfo ProxyIPInfo) Deletable() bool {
	return len(ipinfo.Results) >= maxResult && ipinfo.Rate < maxReliability
}

type Store struct {
	db *gorm.DB
}

var store *Store
var once sync.Once

func NewStore() *Store {
	once.Do(func() {
		store = &Store{}
		var err error
		dbpath := `./data/proxypool.db`
		store.db, err = gorm.Open("sqlite3", dbpath)
		if err != nil {
			log.Fatal(err.Error())
		}
		// store.db.DropTableIfExists(&ProxyIPInfo{})
		if !store.db.HasTable(&ProxyIPInfo{}) {
			store.db.CreateTable(&ProxyIPInfo{})
		}
	})
	return store
}

// key -- ipinfo:ip:port
// value -- json

func (s *Store) Put(ipinfo *ProxyIPInfo) (err error) {
	proxy := &ProxyIPInfo{}
	err = s.db.Where("ip_port = ?", ipinfo.IpPort).First(proxy).Error
	ipinfo.Results = proxy.Results + ipinfo.Results
	ipinfo.CalcRate()
	if err == nil { // 找到数据
		ipinfo.ID = proxy.ID
		if ipinfo.Deletable() {
			err = s.db.Delete(ipinfo).Error
		} else {
			// 更新
			ipinfo.Created = proxy.Created
			ipinfo.Updated = function.FormatTime(time.Now())
			err = s.db.Save(ipinfo).Error
		}
	} else {
		// 新建
		ipinfo.Created = function.FormatTime(time.Now())
		ipinfo.Updated = function.FormatTime(time.Now())
		err = s.db.Create(ipinfo).Error
	}
	return
}

func (s *Store) GetReliability(num int, limit float64) (ips []*ProxyIPInfo, err error) {
	err = s.db.Order("rate desc").Limit(num).Where("rate >= ?", limit).Find(&ips).Error
	return
}

func (s *Store) GetAll() (ips []*ProxyIPInfo, err error) {
	// ips = []ProxyIPInfo{}
	err = s.db.Find(&ips).Error
	return
}

func CheckStore(sysChan chan<- string, group *sync.WaitGroup) {
	fmt.Println("数据库中IP检测模块启动...")
	store := NewStore()
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			ips, err := store.GetAll()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("数据库中共有%d个IP需要进行检测\n", len(ips))
			for _, ip := range ips {
				sysChan <- ip.IpPort
			}
		}()
		time.Sleep(time.Duration(10) * time.Minute)
	}
}
