package main

import (
  "fmt"
  "database/sql"
)

type user struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
  Last_Name string `json:"last_name"`
  Email string `json:"email"`
  Password string `json:"password"`
  Birthday string `json:"birthday"`
  Balance  float64    `json:"balance"`
}

func (u *user) getUser(db *sql.DB) error {
  statement := fmt.Sprintf("SELECT name, last_name, email, password, birthday, balance FROM users WHERE id=%d", u.ID)
  return db.QueryRow(statement).Scan(&u.Name, &u.Last_Name, &u.Email, &u.Password, &u.Birthday, &u.Balance)
}

func (u *user) updateUser(db *sql.DB) error {
  statement := fmt.Sprintf("UPDATE users SET name='%s', last_name='%s', email='%s', password='%s', birthday='%s', balance=%f WHERE id=%d", u.Name, u.Last_Name, u.Email, u.Password, u.Birthday, u.Balance, u.ID)
  _, err := db.Exec(statement)
  return err
}

func (u *user) deleteUser(db *sql.DB) error {
  statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
  _, err := db.Exec(statement)
  return err
}

func (u *user) createUser(db *sql.DB) error {
  statement := fmt.Sprintf("INSERT INTO users(name, last_name, email, password, birthday, balance) VALUES('%s', '%s', '%s', '%s', '%s', %f)", u.Name, u.Last_Name, u.Email, u.Password, u.Birthday, u.Balance)
  _, err := db.Exec(statement)

  if err != nil {
    return err
  }

  err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

  if err != nil {
    return err
  }

  return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
  statement := fmt.Sprintf("SELECT id, name, last_name, email, password, birthday, balance FROM users LIMIT %d OFFSET %d", count, start)
  rows, err := db.Query(statement)

  if err != nil {
    return nil, err
  }

  defer rows.Close()

  users := []user{}

  for rows.Next() {
    var u user
    if err := rows.Scan(&u.ID, &u.Name, &u.Last_Name, &u.Email, &u.Password, &u.Birthday, &u.Balance); err != nil {
      return nil, err
    }
    users = append(users, u)
  }

  return users, nil
}