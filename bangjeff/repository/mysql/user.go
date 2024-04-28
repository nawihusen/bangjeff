package mysql

import (
	"bangjeff/domain"
	"context"
	"database/sql"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

// NewMySQLUserRepository is constructor of MySQL repository
func NewMySQLUserRepository(Conn *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{Conn}
}

func (db *mysqlUserRepository) SignUp(ctx context.Context, user domain.User) (err error) {
	query := `INSERT INTO user (username, password, name, phone, email, address, dtm_crt) VALUES (?, ?, ?, ?, ?, ?, NOW())`
	log.Debug("Query : " + query)

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = stmt.ExecContext(ctx, user.Username, user.Password, user.Name, user.Phone, user.Email, user.Address)

	return err
}

func (db *mysqlUserRepository) GetUsers(ctx context.Context, opt domain.Options) (users []domain.User, err error) {
	offset := (opt.Page - 1) * opt.Limit

	query := `SELECT id, username, password, name, phone, email, address, dtm_crt FROM user`
	var params []interface{}
	if opt.Name != "" || opt.Username != "" {
		query += ` WHERE `
	}

	if opt.Name != "" {
		query += ` name LIKE "%` + opt.Name + `%" `
	}

	if opt.Username != "" {
		if opt.Name != "" {
			query += ` AND `
		}

		query += ` username = ? `
		params = append(params, opt.Username)
	}

	query += ` ORDER BY ` + opt.Sort

	query += ` ` + opt.Order

	query += " LIMIT " + strconv.Itoa(int(opt.Limit))

	query += " OFFSET " + strconv.Itoa(int(offset))

	log.Debug(query)

	rows, err := db.Conn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i domain.User
		err := rows.Scan(&i.ID, &i.Username, &i.Password, &i.Name, &i.Phone, &i.Email, &i.Address, &i.DtmCrt)
		if err != nil {
			return nil, err
		}

		users = append(users, i)
	}

	return users, err
}

func (db *mysqlUserRepository) CountUsers(ctx context.Context, opt domain.Options) (total int64, err error) {
	query := `SELECT COUNT(id) FROM user`
	var params []interface{}
	if opt.Name != "" || opt.Username != "" {
		query += ` WHERE `
	}

	if opt.Name != "" {
		query += ` name LIKE "%` + opt.Name + `%" `
	}

	if opt.Username != "" {
		if opt.Name != "" {
			query += ` AND `
		}

		query += ` username = ? `
		params = append(params, opt.Username)
	}

	log.Debug(query)

	row := db.Conn.QueryRowContext(ctx, query, params...)
	err = row.Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("there is no data")
		}
		return total, err
	}

	return total, err
}

func (db *mysqlUserRepository) GetUser(ctx context.Context, id int64) (user domain.User, err error) {
	query := `SELECT id, username, password, name, phone, email, address, dtm_crt FROM user WHERE id = ? `
	log.Debug(query)

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Address,
		&user.DtmCrt,
	)
	if err != nil {
		err = errors.New("id not found")
		return
	}

	return user, err
}

func (db *mysqlUserRepository) GetByToken(ctx context.Context, token string) (user domain.User, active int64, err error) {
	query := `SELECT id, username, password, name, phone, email, address, dtm_crt, active FROM user WHERE token = ? `
	log.Debug(query)

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, token)
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Address,
		&user.DtmCrt,
		&active,
	)
	if err != nil {
		err = errors.New("token not found")
		return
	}

	return user, active, err
}

func (db *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (user domain.User, err error) {
	query := `SELECT id, username, password, name, phone, email, address, dtm_crt FROM user WHERE username = ? `
	log.Debug(query)

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, username)
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Address,
		&user.DtmCrt,
	)
	if err != nil {
		err = errors.New("username not found")
		return
	}

	return user, err
}

func (db *mysqlUserRepository) SaveToken(ctx context.Context, id int64, token domain.Token, active int64) (err error) {
	query := `UPDATE user SET token = ?, active = ? WHERE id = ?`
	log.Debug("Query : " + query)

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, token.Token, active, id)
	if err != nil {
		return
	}

	return err
}
