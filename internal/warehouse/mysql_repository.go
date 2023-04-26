package warehouse

import (
	"database/sql"
	"errors"
	"fmt"

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

type Repository interface {
	GetByID(id int) (domain.Warehouse, error)
	Create(p domain.Warehouse) (domain.Warehouse, error)
	GetAll() ([]domain.Warehouse, error)
	ReportProducts(id int) (domain.ReportProducts, error)
}

// mySQLRepository struct definition
type mySQLRepository struct {
	database *sql.DB
}

// NewMySQLRepository constructor function
func NewMySQLRepository(database *sql.DB) Repository {
	return &mySQLRepository{database}
}

// Create method to insert a new product into the products table
func (repository *mySQLRepository) Create(warehouse domain.Warehouse) (domain.Warehouse, error) {

	statement, err := repository.database.Prepare(`INSERT INTO warehouses(name, address, telephone, capacity) VALUES( ?, ?, ?, ?)`)
	if err != nil {
		return domain.Warehouse{}, err
	}
	defer statement.Close()
	var result sql.Result
	result, err = statement.Exec(warehouse.Name, warehouse.Address, warehouse.Telephone, warehouse.Capacity)
	if err != nil {
		fmt.Println(err)
		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return domain.Warehouse{}, ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return domain.Warehouse{}, ErrAccessDenied
		case 1046:
			return domain.Warehouse{}, ErrNoDatabaseSelected
		case 1048:
			return domain.Warehouse{}, ErrNotNullColumn
		case 1054:
			return domain.Warehouse{}, ErrUnknownColumn
		case 1062:
			return domain.Warehouse{}, ErrDuplicateEntry
		case 1064:
			return domain.Warehouse{}, ErrSyntaxError
		case 1146:
			return domain.Warehouse{}, ErrTableDoesNotExist
		default:
			return domain.Warehouse{}, ErrInternal
		}
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		return domain.Warehouse{}, err
	}
	warehouse.Id = int(insertedId)
	return warehouse, nil
}

func (repository *mySQLRepository) GetByID(id int) (warehouse domain.Warehouse, err error) {
	query := `SELECT id, name, address, telephone, capacity FROM products where id = ?`
	row := repository.database.QueryRow(query, id)
	err = row.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Address, &warehouse.Telephone, &warehouse.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Warehouse{}, ErrNotFound
		}

		mysqlError, ok := err.(*mysql.MySQLError)
		if !ok {
			return domain.Warehouse{}, ErrInternal
		}
		switch mysqlError.Number {
		case 1044, 1045:
			return domain.Warehouse{}, ErrAccessDenied
		case 1046:
			return domain.Warehouse{}, ErrNoDatabaseSelected
		case 1054:
			return domain.Warehouse{}, ErrUnknownColumn
		case 1064:
			return domain.Warehouse{}, ErrSyntaxError
		case 1146:
			return domain.Warehouse{}, ErrTableDoesNotExist
		default:
			return domain.Warehouse{}, ErrInternal
		}
	}

	return warehouse, nil
}

func (repository *mySQLRepository) GetAll() ([]domain.Warehouse, error) {
	query := (`SELECT id, name, address, telephone, capacity FROM warehouses`)
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

	var warehouses []domain.Warehouse
	for rows.Next() {
		var warehouse domain.Warehouse
		if err := rows.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Address, &warehouse.Telephone, &warehouse.Capacity); err != nil {
			return nil, ErrInternal
		}
		warehouses = append(warehouses, warehouse)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return warehouses, nil
}

func (repository *mySQLRepository) ReportProducts(id int) (reportProducts domain.ReportProducts, err error) {
	query := `SELECT w.name, count(*) FROM warehouses w 
	LEFT JOIN products p
	ON w.id = p.id_warehouse
	WHERE w.id = ?
	GROUP BY w.id`
	row := repository.database.QueryRow(query, id)
	err = row.Scan(&reportProducts.WarehouseName, &reportProducts.ProductCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return reportProducts, ErrNotFound
		}
		return reportProducts, ErrInternal
	}
	return reportProducts, nil
}
