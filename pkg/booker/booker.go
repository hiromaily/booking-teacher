package booker

import (
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/models"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	storages "github.com/hiromaily/go-book-teacher/pkg/storage"
)

// ----------------------------------------------------------------------------
// Booker interface
// ----------------------------------------------------------------------------

// Booker is interface
type Booker interface {
	Start() error
	Cleanup()
	Close()
}

// NewBooker is to return booker interface
func NewBooker(
	storager storages.Storager,
	notifier notifier.Notifier,
	siter site.Siter,
	logger *zap.Logger,
	day int,
	interval int,
) Booker {
	return NewBook(
		storager,
		notifier,
		siter,
		logger,
		day,
		interval,
	)
}

// ----------------------------------------------------------------------------
// Book
// ----------------------------------------------------------------------------

// Book is Book object
type Book struct {
	storager storages.Storager
	notifier notifier.Notifier
	siter    site.Siter
	logger   *zap.Logger
	day      int
	interval int
	isLoop   bool
}

// NewBook is to return book object
func NewBook(
	storager storages.Storager,
	notifier notifier.Notifier,
	siter site.Siter,
	logger *zap.Logger,
	day int,
	interval int,
) *Book {

	var isLoop bool
	if interval != 0 {
		isLoop = true
	}

	book := Book{
		storager: storager,
		notifier: notifier,
		siter:    siter,
		logger:   logger,
		day:      day,
		interval: interval,
		isLoop:   isLoop, // TODO: testmode, heroku env should be false
	}
	return &book
}

// Start is to start book execution
func (b *Book) Start() error {
	b.logger.Debug("book Start()")
	defer b.storager.Close()

	b.logger.Debug("book siter.FetchInitialData()")
	if err := b.siter.FetchInitialData(); err != nil {
		return errors.Wrap(err, "fail to call site.FetchInitialData()")
	}

	for {
		// scraping
		b.logger.Debug("book siter.FindTeachers()")
		teachers := b.siter.FindTeachers(b.day)

		// save
		b.logger.Debug("book siter.saveAndNotify()")
		b.saveAndNotify(teachers)

		// execute only once
		if !b.isLoop {
			return nil
		}

		b.logger.Debug("book sleep for next execution")
		time.Sleep(time.Duration(b.interval) * time.Second)
	}
}

// Cleanup is to clean up middleware object
func (b *Book) Cleanup() {
	b.storager.Delete()
}

// Close is to clean up middleware object
func (b *Book) Close() {
	b.storager.Close()
}

// saveAndNotify is to save and notify if something saved
func (b *Book) saveAndNotify(ths []models.TeacherInfo) {
	if len(ths) != 0 {
		// create string from ids slice
		var sum int
		for _, t := range ths {
			sum += t.ID
		}
		newData := strconv.Itoa(sum)

		// save
		isUpdated, err := b.storager.Save(newData)
		if err != nil {
			b.logger.Error("fail to call Save()", zap.Error(err))
		}

		if isUpdated {
			// notify
			b.notifier.Notify(ths)
		}
	}
}

// ----------------------------------------------------------------------------
// DummyBook
// ----------------------------------------------------------------------------

// DummyBook is DummyBook object
type DummyBook struct{}

// NewDummyBook is to return NewDummyBook object
func NewDummyBook() *DummyBook {
	return &DummyBook{}
}

// Start is to do nothing
func (b *DummyBook) Start() error {
	return nil
}

// Cleanup is to do nothing
func (b *DummyBook) Cleanup() {}

// Close is to do nothing
func (b *DummyBook) Close() {}
