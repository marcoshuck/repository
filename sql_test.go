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

func (s *SQLTestSuite) TestGet() {
	s.createMockData()
	result, err := s.repository.Get(context.Background(), 1)

	s.Assert().NoError(err)
	s.Assert().Equal("Marcos", result.FirstName)

	result, err = s.repository.Get(context.Background(), 55)
	s.Assert().Error(err)
	s.Assert().ErrorIs(err, gorm.ErrRecordNotFound)
	s.Assert().Zero(result)
}

func (s *SQLTestSuite) TestCreate() {
	s.createMockData()
	created, err := s.repository.Create(context.Background(), Test{
		FirstName: "Another",
		LastName:  "Test",
	})
	s.Assert().NoError(err)

	result, err := s.repository.Get(context.Background(), created.ID)
	s.Assert().NoError(err)

	s.Assert().Equal(created.FirstName, result.FirstName)
}

func (s *SQLTestSuite) TestUpdate() {
	s.createMockData()

	before, err := s.repository.Get(context.Background(), 2)
	s.Assert().NoError(err)
	s.Assert().Equal("Andres", before.FirstName)

	update := before
	update.FirstName = "Marcos"
	updated, err := s.repository.Update(context.Background(), 2, update)
	s.Assert().NoError(err)
	s.Assert().Equal(update.ID, updated.ID)
	s.Assert().Equal(update.FirstName, updated.FirstName)

	after, err := s.repository.Get(context.Background(), 2)
	s.Assert().NoError(err)
	s.Assert().Equal(before.ID, after.ID)
	s.Assert().NotEqual(before.FirstName, after.FirstName)

	updated, err = s.repository.Update(context.Background(), 100, update)
	s.Assert().Error(err)
	s.Assert().ErrorIs(err, gorm.ErrRecordNotFound)
	s.Assert().Zero(updated)
}

func (s *SQLTestSuite) TestRemove() {
	s.createMockData()

	before, err := s.repository.Get(context.Background(), 2)
	s.Assert().NoError(err)
	s.Assert().False(before.DeletedAt.Valid)

	result, err := s.repository.Remove(context.Background(), 2)
	s.Assert().NoError(err)
	s.Assert().True(result.DeletedAt.Valid)

	after, err := s.repository.Get(context.Background(), 2)
	s.Assert().Error(err)
	s.Assert().Zero(after)
	s.Assert().ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *SQLTestSuite) TestCreateBulk() {}

func (s *SQLTestSuite) TestUpdateBulk() {}

func (s *SQLTestSuite) TestRemoveBulk() {}
