package sqlstore_test

import (
	"testing"

	"github.com/qeery8/rest/internal/app/model"
	"github.com/qeery8/rest/internal/app/store"
	"github.com/qeery8/rest/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u.ID)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	s.User().Create(u1)
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	_, err := s.User().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(u1)
	u2, err := s.User().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	foundUser, err := s.User().FindByEmail(u.Email)
	assert.NoError(t, err)
	assert.True(t, foundUser.ComparePassword("password"))

	foundUser.Password = "newpassword123"
	err = s.User().Update(foundUser)
	assert.NoError(t, err, "OK")

	updatedUser, err := s.User().Find(foundUser.ID)
	assert.NoError(t, err)

	assert.False(t, updatedUser.ComparePassword("password"))

	assert.True(t, updatedUser.ComparePassword("newpassword123"))

	assert.Equal(t, u.Email, updatedUser.Email)
}

func TestUserREpository_Delete(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	foundUser, err := s.User().Find(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, foundUser.Email)

	err = s.User().Delete(u.ID)
	assert.NoError(t, err)

	_, err = s.User().Find(u.ID)
	assert.Error(t, err)
	assert.Equal(t, store.ErrRecordNotFound, err)

	_, err = s.User().FindByEmail(u.Email)
	assert.Error(t, err)
	assert.Equal(t, store.ErrRecordNotFound, err)
}
