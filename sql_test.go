package repository

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type Test struct {
	gorm.Model
	FirstName string
	LastName  string
}

func TestSQL(t *testing.T) {
	suite.Run(t, new(SQLTestSuite))
}

type SQLTestSuite struct {
	suite.Suite

	db *gorm.DB

	tx *gorm.DB

	repository Repository[Test]
}

func (s *SQLTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open("sql_test.db"))
	s.Require().NoError(err)
	s.Require().NotNil(db)
	s.db = db
}

func (s *SQLTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repository = &SQL[Test]{
		db: s.tx,
	}
	s.Require().NoError(s.tx.Migrator().AutoMigrate(&Test{}))
}

func (s *SQLTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *SQLTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	s.Require().NoError(sqlDB.Close())
	s.Require().NoError(os.Remove("sql_test.db"))
}

func (s *SQLTestSuite) createMockData() {
	s.Require().NoError(s.tx.Create(&Test{
		FirstName: "Marcos",
		LastName:  "Huck",
	}).Error)

	s.Require().NoError(s.tx.Create(&Test{
		FirstName: "Andres",
		LastName:  "Huck",
	}).Error)

	s.Require().NoError(s.tx.Create(&Test{
		FirstName: "Andrew",
		LastName:  "Baker",
	}).Error)
}

func (s *SQLTestSuite) TestFind() {
	s.createMockData()
	list, err := s.repository.Find(context.Background(), []uint{1, 2, 3})
	s.Assert().NoError(err)
	s.Assert().Len(list, 3)

	list, err = s.repository.Find(context.Background(), []uint{1, 2})
	s.Assert().NoError(err)
	s.Assert().Len(list, 2)

	list, err = s.repository.Find(context.Background(), []uint{55, 12})
	s.Assert().NoError(err)
	s.Assert().Empty(list)
}

// 	// Create creates an entity in a persistence layer.
//	Create(ctx context.Context, entity E) (E, error)
//	// CreateBulk creates a set of entities in a persistence layer.
//	CreateBulk(ctx context.Context, entities []E) ([]E, error)
//	// Get returns an entity from a persistence layer identified by its ID. It returns an error if the entity doesn't exist.
//	Get(ctx context.Context, id uint) (E, error)
//	// Find returns a set of entities from a persistence layer identified by their ID. It returns
//	// an empty slice if no records were found.
//	Find(ctx context.Context, ids []uint) ([]E, error)
//	// Update updates with the values of entity the entity identified by id.
//	Update(ctx context.Context, id uint, entity E) (E, error)
//	// UpdateBulk updates with the values of entity all the elements identified by the slice of ids.
//	UpdateBulk(ctx context.Context, ids []uint, entity E) ([]E, error)
//	// Remove removes the given id from a persistence layer.
//	Remove(ctx context.Context, id uint) (E, error)
//	// RemoveBulk removes a set of elements from a persistence layer.
//	RemoveBulk(ctx context.Context, ids []uint) ([]E, error)
