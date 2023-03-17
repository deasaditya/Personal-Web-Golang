package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	databaseUrl := "postgres://postgres:07210612@localhost:5432/personalweb"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Success to connect database")
}
