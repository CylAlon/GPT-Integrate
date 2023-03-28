package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Mysqls struct {
	DB *sql.DB
}

type User struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	DingUserName string `json:"dingUserName"`
	DingUserId   string `json:"dingUserId"`
}
type Context struct {
	Id       int    `json:"id"`
	Uid      int    `json:"uid"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func init() {
	// 初始化DB
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// defer DB.Close()
}
func NewSql() *Mysqls {
	return &Mysqls{DB: db}
}

// 修改用户信息
func (m *Mysqls) UpdateUser(user User) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("update user set name=?,dingUserName=?,dingUserId=? where id=?", user.Name, user.DingUserName, user.DingUserId, user.Id)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// 添加用户
func (m *Mysqls) AddUser(user User) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("insert into user(name,dingUserName,dingUserId) values(?,?,?)", user.Name, user.DingUserName, user.DingUserId)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil

}
func (m *Mysqls) AddUserCheckDingId(user User) {
	// 如果用户存在则更新，不存在则添加
	_, err := m.GetUserByDingId(user.DingUserId)
	if err != nil {
		m.AddUser(user)
	}
}

// 删除用户
func (m *Mysqls) DeleteUser(id int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from user where id=?", id)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// 通过dingid删除用户
func (m *Mysqls) DeleteUserByDingId(dingId string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from user where dingUserId=?", dingId)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// 通过id获取用户信息
func (m *Mysqls) GetUser(id int) (User, error) {
	var user User
	err := m.DB.QueryRow("select * from user where id=?", id).Scan(&user.Id, &user.Name, &user.DingUserName, &user.DingUserId)
	return user, err
}

// 通过钉钉id获取用户信息
func (m *Mysqls) GetUserByDingId(dingId string) (User, error) {
	var user User
	err := m.DB.QueryRow("select * from user where dingUserId=?", dingId).Scan(&user.Id, &user.Name, &user.DingUserName, &user.DingUserId)
	return user, err
}

// 添加上下文
func (m *Mysqls) AddContext(context Context) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("insert into context(uid,question,answer) values(?,?,?)", context.Uid, context.Question, context.Answer)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// 删除uid下的所有上下文
func (m *Mysqls) DeleteContextByUid(uid int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from context where uid=?", uid)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (m *Mysqls) GetContextByUidAndSize(uid int, size int) ([]Context, error) {
	var contexts []Context = nil
	var length int = 0
	rows, err := m.DB.Query("select * from context where uid=? order by id desc", uid)
	if err != nil {
		return contexts, err
	}
	defer rows.Close()
	for rows.Next() {
		var context Context
		err = rows.Scan(&context.Id, &context.Uid, &context.Question, &context.Answer)
		if err != nil {
			return contexts, err
		}
		length += len(context.Question) + len(context.Answer)
		if length > size {
			break
		}
		contexts = append([]Context{context}, contexts...)
	}

	return contexts, nil
}
