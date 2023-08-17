package dockertest_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "timescale/timescaledb-postgis",
		Tag:        "latest-pg13",
		Env: []string{
			"POSTGRES_DB=mauzi",
			"POSTGRES_USER=mauzi",
			"POSTGRES_PASSWORD=mauzipw",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds
	hostAndPort := resource.GetPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://mauzi:mauzipw@localhost:%s/mauzi?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url", databaseUrl)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	log.Println("Applying migrations")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mi, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	mi.Up()

	log.Println("Running tests")

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSomething(t *testing.T) {
	var count int
	sql_query := `
		SELECT COUNT(table_name)
		FROM information_schema.tables 
		WHERE 
			table_schema LIKE 'public' AND 
			table_type LIKE 'BASE TABLE' AND
			table_name = 'trecord';
	`

	rows, err := db.Query(sql_query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&count)
	}

	if count != 1 {
		t.Errorf("Expected 1 table, got %d", count)
	}
}
