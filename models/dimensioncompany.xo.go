// Package models contains the types for schema 'mydb'.
package models

// GENERATED BY XO. DO NOT EDIT.

// DimensionCompany represents a row from 'mydb.dimension_company'.
type DimensionCompany struct {
	Company string `json:"company"` // company
}

// DimensionCompanyByCompany retrieves a row from 'mydb.dimension_company' as a DimensionCompany.
//
// Generated from index 'dimension_company_company_pkey'.
func DimensionCompanyByCompany(db XODB, company string) (*DimensionCompany, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`company ` +
		`FROM mydb.dimension_company ` +
		`WHERE company = ?`

	// run query
	XOLog(sqlstr, company)
	dc := DimensionCompany{}

	err = db.QueryRow(sqlstr, company).Scan(&dc.Company)
	if err != nil {
		return nil, err
	}

	return &dc, nil
}
