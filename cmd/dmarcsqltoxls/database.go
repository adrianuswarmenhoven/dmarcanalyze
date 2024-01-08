package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	// Database drivers
	_ "github.com/go-sql-driver/mysql" // mysql
	_ "github.com/jackc/pgx/v5"        // postgres
	_ "modernc.org/sqlite"             // sqlite3
)

type database struct {
	backendDB *sql.DB
	driver    string

	preparedStatements map[string]*sql.Stmt
}

var (
	db database
)

var (
	preparedStatements = map[string]string{
		// Fetch metadata
		"fetch metadata": `
		SELECT * FROM metadata;
		`,
		// Fetch policy published
		"fetch policy published": `
		SELECT * FROM policy_published;
		`,
		// Fetch record
		"fetch record": `
		SELECT * FROM record;
		`,
	}
)

const (
	maxParams = 14 // Maximum number of parameters in a prepared statement, used to replace $1, $2, etc. with ? (for sqlite and mysql)
)

func initDB() error {
	err := db.Open(Configuration.Database.Driver, Configuration.Database.ConnectionString)
	if err != nil {
		return err
	}
	return nil
}

/*
insert full record variables ) as transaction
*/
func (db *database) Open(driver, connectionstring string) error {
	var err error
	db.backendDB, err = sql.Open(driver, connectionstring)
	if err != nil {
		slog.Error("error opening database", "error", err)
		return err
	}
	err = db.backendDB.Ping()
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		return err
	}
	db.driver = driver

	db.preparedStatements = make(map[string]*sql.Stmt)
	err = db.initStatements()
	if err != nil {
		slog.Error("error initializing statements", "error", err)
		return err
	}
	return nil
}

func (db *database) initStatements() error {
	for name, query := range preparedStatements {
		// Need to modify query placeholders
		switch db.driver {
		case "sqlite", "mysql":
			for cnt := 0; cnt <= maxParams; cnt++ {
				query = strings.ReplaceAll(query, fmt.Sprintf("$%dd", cnt), "?")
			}
		case "postgres": //do nothing
		default:
			slog.Error("prepared statements not implemented for driver", "driver", db.driver)
			return fmt.Errorf("prepared statements not implemented for driver %s", db.driver)
		}

		stmt, err := db.backendDB.Prepare(query)
		if err != nil {
			slog.Error("error preparing statement", "error", err)
			return err
		}
		db.preparedStatements[name] = stmt
	}
	return nil
}

func (db *database) Close() error {
	return db.backendDB.Close()
}

type Metadata struct {
	OrgName          string
	Email            string
	ExtraContactInfo string
	ReportID         string
	Begin            int64
	End              int64
}

func (db *database) FetchMetadata() ([]*Metadata, error) {
	rows, err := db.preparedStatements["fetch metadata"].Query()
	if err != nil {
		slog.Error("error querying metadata", "error", err)
		return nil, err
	}
	defer rows.Close()
	metadata := make([]*Metadata, 0)
	id := 0
	for rows.Next() {
		m := Metadata{}
		if err := rows.Scan(&id, &m.OrgName, &m.Email, &m.ExtraContactInfo, &m.ReportID, &m.Begin, &m.End); err != nil {
			slog.Error("error scanning metadata", "error", err)
			return nil, err
		}
		metadata = append(metadata, &m)
	}
	if err := rows.Err(); err != nil {
		slog.Error("error scanning metadata", "error", err)
		return nil, err
	}
	return metadata, nil
}

// PolicyPublished represents feedback>policy_published section
type PolicyPublished struct {
	Domain     string `xml:"domain"`
	ADKIM      string `xml:"adkim"`
	ASPF       string `xml:"aspf"`
	Policy     string `xml:"p"`
	SPolicy    string `xml:"sp"`
	Percentage int    `xml:"pct"`
	ReportID   string
}

func (db *database) FetchPolicyPublished() ([]*PolicyPublished, error) {
	rows, err := db.preparedStatements["fetch policy published"].Query()
	if err != nil {
		slog.Error("error querying policy published", "error", err)
		return nil, err
	}
	defer rows.Close()
	policyPublished := make([]*PolicyPublished, 0)
	id := 0
	for rows.Next() {
		p := PolicyPublished{}
		if err := rows.Scan(&id, &p.Domain, &p.ADKIM, &p.ASPF, &p.Policy, &p.SPolicy, &p.Percentage, &p.ReportID); err != nil {
			slog.Error("error scanning policy published", "error", err)
			return nil, err
		}
		policyPublished = append(policyPublished, &p)
	}
	if err := rows.Err(); err != nil {
		slog.Error("error scanning policy published", "error", err)
		return nil, err
	}
	return policyPublished, nil
}

type Record struct {
	SourceIP               string
	Count                  int
	Disposition            string
	DKIM                   string
	SPF                    string
	HeaderFrom             string
	DKIMAuthResultDomain   string
	DKIMAuthResultResult   string
	DKIMAuthResultSelector string
	SPFAuthResultDomain    string
	SPFAuthResultResult    string
	SPFAuthResultScope     string
	ReportID               string
}

func (db *database) FetchRecords() ([]*Record, error) {
	rows, err := db.preparedStatements["fetch record"].Query()
	if err != nil {
		slog.Error("error querying record", "error", err)
		return nil, err
	}
	defer rows.Close()
	records := make([]*Record, 0)
	id := 0
	for rows.Next() {
		r := Record{}
		if err := rows.Scan(&id, &r.SourceIP, &r.Count, &r.Disposition, &r.DKIM, &r.SPF, &r.HeaderFrom, &r.DKIMAuthResultDomain, &r.DKIMAuthResultResult, &r.DKIMAuthResultSelector, &r.SPFAuthResultDomain, &r.SPFAuthResultResult, &r.SPFAuthResultScope, &r.ReportID); err != nil {
			slog.Error("error scanning record", "error", err)
			return nil, err
		}
		records = append(records, &r)
	}
	if err := rows.Err(); err != nil {
		slog.Error("error scanning record", "error", err)
		return nil, err
	}
	return records, nil
}
