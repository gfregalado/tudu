package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Product struct {
	Id        float64
	Name      string
	Price     float64
	Available bool
	createdAt time.Time
	updatedAt time.Time
}

const ProductsTableKey string = "products"

func InitProductsTable(dbClient *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    	id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        price NUMERIC(6,2) NOT NULL,
        available BOOLEAN,
        createdAt timestamp DEFAULT NOW(),
    	updatedAt timestamp DEFAULT NOW()
	)`, ProductsTableKey)

	_, err := dbClient.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (i *Product) Create(dbClient *sql.DB) (Product, error) {
	query := fmt.Sprintf(`INSERT INTO %s (name, price, available)
		VALUES ($1, $2,$3) RETURNING *`, ProductsTableKey)

	var p Product

	err := dbClient.QueryRow(query, i.Name, i.Price, i.Available).Scan(&p.Id, &p.Name, &p.Price, &p.Available, &p.createdAt, &p.updatedAt)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func (i *Product) GetById(dbClient *sql.DB, id int) (Product, error) {
	query := fmt.Sprintf(`SELECT id ,name, available, price, createAt, updateAt FROM %s where id = $1`, ProductsTableKey)

	var p Product

	err := dbClient.QueryRow(query, id).Scan(&p.Id, &p.Name, &p.Price, &p.Available, &p.createdAt, &p.updatedAt)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}
