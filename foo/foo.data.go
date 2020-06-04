package foo

import (
	"database/sql"

	"github.com/pluralsight/inventoryservice/database"
)

func getFoo(productID int) (*Foo, error) {
	row := database.DbConn.QueryRow(
		`SELECT 
			productId, 
			message, 
			age, 
			name, 
			surname 
		FROM foos
		WHERE productId = ?`,
		productID,
	)

	foo := &Foo{}

	err := row.Scan(
		&foo.ProductID,
		&foo.Message,
		&foo.Age,
		&foo.Name,
		&foo.Surname,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return foo, nil
}

func removeFoo(productID int) error {
	_, err := database.DbConn.Query(
		`DELETE FROM foos WHERE productId = ?`, productID,
	)

	if err != nil {
		return err
	}

	return nil
}

func getFooList() ([]Foo, error) {
	results, err := database.DbConn.Query(
		`SELECT 
			productId, 
			message, 
			age, 
			name, 
			surname 
		from foos`,
	)

	if err != nil {
		return nil, err
	}

	defer results.Close()

	foos := make([]Foo, 0)

	for results.Next() {
		var foo Foo

		results.Scan(
			&foo.ProductID,
			&foo.Message,
			&foo.Age,
			&foo.Name,
			&foo.Surname,
		)

		foos = append(foos, foo)
	}

	return foos, nil
}

func updateFoo(foo Foo) error {
	_, err := database.DbConn.Exec(
		`UPDATE foos SET 
			message = ?, 
			age = ?, 
			name = ?, 
			surname = ?
		WHERE productId = ? `,
		foo.Message,
		foo.Age,
		foo.Name,
		foo.Surname,
	)

	if err != nil {
		return err
	}

	return nil
}

func insertFoo(foo Foo) (int, error) {
	result, err := database.DbConn.Exec(
		`INSERT INTO foos 
			(
				message, 
				age, 
				name, 
				surname
			) 
		VALUES (?, ?, ?, ?)`,
		foo.Message,
		foo.Age,
		foo.Name,
		foo.Surname,
	)

	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(insertID), err
}
