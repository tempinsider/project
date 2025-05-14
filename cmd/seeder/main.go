package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	query := `
	INSERT INTO messages (id, created_at, updated_at, message_content, phone_number, status)
SELECT 
    uuid_generate_v4() as id,
    NOW() as created_at,
    NOW() as updated_at,
    md5(random()::text) as message_content,
    '+90' || (5000000000 + floor(random() * 999999999)::bigint)::text as phone_number,
    0 as status
FROM generate_series(1, 1000);
`

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	// check if we can actually connect to it.
	var data any
	err = conn.QueryRow(ctx, "select version();").Scan(&data)
	if err != nil {
		panic(err)
	}

	tag, err := conn.Exec(ctx, query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v rows added\n", tag.RowsAffected())
}
