package main

import (
  "fmt"
  "database/sql"
  // "errors"
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

type transaction struct {
  ID   int    `json:"id"`
  Date string `json:"date"`
  Hour string `json:"hour"`
  Type string `json:"type"` 
  Bitcoins float64 `json:"bitcoins"`
  Convert_Tx float64 `json:"convert_tx"`
  Final_Value float64 `json:"final_value"`
  User_ID_1 int `json:"user_id_1"`
  User_ID_2 int `json:"user_id_2"`
}

func (u *user) getUser(db *sql.DB) error {
  statement := fmt.Sprintf("SELECT name, last_name, email, password, birthday, balance FROM users WHERE id=%d", u.ID)
  return db.QueryRow(statement).Scan(&u.Name, &u.Last_Name, &u.Email, &u.Password, &u.Birthday, &u.Balance)
}

func (u *transaction) getTransaction(db *sql.DB) error {
  statement := fmt.Sprintf("SELECT date, hour, type, bitcoins, convert_tx, final_value, user_id_1, user_id_2 FROM transactions WHERE id=%d", u.ID)
  return db.QueryRow(statement).Scan(&u.Date, &u.Hour, &u.Type, &u.Bitcoins, &u.Convert_Tx, &u.Final_Value, &u.User_ID_1, &u.User_ID_2)
}

func (u *user) updateUser(db *sql.DB) error {
  statement := fmt.Sprintf("UPDATE users SET name='%s', last_name='%s', email='%s', password='%s', birthday='%s', balance=%f WHERE id=%d", u.Name, u.Last_Name, u.Email, u.Password, u.Birthday, u.Balance, u.ID)
  _, err := db.Exec(statement)
  return err
}

func (u *transaction) updateTransaction(db *sql.DB) error {
  statement := fmt.Sprintf("UPDATE transactions SET date='%s', hour='%s', type='%s', bitcoins=%f, convert_tx=%f, final_value=%f, user_id_1=%d, user_id_2=%d WHERE id=%d", u.Date, u.Hour, u.Type, u.Bitcoins, u.Convert_Tx, u.Final_Value, u.User_ID_1, u.User_ID_2, u.ID)
  _, err := db.Exec(statement)
  return err
}

func (u *user) deleteUser(db *sql.DB) error {
  statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
  _, err := db.Exec(statement)
  return err
}

func (u *transaction) deleteTransaction(db *sql.DB) error {
  statement := fmt.Sprintf("DELETE FROM transactions WHERE id=%d", u.ID)
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

func (u *transaction) createTransaction(db *sql.DB) error {
  statement := fmt.Sprintf("INSERT INTO transactions(date, hour, type, bitcoins, convert_tx, final_value, user_id_1, user_id_2) VALUES('%s', '%s', '%s', %f, %f, %f, %d, %d)", u.Date, u.Hour, u.Type, u.Bitcoins, u.Convert_Tx, u.Final_Value, u.User_ID_1, u.User_ID_2)
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

func getTransactions(db *sql.DB, start, count int) ([]transaction, error) {
  statement := fmt.Sprintf("SELECT id, date, hour, type, bitcoins, convert_tx, final_value, user_id_1, user_id_2 FROM transactions LIMIT %d OFFSET %d", count, start)
  rows, err := db.Query(statement)

  if err != nil {
    return nil, err
  }

  defer rows.Close()

  transactions := []transaction{}

  for rows.Next() {
    var u transaction
    if err := rows.Scan(&u.ID, &u.Date, &u.Hour, &u.Type, &u.Bitcoins, &u.Convert_Tx, &u.Final_Value, &u.User_ID_1, &u.User_ID_2); err != nil {
      return nil, err
    }
    transactions = append(transactions, u)
  }

  return transactions, nil
}

// func (u *transaction) getTransaction(db *sql.DB) error {
//   statement := fmt.Sprintf("SELECT date, hour, type, bitcoins, convert_tx, final_value, user_id_1, user_id_2 FROM users WHERE id=%d", u.ID)
//   return db.QueryRow(statement).Scan(&u.Date, &u.Hour, &u.Type, &u.Bitcoins, &u.Convert_Tx, &u.Final_Value, &u.User_ID_1, &u.User_ID_2)
// }

// func (u *transaction) updateTransaction(db *sql.DB) error {
//   // statement := fmt.Sprintf("UPDATE users SET date='%s', last_name='%s', email='%s', password='%s', birthday='%s', balance=%f WHERE id=%d", u.Name, u.Last_Name, u.Email, u.Password, u.Birthday, u.Balance, u.ID)
//   // _, err := db.Exec(statement)
//   // return err
//   return errors.New("Not implemented")
// }

// func (u *transaction) deleteTransaction(db *sql.DB) error {
//     return errors.New("Not implemented")
// }
// func (u *transaction) createTransaction(db *sql.DB) error {
//     return errors.New("Not implemented")
// }
// func getTransactions(db *sql.DB, start, count int) ([]transaction, error) {
//     return nil, errors.New("Not implemented")
// }