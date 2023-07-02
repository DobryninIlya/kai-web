package sqlstore

import (
	"database/sql"
	"main/internal/app/model"
	"main/internal/app/store"
	_ "main/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	return r.store.db.QueryRow(
		"INSERT INTO public.users(id_vk, name, groupp, distribution, admlevel, groupreal, \"dateChange\", balance, distr, warn, expiration, banhistory, ischeked, role, login, potok_lecture, has_own_shed, affiliate)"+
			"VALUES ($1, '', $2, 1, 1, $3, Now(), 0, 0, 0, '2020-01-01', 0, 0, $4, $5, true, false, false) RETURNING id_vk",
		u.ID,
		u.Group,
		u.GroupReal,
		u.Role,
		u.Login,
	).Scan(&u.ID)
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	var login sql.NullString
	if err := r.store.db.QueryRow(
		"SELECT id_vk, name, groupp, distribution, admlevel, groupreal, \"dateChange\", distr, warn, "+
			"expiration, banhistory, ischeked, role, login, potok_lecture, has_own_shed, affiliate FROM public.users"+
			" WHERE id_vk = $1",
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Group,
		&u.Distribution,
		&u.Admlevel,
		&u.GroupReal,
		&u.DateChanged,
		&u.Distr,
		&u.Warn,
		&u.Expiration,
		&u.BanHistory,
		&u.IsChecked,
		&u.Role,
		&login,
		&u.PotokLecture,
		&u.HasOwnShed,
		&u.Affiliate,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	if login.Valid {
		u.Login = login.String
	} else {
		u.Login = ""
	}

	return u, nil
}

func (r *UserRepository) MakeVerification(v *model.VerificationParams, u *model.User) error {
	var id int
	url := "INSERT INTO public.user_verification(id, date_update, faculty, course, group_id, idcard, groupname, studentid) VALUES ($1, Now(),$2, $3, $4, $5, $6, $7) RETURNING id"
	if err := r.store.db.QueryRow(
		url,
		u.ID,
		v.Faculty,
		v.Course,
		v.Group,
		v.ID,
		v.Groupname,
		v.Student,
	).Scan(&id); err != nil {
		return err
	}

	if res, _ := r.Find(u.ID); res == nil {
		err := r.Create(u)

		return err
	}
	return nil
}

// Find ...
//func (r *UserRepository) Find(id int) (*model.User, error) {
//	u := &model.User{}
//	if err := r.store.db.QueryRow(
//		"SELECT id, email, encrypted_password FROM users WHERE id = $1",
//		id,
//	).Scan(
//		&u.ID,
//		&u.Email,
//		&u.EncryptedPassword,
//	); err != nil {
//		if err == sql.ErrNoRows {
//			return nil, store.ErrRecordNotFound
//		}
//
//		return nil, err
//	}
//
//	return u, nil
//}
//
//// FindByEmail ...
//func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
//	u := &model.User{}
//	if err := r.store.db.QueryRow(
//		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
//		email,
//	).Scan(
//		&u.ID,
//		&u.Email,
//		&u.EncryptedPassword,
//	); err != nil {
//		if err == sql.ErrNoRows {
//			return nil, store.ErrRecordNotFound
//		}
//
//		return nil, err
//	}
//
//	return u, nil
//}
