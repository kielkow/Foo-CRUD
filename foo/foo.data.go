package foo

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/pluralsight/inventoryservice/database"
)

func getFoo(productID int) (*Foo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext( 
		ctx,
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
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`DELETE FROM foos WHERE productId = ?`, productID,
	)

	if err != nil {
		return err
	}

	return nil
}

func getFooList() ([]Foo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(
		ctx,
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
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
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
		foo.ProductID,
	)

	if err != nil {
		return err
	}

	return nil
}

func insertFoo(foo Foo) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(
		ctx,
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

func getToptenFoos() ([]Foo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(
		ctx,
		`SELECT 
			productId, 
			message, 
			age, 
			name, 
			surname 
		from foos ORDER BY productId DESC LIMIT 10`,
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

func searchFooData(fooFilter ReportFilter) ([]Foo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`SELECT
		productId,
		message,
		age,
		LOWER(name), 
		LOWER(surname)
		FROM foos WHERE
	`)

	if fooFilter.Name != "" {
		queryBuilder.WriteString(`name LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(fooFilter.Name)+"%")
	}

	if fooFilter.Surname != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}

		queryBuilder.WriteString(`surname LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(fooFilter.Surname)+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)

	if err != nil {
		log.Println(err.Error())
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