package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	// 初始化DB
	var err error
	DB, err = sql.Open("sqlite3", "./cyl_chat.db")
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	// defer DB.Close()
}

// sqlite数据库字段如下，给我写增删改查的函数，开启事物
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}
type Context struct {
	Id       int    `json:"id"`
	Uid      int    `json:"uid"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// 增加用户
func addUser(user User) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO user(name, key) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Name, user.Key)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func SqlAddUser(user User) error {
	_, err := SqlGetUserForName(user.Name)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return addUser(user)
		}
		return err
	}
	return SqlUpdateUser(user)
}

// 删除用户
func SqlDeleteUserForId(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM user WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func SqlDeleteUserForName(name string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM user WHERE name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 修改用户
func SqlUpdateUser(user User) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE user SET name=?, key=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Name, user.Key, user.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 查询用户
func SqlGetUserForid(id int) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, name, key FROM user WHERE id=?", id).Scan(&user.Id, &user.Name, &user.Key)
	if err != nil {
		return user, err
	}
	return user, nil
}
func SqlGetUserForName(name string)(User, error){
	var user User
	err := DB.QueryRow("SELECT id, name, key FROM user WHERE name=?", name).Scan(&user.Id, &user.Name, &user.Key)
	if err != nil {
		return user, err
	}
	return user, nil
}

//增加记录
func SqlAddContext(record Context) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }
    stmt, err := tx.Prepare("INSERT INTO context(uid, question, answer) values(?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(record.Uid, record.Question, record.Answer)
    if err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return nil
}

//删除记录
func SqlDeleteContextForId(id int) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }
    stmt, err := tx.Prepare("DELETE FROM context WHERE id=?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(id)
    if err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return nil
}
func SqlDeleteContextForUid(uid int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM context WHERE uid=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(uid)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//修改记录
func updateContext(record Context) error {
    tx, err := DB.Begin()
    if err != nil {
        return err
    }
    stmt, err := tx.Prepare("UPDATE context SET uid=?, question=?, answer=? WHERE id=?")
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(record.Uid, record.Question, record.Answer, record.Id)
    if err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return nil
}

//查询记录
func SqlGetContextForId(id int) (Context, error) {
    var record Context
    err := DB.QueryRow("SELECT id, uid, question, answer FROM context WHERE id=?", id).Scan(&record.Id, &record.Uid, &record.Question, &record.Answer)
    if err != nil {
        return record, err
    }
    return record, nil
}
func SqlGetContextsByUid(uid int) ([]Context, error) {
    var contexts []Context
    rows, err := DB.Query("SELECT id, uid, question, answer FROM context WHERE uid=?", uid)
    if err != nil {
        return contexts, err
    }
    defer rows.Close()
    for rows.Next() {
        var context Context
        err := rows.Scan(&context.Id, &context.Uid, &context.Question, &context.Answer)
        if err != nil {
            return contexts, err
        }
        contexts = append(contexts, context)
    }
    if err = rows.Err(); err != nil {
        return contexts, err
    }
    return contexts, nil
}

func SqlAddContextLimit(uid int, question string, answer string) ([]Context,error) {
	// 根据uid查询记录到变量ctx中
	var ctx []Context
	ctx, err := SqlGetContextsByUid(uid)
	if err != nil {
		return ctx, err
	}
	// 判断ctx的长度
	if len(ctx) == 5 {
		// 删除最开始的那一条记录
		err = SqlDeleteContextForId(ctx[0].Id)
		if err != nil {
			return ctx, err
		}
		// 去掉最开始的那一条记录
		ctx = ctx[1:]
	}
	// 插入新的记录
	var record Context
	record.Uid = uid
	record.Question = question
	record.Answer = answer
	err = SqlAddContext(record)
	if err != nil {
		return  ctx,err
	}
	// 查询新的记录
	ctx, err = SqlGetContextsByUid(uid)
	if err != nil {
		return ctx, err
	}
	return ctx,nil
}




