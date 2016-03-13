package fields

import (
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"io"
	"log"

	sql "github.com/aodin/sol"
	pg "github.com/aodin/sol/postgres"
)

// Copyright 2011 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// UUID code is a variant of code.google.com/p/go-uuid
// Added database driver and restriction to v4

var UUIDv4 = pg.UUID().NotNull()

type UUID [16]byte

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (uuid UUID) String() string {
	//b := [16]byte(uuid)
	b := uuid
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x", b[:4], b[4:6], b[6:8], b[8:10], b[10:],
	)
}

// TODO remove the nesting
func (uuid UUID) MarshalJSON() ([]byte, error) {
	if len(uuid) == 0 {
		return []byte(`""`), nil
	}
	return []byte(`"` + uuid.String() + `"`), nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == `""` {
		return nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("UUIDs must be a valid UUID version 4")
	}
	data = data[1 : len(data)-1]
	uu, err := ParseUUID(string(data))
	if err != nil {
		return err
	}
	*u = uu
	return nil
}

// Scan converts an SQL value into a UUID
func (uuid *UUID) Scan(value interface{}) error {
	uu, _ := ParseUUID(string(value.([]byte)))
	*uuid = uu
	return nil
}

// Value returns the UUID formatted for insert into SQL
func (uuid UUID) Value() (driver.Value, error) {
	return uuid.String(), nil
}

// Equals returns true if the UUIDs are equal
func (uuid UUID) Equals(other UUID) bool {
	return uuid == other
}

// Exists returns true if the UUID is a valid UUID
func (uuid UUID) Exists() bool {
	return uuid[6] >= 0x40 && uuid[6] < 0x50
}

var _ sql.Modifier = UUID{}

func (uuid UUID) Modify(table *sql.TableElem) error {
	return sql.Column("uuid", UUIDv4).Modify(table)
}

// http://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
func NewUUID() (u UUID) {
	if _, err := io.ReadFull(rand.Reader, u[:]); err != nil {
		log.Panic(err) // rand should never fail
	}
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10
	return
}

// TODO require v4
func ParseUUID(s string) (UUID, error) {
	if len(s) != 36 {
		return UUID{}, fmt.Errorf("UUIDs must have a length of 36 characters")
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return UUID{}, fmt.Errorf(
			"UUIDs must have a dash in the 8, 13, 18, and 23 positions",
		)
	}
	var uuid [16]byte

	// Convert the hex characters to bytes
	for i, x := range []int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		if v, ok := xtob(s[x:]); !ok {
			return UUID{}, fmt.Errorf(
				"UUIDs must have a valid hex encoded byte starting at position %d", x)
		} else {
			uuid[i] = v
		}
	}
	return uuid, nil
}

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = []byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob converts the the first two hex bytes of x into a byte.
func xtob(x string) (byte, bool) {
	b1 := xvalues[x[0]]
	b2 := xvalues[x[1]]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}
