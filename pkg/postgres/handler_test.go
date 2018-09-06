package postgres

import (
	"testing"
)

func TestShouldCreateRandomRole(t *testing.T) {
	dsn := "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable"
	pgh := NewPostgresHandler(dsn)
	err := pgh.EstablishConnection()
	if err != nil {
		t.Logf("Connecting to dsn %s failed because of %v", dsn, err)
		t.FailNow()
	}
	defer pgh.db.Close()
	roleName, err := pgh.CreateRandomRole()
	if err != nil {
		t.Logf("Test - Unable to create random role %v", err)
		t.FailNow()
	}
	defer func() {
		_, err = pgh.db.Exec("DROP ROLE %s", roleName)
		if err != nil {
			t.Logf("Test - Unable to drop the created role name %s", roleName)
		}
	}()
}

func TestShouldCreateDatabase(t *testing.T) {
	dsn := "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable"
	pgh := NewPostgresHandler(dsn)
	err := pgh.EstablishConnection()
	if err != nil {
		t.Logf("Connecting to dsn %s failed because of %v", dsn, err)
		t.FailNow()
	}
	defer pgh.db.Close()
	roleName, err := pgh.CreateRandomRole()
	if err != nil {
		t.Logf("Test - Unable to create random role %v", err)
		t.FailNow()
	}
	defer func() {
		_, err := pgh.db.Exec("DROP ROLE %s", roleName)
		if err != nil {
			t.Logf("Test - Unable to drop the created role name %s - %+v", roleName, err)
		}
	}()
	err = pgh.CreateDatabase("testdb", roleName)
	if err != nil {
		t.Logf("Test - Unable to create database for %s for %s - %+v", "testdb", roleName, err)
	}
	defer func() {
		_, err := pgh.db.Exec("drop database testdb")
		if err != nil {
			t.Logf("Test - Unable to drop the created database name testdb, %+v", err)
		}
	}()
}
