package product

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bootcamp-go/consignas-go-db.git/internal/domain"
	"github.com/go-sql-driver/mysql"
)

// Define custom error variables
var (
	ErrAccessDenied       = errors.New("access denied for user")
	ErrNoDatabaseSelected = errors.New("no database selected")
	ErrNotNullColumn      = errors.New("column cannot be null")
	ErrUnknownColumn      = errors.New("unknown column")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrSyntaxError        = errors.New("syntax error")
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrParsingDate        = errors.New("error parsing date")
)

// mySQLRepository struct definition
type mySQLRepository struct {
	database *sql.DB
}

// NewMySQLRepository constructor function
func NewMySQLRepository(database *sql.DB) Repository {
	return &mySQLRepository{database}
}

// Create method to insert a new product into the products table
func (repository *mySQLRepository) Create(product domain.Product) (domain.Product, error) {
	parsedDate, err := time.Parse("02/01/2006", product.Expiration)
	if err != nil {
		return domain.Product{}, ErrParsingDate
	}
	formattedDate := parsedDate.Format("2006-01-02")

	statement, err := repository.database.Prepare(`INSERT INTO products(name, quantity, code_value, is_published, expiration, price, id_warehouse) VALUES( ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return domain.Product{}, err
	}
	defer statement.Close()
	var result sql.Result
	result, err = statement.Exec(product.Name, product.Quantity, product.CodeValue, product.IsPublished, formattedDate, product.Price, product.WarehouseId)
	if err != nil {
		fmt.Println(err)
		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return domain.Product{}, ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return domain.Product{}, ErrAccessDenied
		case 1046:
			return domain.Product{}, ErrNoDatabaseSelected
		case 1048:
			return domain.Product{}, ErrNotNullColumn
		case 1054:
			return domain.Product{}, ErrUnknownColumn
		case 1062:
			return domain.Product{}, ErrDuplicateEntry
		case 1064:
			return domain.Product{}, ErrSyntaxError
		case 1146:
			return domain.Product{}, ErrTableDoesNotExist
		default:
			return domain.Product{}, ErrInternal
		}
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		return domain.Product{}, err
	}
	product.Id = int(insertedId)
	return product, nil
}

func (repository *mySQLRepository) GetAll() ([]domain.Product, error) {
	query := (`SELECT id, name, quantity, code_value, is_published, expiration, price, id_warehouse FROM products`)
	rows, err := repository.database.Query(query)
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return nil, ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return nil, ErrAccessDenied
		case 1046:
			return nil, ErrNoDatabaseSelected
		case 1054:
			return nil, ErrUnknownColumn
		case 1064:
			return nil, ErrSyntaxError
		case 1146:
			return nil, ErrTableDoesNotExist
		default:
			return nil, ErrInternal
		}
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price, &product.WarehouseId); err != nil {
			return nil, ErrInternal
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return products, nil
}

func (repository *mySQLRepository) GetByID(id int) (product domain.Product, err error) {
	query := `SELECT id, name, quantity, code_value, is_published, expiration, price, id_warehouse FROM products where id = ?`
	row := repository.database.QueryRow(query, id)
	err = row.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price, &product.WarehouseId)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return domain.Product{}, ErrNotFound
		}

		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return domain.Product{}, ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return domain.Product{}, ErrAccessDenied
		case 1046:
			return domain.Product{}, ErrNoDatabaseSelected
		case 1054:
			return domain.Product{}, ErrUnknownColumn
		case 1064:
			return domain.Product{}, ErrSyntaxError
		case 1146:
			return domain.Product{}, ErrTableDoesNotExist
		default:
			return domain.Product{}, ErrInternal
		}
	}
	fmt.Println(product.Name)

	return product, nil
}

func (repository *mySQLRepository) Update(id int, product domain.Product) (domain.Product, error) {
	parsedDate, err := time.Parse("02/01/2006", product.Expiration)
	if err != nil {
		return domain.Product{}, ErrParsingDate
	}
	formattedDate := parsedDate.Format("2006-01-02")
	statement, err := repository.database.Prepare(`UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ?, id_warehouse = ? WHERE id = ?`)
	if err != nil {
		return domain.Product{}, ErrInternal
	}
	defer statement.Close()
	result, err := statement.Exec(product.Name, product.Quantity, product.CodeValue, product.IsPublished, formattedDate, product.Price, product.WarehouseId, id)

	if err != nil {
		fmt.Println(err)
		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return domain.Product{}, ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return domain.Product{}, ErrAccessDenied
		case 1046:
			return domain.Product{}, ErrNoDatabaseSelected
		case 1048:
			return domain.Product{}, ErrNotNullColumn
		case 1054:
			return domain.Product{}, ErrUnknownColumn
		case 1062:
			return domain.Product{}, ErrDuplicateEntry
		case 1064:
			return domain.Product{}, ErrSyntaxError
		case 1146:
			return domain.Product{}, ErrTableDoesNotExist
		default:
			return domain.Product{}, ErrInternal
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.Product{}, ErrInternal
	}
	if rowsAffected == 0 && err != nil {
		return domain.Product{}, ErrNotFound
	}
	product.Id = id
	return product, nil
}

func (repository *mySQLRepository) Delete(id int) error {
	statement, err := repository.database.Prepare(`DELETE FROM products WHERE id = ?`)
	if err != nil {
		return ErrInternal
	}
	defer statement.Close()
	result, err := statement.Exec(id)

	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return ErrAccessDenied
		case 1046:
			return ErrNoDatabaseSelected
		case 1054:
			return ErrUnknownColumn
		case 1064:
			return ErrSyntaxError
		case 1146:
			return ErrTableDoesNotExist
		default:
			return ErrInternal
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrInternal
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
