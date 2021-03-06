// Package models contains the types for schema 'mydb'.
package models

// GENERATED BY XO. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// DimensionCompany represents a row from 'mydb.dimension_company'.
type DimensionCompany struct {
	Company string         `json:"company"` // company
	Opaque  sql.NullString `json:"opaque"`  // opaque

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the DimensionCompany exists in the database.
func (dc *DimensionCompany) Exists() bool {
	return dc._exists
}

// Deleted provides information if the DimensionCompany has been deleted from the database.
func (dc *DimensionCompany) Deleted() bool {
	return dc._deleted
}

// Insert inserts the DimensionCompany to the database.
func (dc *DimensionCompany) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if dc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO mydb.dimension_company (` +
		`company,` +
		`opaque` +
		`) VALUES (` +
		`?,` +
		`?` +
		`)`

	// run query
	XOLog(sqlstr, dc.Opaque)
	_, err = db.Exec(sqlstr, dc.Company, dc.Opaque)
	if err != nil {
		return err
	}

	dc._exists = true

	return nil
}

// Update updates the DimensionCompany in the database.
func (dc *DimensionCompany) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !dc._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if dc._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE mydb.dimension_company SET ` +
		`opaque = ?` +
		` WHERE company = ?`

	// run query
	XOLog(sqlstr, dc.Opaque, dc.Company)
	_, err = db.Exec(sqlstr, dc.Opaque, dc.Company)
	return err
}

// Save saves the DimensionCompany to the database.
func (dc *DimensionCompany) Save(db XODB) error {
	if dc.Exists() {
		return dc.Update(db)
	}

	return dc.Insert(db)
}

// Delete deletes the DimensionCompany from the database.
func (dc *DimensionCompany) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !dc._exists {
		return nil
	}

	// if deleted, bail
	if dc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM mydb.dimension_company WHERE company = ?`

	// run query
	XOLog(sqlstr, dc.Company)
	_, err = db.Exec(sqlstr, dc.Company)
	if err != nil {
		return err
	}

	// set deleted
	dc._deleted = true

	return nil
}

// DimensionCompanyByCompany retrieves a row from 'mydb.dimension_company' as a DimensionCompany.
//
// Generated from index 'dimension_company_company_pkey'.
func DimensionCompanyByCompany(db XODB, company string) (*DimensionCompany, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`company, opaque ` +
		`FROM mydb.dimension_company ` +
		`WHERE company = ?`

	// run query
	XOLog(sqlstr, company)
	dc := DimensionCompany{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, company).Scan(&dc.Company, &dc.Opaque)
	if err != nil {
		return nil, err
	}

	return &dc, nil
}
