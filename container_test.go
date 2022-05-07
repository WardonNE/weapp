package weapp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	assert.Equal(t, &container{constructor: "Init"}, newContainer())
}

func TestSetConstruct(t *testing.T) {
	container := newContainer()
	container.SetConstructor("Construct")
	assert.Equal(t, "Construct", container.constructor)
}

type TestBase struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ITestLogger interface {
	Output() string
}

type TestLogger struct {
}

func (logger *TestLogger) Output() string {
	return "logger"
}

type ITestLogger2 interface {
	Output2() string
}

type TestLogger2 struct {
}

func (logger *TestLogger2) Output2() string {
	return "logger2"
}

type TestUser struct {
	Attributes map[string]any `inject:""`
	ID         int
	Name       string
	Tag        []string     `inject:""`
	Company    *TestCompany `inject:""`
	Group      *TestGroup   `inject:"group"`
	Logger     ITestLogger  `inject:"logger"`
	Logger2    ITestLogger2 `inject:""`
	TestBase   `inject:""`
}

func (t *TestUser) Init() {
	t.ID = 1
	t.Name = "testUser"
	t.Tag = append(t.Tag, "Tag1")
	t.Attributes["gender"] = "female"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

type TestGroup struct {
	ID int
}

type TestCompany struct {
	ID       int
	Name     string
	Country  *TestCountry `inject:"country"`
	TestBase `inject:""`
}

func (t *TestCompany) Init() {
	t.ID = 2
	t.Name = "testCompany"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

type TestCountry struct {
	ID       int
	Name     string
	TestBase `inject:""`
}

func (t *TestCountry) Init() {
	t.ID = 3
	t.Name = "testCountry"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func TestStoreWithDefaultInit(t *testing.T) {
	container := newContainer()
	testUser := new(TestUser)
	err := container.Store("logger", new(TestLogger), nil)
	err = container.Store("logger2", new(TestLogger2), nil)
	err = container.Store("country", new(TestCountry), nil)
	err = container.Store("user", testUser, nil)
	testUser.Attributes["gender"] = "female"
	assert.Equal(t, nil, err)
	assert.Equal(t, "female", testUser.Attributes["gender"])
	assert.Equal(t, "testUser", testUser.Name)
	assert.Equal(t, 1, testUser.ID)
	assert.Equal(t, 2, testUser.Company.ID)
	assert.Equal(t, "testCompany", testUser.Company.Name)
	assert.Equal(t, 3, testUser.Company.Country.ID)
	assert.Equal(t, "testCountry", testUser.Company.Country.Name)
}

func TestStoreWithCustomInit(t *testing.T) {
	container := newContainer()
	testUser := new(TestUser)
	err := container.Store("logger", new(TestLogger), nil)
	err = container.Store("logger2", new(TestLogger2), nil)
	err = container.Store("country", new(TestCountry), nil)
	err = container.Store("user", testUser, func(user any) {
		user.(*TestUser).ID = 1
		user.(*TestUser).Company.ID = 2
		user.(*TestUser).Company.Country.ID = 3
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, testUser.ID)
	assert.Equal(t, 2, testUser.Company.ID)
	assert.Equal(t, 3, testUser.Company.Country.ID)
}
