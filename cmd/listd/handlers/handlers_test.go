package handlers

import (
	"os"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/configuration"
	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/item"
	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// testSuite is a struct type that contains necessary fields to carry out
// tasks to fully test the handlers package along with it's integrations.
type testSuite struct {
	a     *Application
	lists []list.List
	items []item.Item
}

// reseedDatabase is a function attached to the testSuite type that attempts
// to reseed the database back to its original testing state.
func (ts *testSuite) reseedDatabase(t *testing.T) {
	var err error

	if err = testdb.Truncate(ts.a.db); err != nil {
		t.Errorf("error truncating database: %v", err)
	}

	if ts.lists, err = testdb.SeedLists(ts.a.db); err != nil {
		t.Errorf("error seeding lists: %v", err)
	}

	if ts.items, err = testdb.SeedItems(ts.a.db, ts.lists); err != nil {
		t.Errorf("error seeding items: %v", err)
	}
}

// ts is the global variable that is of type testSuite which helps test the
// entirety of the handlers package and it's integrations.
var ts testSuite

// TestMain handles the setup of the testSuite, runs all of the unit tests within
// the handlers package, and cleans up afterward.
func TestMain(m *testing.M) {
	var err error
	var dbc *sqlx.DB

	exitCode := 1

	defer func() {
		if err != nil {
			log.WithError(err).Info("error in handlers TestMain")
		}

		if dbc != nil {
			if err = dbc.Close(); err != nil {
				log.WithError(err).Info("close test database connection")
			}
		}

		os.Exit(exitCode)
	}()

	if dbc, err = db.NewConnection(&configuration.Config{
		DBUser: configuration.DefaultDBUser,
		DBPass: configuration.DefaultDBPass,
		DBName: testdb.DatabaseName,
		DBHost: configuration.DefaultDBHost,
		DBPort: configuration.DefaultDBPort,
	}); err != nil {
		err = errors.Wrap(err, "create test database connection")
		return
	}

	ts.a = NewApplication(dbc, &configuration.Config{})

	// Initial seeding of the test database using test values defined within
	// the testdb package.
	if err = testdb.Truncate(ts.a.db); err != nil {
		err = errors.Wrap(err, "truncate database")
		return
	}

	if ts.lists, err = testdb.SeedLists(ts.a.db); err != nil {
		err = errors.Wrap(err, "seed lists")
		return
	}

	if ts.items, err = testdb.SeedItems(ts.a.db, ts.lists); err != nil {
		err = errors.Wrap(err, "seed items")
		return
	}

	exitCode = m.Run()
}
