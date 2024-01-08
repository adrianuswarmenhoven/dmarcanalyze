package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	// Database drivers
	_ "github.com/go-sql-driver/mysql" // mysql
	_ "github.com/jackc/pgx/v5"        // postgres
	report "github.com/oliverpool/go-dmarc-report"
	_ "modernc.org/sqlite" // sqlite3
)

type database struct {
	backendDB *sql.DB
	driver    string

	preparedStatements map[string]*sql.Stmt
}

var (
	preparedStatements = map[string]string{
		// CREATE TABLE system
		"create table system": `
		CREATE TABLE IF NOT EXISTS system (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		last_run INTEGER(8)
		);
		`,
		// CREATE TABLE metadata
		"create table metadata": `
		CREATE TABLE IF NOT EXISTS metadata (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		organization TEXT,
		email TEXT,
		extra_contact_info TEXT,
		report_id TEXT UNIQUE NOT NULL,
		begin_date INTEGER(8),
		end_date INTEGER(8)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS metadata_report_id ON metadata (report_id);
		CREATE INDEX IF NOT EXISTS metadata_organization ON metadata (organization);
		CREATE INDEX IF NOT EXISTS metadata_email ON metadata (email);
		CREATE INDEX IF NOT EXISTS metadata_begin_date ON metadata (begin_date);
		CREATE INDEX IF NOT EXISTS metadata_end_date ON metadata (end_date);
		`,
		// INSERT INTO metadata
		"insert into metadata": `
		INSERT INTO metadata (
			organization,
			email,
			extra_contact_info,
			report_id,
			begin_date,
			end_date
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		);
		`,
		// CREATE TABLE policy_published
		"create table policy_published": `
		CREATE TABLE IF NOT EXISTS policy_published (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT,
		adkim TEXT,
		aspf TEXT,
		policy TEXT,
		spolicy TEXT,
		percentage INTEGER(1),
		report_id TEXT,
		FOREIGN KEY (report_id)	REFERENCES metadata (report_id) 
		   ON UPDATE CASCADE
		   ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS policy_published_domain ON policy_published (domain);
		CREATE INDEX IF NOT EXISTS policy_published_adkim ON policy_published (adkim);
		CREATE INDEX IF NOT EXISTS policy_published_aspf ON policy_published (aspf);
		CREATE INDEX IF NOT EXISTS policy_published_policy ON policy_published (policy);
		CREATE INDEX IF NOT EXISTS policy_published_spolicy ON policy_published (spolicy);
		CREATE INDEX IF NOT EXISTS policy_published_percentage ON policy_published (percentage);
		CREATE INDEX IF NOT EXISTS policy_published_report_id ON policy_published (report_id);
		`,
		// INSERT INTO policy_published
		"insert into policy_published": `
		INSERT INTO policy_published (
			domain,
			adkim,
			aspf,
			policy,
			spolicy,
			percentage,
			report_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		);
		`,
		// CREATE TABLE record
		"create table record": `
		CREATE TABLE IF NOT EXISTS record (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source_ip TEXT,
		count int(8),
		disposition TEXT,
		dkim TEXT,
		spf TEXT,
		header_from TEXT,
		dkim_auth_result_domain TEXT,
		dkim_auth_result_result TEXT,
		dkim_auth_result_selector TEXT,	
		spf_auth_result_domain TEXT,
		spf_auth_result_result TEXT,
		spf_auth_result_scope TEXT,
		report_id TEXT,
		FOREIGN KEY (report_id)	REFERENCES metadata (report_id) 
	   		ON UPDATE CASCADE
	   		ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS record_source_ip ON record (source_ip);
		CREATE INDEX IF NOT EXISTS record_count ON record (count);
		CREATE INDEX IF NOT EXISTS record_disposition ON record (disposition);
		CREATE INDEX IF NOT EXISTS record_dkim ON record (dkim);
		CREATE INDEX IF NOT EXISTS record_spf ON record (spf);
		CREATE INDEX IF NOT EXISTS record_header_from ON record (header_from);
		CREATE INDEX IF NOT EXISTS record_dkim_auth_result_domain ON record (dkim_auth_result_domain);
		CREATE INDEX IF NOT EXISTS record_dkim_auth_result_result ON record (dkim_auth_result_result);
		CREATE INDEX IF NOT EXISTS record_dkim_auth_result_selector ON record (dkim_auth_result_selector);
		CREATE INDEX IF NOT EXISTS record_spf_auth_result_domain ON record (spf_auth_result_domain);
		CREATE INDEX IF NOT EXISTS record_spf_auth_result_result ON record (spf_auth_result_result);
		CREATE INDEX IF NOT EXISTS record_spf_auth_result_scope ON record (spf_auth_result_scope);
		CREATE INDEX IF NOT EXISTS record_report_id ON record (report_id);
		`,
		// INSERT INTO record
		"insert into record": `
		INSERT INTO record (
			source_ip,
			count,
			disposition,
			dkim,
			spf,
			header_from,
			dkim_auth_result_domain,
			dkim_auth_result_result,
			dkim_auth_result_selector,
			spf_auth_result_domain,
			spf_auth_result_result,
			spf_auth_result_scope,
			report_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
		);
		`,
	}
)

const (
	maxParams = 14 // Maximum number of parameters in a prepared statement, used to replace $1, $2, etc. with ? (for sqlite and mysql)
)

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

	err = db.initTables()
	if err != nil {
		slog.Error("error initializing tables", "error", err)
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

func (db *database) initTables() error {
	for pstmt, stmt := range db.preparedStatements {
		if strings.HasPrefix(pstmt, "create table") {
			_, err := stmt.Exec()
			if err != nil {
				slog.Error("error creating table", "statement", pstmt, "error", err)
				return err
			}
		}
	}
	return nil
}

func storeReports(reps []*report.Aggregate) error {
	db := database{}
	err := db.Open(Configuration.Database.Driver, Configuration.Database.ConnectionString)
	if err != nil {
		slog.Error("error opening database", "error", err)
		return err
	}
	defer db.Close()
	total := len(reps)
reportLoop:
	for id, report := range reps {
		if Configuration.LogProgress > 0 && (id+1)%Configuration.LogProgress == 0 {
			slog.Info("Storing reports:", "current", id+1, "total", total) //id+1 because id starts at 0
		}
		slog.Debug("storing report", "report", report.Metadata.ReportID)
		// Add metadata
		_, err := db.preparedStatements["insert into metadata"].Exec(
			report.Metadata.OrgName,
			report.Metadata.Email,
			report.Metadata.ExtraContactInfo,
			report.Metadata.ReportID,
			report.Metadata.DateRange.Begin.Unix(),
			report.Metadata.DateRange.End.Unix(),
		)
		if err != nil {
			switch {
			// NB: need to update for other database drivers
			case strings.Contains(err.Error(), "UNIQUE constraint failed: metadata.report_id"):
				slog.Debug("report already exists", "report", report.Metadata.ReportID)
				continue reportLoop
			default:
				slog.Error("error inserting metadata", "error", err)
				return err
			}
		}
		// Add policy_published
		_, err = db.preparedStatements["insert into policy_published"].Exec(
			report.PolicyPublished.Domain,
			report.PolicyPublished.ADKIM,
			report.PolicyPublished.ASPF,
			report.PolicyPublished.Policy,
			report.PolicyPublished.SPolicy,
			report.PolicyPublished.Percentage,
			report.Metadata.ReportID,
		)
		if err != nil {
			slog.Error("error inserting policy_published", "error", err)
			return err
		}
		// Add records
		for _, record := range report.Records {
			_, err = db.preparedStatements["insert into record"].Exec(
				record.Row.SourceIP,
				record.Row.Count,
				record.Row.PolicyEvaluated.Disposition,
				record.Row.PolicyEvaluated.DKIM,
				record.Row.PolicyEvaluated.SPF,
				record.Identifiers.HeaderFrom,
				record.AuthResults.DKIM.Domain,
				record.AuthResults.DKIM.Result,
				record.AuthResults.DKIM.Selector,
				record.AuthResults.SPF.Domain,
				record.AuthResults.SPF.Result,
				record.AuthResults.SPF.Scope,
				report.Metadata.ReportID,
			)
			if err != nil {
				slog.Error("error inserting record", "error", err)
				return err
			}
		}
	}

	return nil
}

func getLastRun() (time.Time, error) {
	db := database{}
	err := db.Open(Configuration.Database.Driver, Configuration.Database.ConnectionString)
	if err != nil {
		return time.Now(), err
	}
	defer db.Close()
	rows, err := db.backendDB.Query("SELECT last_run FROM system")
	if err != nil {
		return time.Now(), err
	}
	defer rows.Close()
	last_run := int64(0)
	for rows.Next() {
		if err := rows.Scan(&last_run); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	if last_run == 0 {
		return time.Now().AddDate(-100, 0, -1), nil
	}
	return time.Unix(last_run, 0), nil
}

func setLastRun(last_run time.Time) error {
	db := database{}
	err := db.Open(Configuration.Database.Driver, Configuration.Database.ConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()
	// FIXME: Variables for postgress
	_, err = db.backendDB.Exec("INSERT INTO system (last_run) VALUES (?)", last_run.Unix())
	if err != nil {
		return err
	}
	return nil
}
