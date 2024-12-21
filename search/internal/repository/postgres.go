package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/postgres"
)

const create string = `
CREATE TABLE IF NOT EXISTS Locations
(
	id uuid primary key NOT NULL,
	country varchar NOT NULL,
	city varchar NOT NULL,
	district varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS Finds
(
	id uuid primary key NOT NULL,
	name varchar NOT NULL,
	description varchar NOT NULL,
	type varchar NOT NULL,
	user_login varchar NOT NULL,
	location_id uuid references Locations(id) ON DELETE CASCADE
);
`

type PostgresRepo struct {
	db *sqlx.DB
}

type dest struct {
	Id         string `db:"id"`
	LocationId string `db:"location_id"`
	AddReq
	Country  string `db:"country"`
	City     string `db:"city"`
	District string `db:"district"`
}

func NewPostgresRepo(cfg postgres.PostgresConfig) (*PostgresRepo, error) {
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresRepo: %w", err)
	}
	if _, err := db.Exec(create); err != nil {
		return nil, fmt.Errorf("NewPostrgresRepo: %w", err)
	}
	return &PostgresRepo{
		db: db,
	}, nil
}

func (r *PostgresRepo) AddFind(req AddReq) (AddResp, error) {
	var (
		row      *sql.Row
		lastUUID string
	)

	const (
		insertLocation string = `
	INSERT INTO Locations VALUES ($1, $2, $3, $4) RETURNING id;
	`
		insertFind string = `
	INSERT INTO Finds VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`
	)

	tx, err := r.db.Beginx()
	if err != nil {
		return AddResp{}, fmt.Errorf("PostgresRepo AddFind: %w", err)
	}
	defer tx.Rollback()

	row = tx.QueryRow(
		sqlx.Rebind(sqlx.DOLLAR, insertLocation),
		uuid.New().String(),
		req.Location.Country,
		req.Location.City,
		req.Location.District,
	)
	if err := row.Scan(&lastUUID); err != nil {
		return AddResp{}, fmt.Errorf("PostgresRepo AddFind: %w", err)
	}

	row = tx.QueryRow(
		sqlx.Rebind(sqlx.DOLLAR, insertFind),
		uuid.New().String(),
		req.Name,
		req.Description,
		req.Type,
		req.UserLogin,
		lastUUID,
	)
	if err := row.Scan(&lastUUID); err != nil {
		return AddResp{}, fmt.Errorf("PostgresRepo AddFind: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return AddResp{}, fmt.Errorf("PostgresRepo AddFind: %w", err)
	}

	return AddResp{
		FindUUID: lastUUID,
	}, nil
}

func (r *PostgresRepo) RespondToFind(req RespondReq) (RespondResp, error) {
	var d []dest = []dest{}

	const (
		deleteFind string = `
	DELETE FROM finds WHERE id = $1;
	`
		selectFind string = `
	SELECT user_login FROM finds WHERE id = $1;
	`
	)
	tx, err := r.db.Begin()
	if err != nil {
		return RespondResp{}, fmt.Errorf("PostgresRepo RespondToFind: %w", err)
	}
	defer tx.Rollback()
	if err := r.db.Select(&d, sqlx.Rebind(sqlx.DOLLAR, selectFind), req.FindUUID); err != nil {
		return RespondResp{}, fmt.Errorf("PostgresRepo RespondToFind: %w", err)
	}
	if len(d) <= 0 {
		return RespondResp{}, fmt.Errorf("PostgresRepo RespondToFind: no rows found")
	}

	// if _, err := tx.Exec(sqlx.Rebind(sqlx.DOLLAR, deleteFind), req.FindUUID); err != nil {
	// 	return RespondResp{}, fmt.Errorf("PostgresRepo RespondToFind: %w", err)
	// }
	// if err := tx.Commit(); err != nil {
	// 	return RespondResp{}, fmt.Errorf("PostgresRepo RespondToFind: %w", err)
	// }
	return RespondResp{UserLogin: d[0].UserLogin}, nil
}

func (r *PostgresRepo) GetFind(req GetReq) (GetResp, error) {
	var (
		// row   *sql.Row
		dests []dest   = []dest{}
		finds []AddReq = []AddReq{}
	)

	// conds := make([]string, 0, 3)

	var selectFind string = fmt.Sprintf(
		`
	SELECT * FROM finds LEFT JOIN locations ON finds.location_id = locations.id %s;
	`,
		getCondition(req),
	)
	if err := r.db.Select(&dests, selectFind); err != nil {
		return GetResp{}, fmt.Errorf("PostgresRepo GetFind: %w", err)
	}

	for _, v := range dests {
		r := v.AddReq
		r.Location = location{Country: v.Country, City: v.City, District: v.District}
		finds = append(finds, r)
	}
	return GetResp{Finds: finds}, nil
}

func getCondition(req GetReq) string {
	conds := make([]string, 0, 3)
	if req.Name != "" {
		conds = append(conds, fmt.Sprintf("finds.name LIKE '%s%%'", req.Name))
	}
	if req.Type != "" {
		conds = append(conds, fmt.Sprintf("finds.type LIKE '%s%%'", req.Type))
	}
	if req.Location.Country != "" {
		conds = append(conds, fmt.Sprintf("locations.country LIKE '%s%%' AND locations.city LIKE '%s%%' AND locations.district LIKE '%s%%'", req.Location.Country, req.Location.City, req.Location.District))
	}
	if len(conds) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conds, " AND "))
}
