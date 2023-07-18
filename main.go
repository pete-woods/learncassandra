package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Clean exit")
}

func run() error {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}

	// connect to the cluster
	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}
	defer session.Close()

	ctx := context.Background()

	// insert a tweet
	err = session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("error inserting tweet: %w", err)
	}

	var id gocql.UUID
	var text string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	err = session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
		"me").WithContext(ctx).Consistency(gocql.One).Scan(&id, &text)
	if err != nil {
		return fmt.Errorf("error querying: %w", err)
	}
	fmt.Printf("Tweet: %s %s\n", id, text)

	fmt.Println("Finding tweets:")
	// list all tweets
	scanner := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`,
		"me").WithContext(ctx).Iter().Scanner()
	for scanner.Next() {
		err = scanner.Scan(&id, &text)
		if err != nil {
			return fmt.Errorf("error scanning: %w", err)
		}
		fmt.Printf("Tweet: %s %s\n", id, text)
	}
	// scanner.Err() closes the iterator, so scanner nor iter should be used afterwards.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error closing scanner: %w", err)
	}

	return nil
}
