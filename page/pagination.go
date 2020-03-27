package page

import (
	"math"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

//Paginator for page use
type Paginator struct {
	TotalCount int64 //`json:"total_record"`
	TotalPage  int64 // `json:"total_page"`
	Page       int64 //`json:"page"`
	PageSize   int64 //`json:"page_size"`
}

// DoPage ip
func (p *Paginator) DoPage(table *gorm.DB, list interface{}, order []string) error {
	//table.Or
	return Page(table, p, list, order)
}

//Page so
func Page(table *gorm.DB, p *Paginator, list interface{}, order []string) error {

	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	done := make(chan bool, 1)
	go countRecords(table, done, &p.TotalCount)
	offset := (p.Page - 1) * p.PageSize
	if len(order) > 0 {
		for _, element := range order {
			table = table.Order(element)
		}
	}
	err := table.Limit(p.PageSize).Offset(offset).Find(list).Error
	<-done
	if err != nil {
		log.Errorf("Query countRecords %v", err)
		return err
	}
	p.TotalPage = int64(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
	if p.TotalPage < p.Page {
		p.Page = p.TotalPage
	}

	return nil
}

func countRecords(table *gorm.DB, done chan bool, count *int64) {
	err := table.Count(count).Error
	if err != nil {
		log.Errorf("Query countRecords %v", err)
	}
	done <- true
}
