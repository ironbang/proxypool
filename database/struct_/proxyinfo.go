package struct_

import (
	"github.com/ironbang/proxypool/common/function"
	"github.com/ironbang/proxypool/database"
	"time"
)

const maxResult = 10
const maxReliability = 0.3

func init() {
	db := database.GetDB()
	if !db.HasTable(&ProxyIPInfo{}) {
		db.CreateTable(&ProxyIPInfo{})
	}
}

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
	Checked       bool    `json:"-"`
}

func (p *ProxyIPInfo) Update() (err error) {
	db := database.GetDB()
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	p.Checked = true
	p.CalcRate()
	if p.Deletable() {
		err = tx.Delete(p).Error
	} else {
		// 更新
		p.Updated = function.FormatTime(time.Now())
		err = tx.Save(p).Error
	}
	return
}

func (p *ProxyIPInfo) Insert() (err error) {
	db := database.GetDB()
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if tx.Where("ip_port = ?", p.IpPort).First(&ProxyIPInfo{}).RecordNotFound() {
		p.Created = function.FormatTime(time.Now())
		p.Updated = function.FormatTime(time.Now())
		err = tx.Create(p).Error
	}
	return
}

func (p *ProxyIPInfo) Deletable() bool {
	return len(p.Results) >= maxResult && p.Rate < maxReliability
}

func (p *ProxyIPInfo) CalcRate() {
	total := len(p.Results)
	if total == 0 {
		p.Rate = 0
	} else {
		if total > maxResult {
			p.Results = p.Results[total-maxResult : total-1]
		}
		succ := 0
		for _, v := range p.Results {
			if v == '1' {
				succ++
			}
		}
		p.Rate = float64(succ) / float64(total)
	}
}

func GetReliability(num int, limit float64) (ips []*ProxyIPInfo, err error) {
	db := database.GetDB()
	err = db.Order("rate desc").Limit(num).Where("rate >= ? and checked=true", limit).Find(&ips).Error
	return
}

func GetAll() (ips []*ProxyIPInfo, err error) {
	// ips = []ProxyIPInfo{}
	db := database.GetDB()
	err = db.Find(&ips).Error
	return
}
