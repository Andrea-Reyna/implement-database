package warehouse

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/bootcamp-go/consignas-go-db.git/internal/domain"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("TOKEN", "123456")
	txdb.Register("txdb", "mysql", "root@tcp(localhost:3306)/my_db?parseTime=true")
}

func TestRepositoryMySQL_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", "my_db")
		assert.NoError(t, err)
		defer db.Close()

		rp := NewMySQLRepository(db)

		exp := []domain.Warehouse{
			{Id: 1, Name: "Main Warehouse", Address: "221 Baker Street", Telephone: "4555666", Capacity: 100},
			{Id: 2, Name: "SuperMarket", Address: "123 Main Street", Telephone: "555-555-5555", Capacity: 2222},
		}

		// act
		wr, err := rp.GetAll()

		t.Log(wr[0])

		// assert
		assert.NoError(t, err)
		assert.Equal(t, exp, wr)
	})
}

func TestRepositoryMySQL_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", "my_db")
		assert.NoError(t, err)
		defer db.Close()

		rp := NewMySQLRepository(db)

		exp := domain.Warehouse{Id: 1, Name: "Main Warehouse", Address: "221 Baker Street", Telephone: "4555666", Capacity: 100}

		// act
		wr, err := rp.GetByID(1)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, exp, wr)
	})

	t.Run("failed, not found", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", "my_db")
		assert.NoError(t, err)
		defer db.Close()

		st := NewMySQLRepository(db)

		// act
		pr, err := st.GetByID(100)

		// assert
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Empty(t, pr)
	})
}

func TestRepositoryMySQL_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", "my_db")
		assert.NoError(t, err)
		defer db.Close()

		rp := NewMySQLRepository(db)

		warehouse := domain.Warehouse{Name: "New Warehouse", Address: "221 Baker Street", Telephone: "4555666", Capacity: 100}

		// act
		wr, err := rp.Create(warehouse)
		exp := warehouse
		exp.Id = wr.Id

		// assert
		assert.NoError(t, err)
		assert.Equal(t, exp, wr)
	})
}
