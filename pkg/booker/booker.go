package booker

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/siter"
	"github.com/hiromaily/go-book-teacher/pkg/storages"
	lg "github.com/hiromaily/golibs/log"
)

// ----------------------------------------------------------------------------
// Booker interface
// ----------------------------------------------------------------------------

type Booker interface {
	Start() error
	Cleanup()
}

func NewBooker(
	conf *config.Config,
	interval int,
	storager storages.Storager,
	notifier notifier.Notifier,
	siter siter.Siter) Booker {

	return NewBook(conf, interval, storager, notifier, siter)
}

// ----------------------------------------------------------------------------
// Book
// ----------------------------------------------------------------------------

type Book struct {
	conf     *config.Config
	interval int
	storager storages.Storager
	notifier notifier.Notifier
	siter    siter.Siter
	isLoop   bool
}

func NewBook(
	conf *config.Config,
	interval int,
	storager storages.Storager,
	notifier notifier.Notifier,
	siter siter.Siter) *Book {

	var isLoop bool
	if interval != 0 {
		isLoop = true
	}

	book := Book{
		conf:     conf,
		interval: interval,
		storager: storager,
		notifier: notifier,
		siter:    siter,
		isLoop:   isLoop, //FIXME: it should be changed from dynamic data, testmode, heroku env should be false
	}
	return &book
}

func (b *Book) Start() error {

	if err := b.siter.FetchInitialData(); err != nil {
		return errors.Wrap(err, "fail to call siter.FetchInitialData()")
	}

	for {
		//FIXME: this logic would be better to move into siter/dmm.go
		//reset
		b.siter.InitializeSavedTeachers()

		//scraping
		b.siter.HandleTeachers()

		//save
		b.saveAndNotify()

		//TODO:when integration test, send channel

		//execute only once
		if !b.isLoop {
			b.storager.Close()
			return nil
		}

		time.Sleep(time.Duration(b.interval) * time.Second)
	}
}

func (b *Book) Cleanup() {
	b.storager.Close()
}

//check saved data and run browser if needed
func (b *Book) saveAndNotify() {
	ths := b.siter.GetSavedTeachers()
	if len(ths) != 0 {
		//create string from ids slice
		var sum int
		for _, t := range ths {
			sum += t.ID
		}
		newData := strconv.Itoa(sum)

		// save
		ok, err := b.storager.Save(newData)
		if err != nil {
			lg.Errorf("fail to save() %v", err)
		}

		if ok {
			//notify
			b.notifier.Send(ths)
		}
	}
}

// ----------------------------------------------------------------------------
// DummyBook
// ----------------------------------------------------------------------------

type DummyBook struct{}

func NewDummyBook() *DummyBook {
	return &DummyBook{}
}

func (b *DummyBook) Start() error {
	return nil
}

func (b *DummyBook) Cleanup() {}