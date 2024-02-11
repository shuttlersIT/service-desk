// agents_test.go
package models

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	//"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	//"github.com/jinzhu/gorm"
	//"github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/mock"
)

type MockClauseBuilder struct {
	mock.Mock
}

func (m *MockClauseBuilder) Build(clause clause.Clause) {
	m.Called(clause)
}

type ClauseBuilder interface {
	Build(clause.Clause)
	// Add more methods as needed
}

func TestCreateAgent_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dsn := "root:1T$hutt!ers@tcp(localhost:3306)/mydb?parseTime=true"
	// Create GORM dialect from mock db
	dialect := mysql.Open(dsn)

	// Open GORM connection using dialect
	gormDB, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening GORM db: %v", err)
	}
	//defer gormDB.Close()

	model := NewAgentDBModel(gormDB)
	//var role Roles
	// Mock expected database interaction for successful agent creation
	expectedAgentID := uint(1)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedAgentID)
	mock.ExpectQuery(`INSERT INTO agents`).
		WithArgs("John Doe", "johndoe@example.com", "1234567890", 1, 1, 1, nil, nil, "").
		WillReturnRows(rows)

	// Call CreateAgent with valid agent data
	agent := &Agents{
		FirstName:              "John",
		LastName:               "Doe",
		AgentEmail:             "johndoe@example.com",
		Credentials:            AgentLoginCredentials{},
		Phone:                  "1234567890",
		RoleID:                 Role{ID: 1},
		Team:                   Teams{ID: 1},
		Unit:                   Unit{ID: 1},
		SupervisorID:           uint(1),
		ResetPasswordRequestID: uint(1),
	}
	err := model.CreateAgent(agent)
	//createdAgent, err := model.CreateAgent(agent)

	// Assert expected database interactions and returned agent data
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedAgentID, createdAgent.ID)
}

func TestCreateAgent_Success2(t *testing.T) {

	cdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer cdb.Close()

	// Initialize the GORM.io database instance
	dsn := "root:1T$hutt!ers@tcp(localhost:3306)/mydb?parseTime=true" // Change this to your actual database connection string
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	model := NewAgentDBModel(db)
	var role Role
	role.ID = 1

	// Mock expected database interaction for successful agent creation
	expectedAgentID := uint(1)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedAgentID)
	mock.ExpectQuery(`INSERT INTO agents`).
		WithArgs("John", "Doe", "johndoe@example.com", "1234567890").
		WillReturnRows(rows)

	// Call CreateAgent with valid agent data
	agent := &Agents{
		FirstName:  "John",
		LastName:   "Doe",
		AgentEmail: "johndoe@example.com",
		//Credentials: AgentLoginCredentials{}, // or appropriate value
		Phone: "1234567890",
	}
	createdAgent, err := model.CreateAgent(agent)
	if err != nil {
		t.Fatalf("unexpected error creating agent: %v", err)
	}

	// Assert expected database interactions and returned agent data
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	mock.ExpectExec(`INSERT INTO agents`).WillReturnResult(sqlmock.NewResult(1, 1))

	assert.NoError(t, err)
	assert.Equal(t, "John", createdAgent.FirstName)
	assert.Equal(t, "johndoe@example.com", createdAgent.AgentEmail)
	assert.Equal(t, expectedAgentID, createdAgent.ID)

}

func TestCreateAgent_Failure(t *testing.T) {
	dbSQL, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbSQL.Close()

	// Use gorm.Open with the *sql.DB object directly:
	db, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{Conn: dbSQL}}, &gorm.Config{})
	if err != nil {
		e := fmt.Errorf("couldnt connect to gorm db")

		fmt.Println(e.Error())
	}

	model := NewAgentDBModel(db)

	// Mock expected database error during agent creation
	mock.ExpectQuery(`INSERT INTO agents`).
		WillReturnError(fmt.Errorf("database error"))

	// Call CreateAgent and assert error handling
	agent := &Agents{} // Fill in with valid agent data
	_, err = model.CreateAgent(agent)

	assert.Error(t, err)
}
