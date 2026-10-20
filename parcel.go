package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (created_at, client, status, address) VALUES (:created_at, :client, :status, :address)",
		sql.Named("created_at", p.CreatedAt),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row, err := s.db.Query("SELECT number, client, address, status, created_at FROM parcel WHERE number = :number", sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}

	defer row.Close()
	p := Parcel{}

	if !row.Next() {
		return Parcel{}, sql.ErrNoRows
	}

	err = row.Scan(&p.Number, &p.Client, &p.Address, &p.Status, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	row, err := s.db.Query("SELECT number, client, address, status, created_at FROM parcel WHERE client = :client", sql.Named("client", client))

	if err != nil {
		fmt.Println(err)
		return []Parcel{}, err
	}

	defer row.Close()

	var res []Parcel

	for row.Next() {
		p := Parcel{}
		err := row.Scan(&p.Number, &p.Client, &p.Address, &p.Status, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return []Parcel{}, err
		}

		res = append(res, p)

	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("number", number),
		sql.Named("status", status))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = 'registered'", sql.Named("address", address), sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = 'registered'", sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
