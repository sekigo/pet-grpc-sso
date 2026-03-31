package migrator

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var storagePath, migrationPath, migrationTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to destination file")
	flag.StringVar(&migrationPath, "storage-path", "", "Path to migration file")
	flag.StringVar(&migrationTable, "storage-path", "", "Path to migration table")
	flag.Parse()

	if storagePath == "" {
		panic("not such storage path")

	}
	if migrationPath == "" {
		panic("not such migration path")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf("slite3://%s?x-migrations=table=%s", storagePath, migrationPath),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully!")

}
