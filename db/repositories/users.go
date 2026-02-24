package db

import (
	"AuthInGo/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create() (error)
	GetById() (*models.User,error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(_db *sql.DB) *UserRepositoryImpl{
	return &UserRepositoryImpl{
		db: _db,
	}
}

func (u *UserRepositoryImpl) Create() (error) {
	query:=`Insert into users (username,email,password) values (?,?,?)`

	res,err:=u.db.Exec(query,"testUser","test@gmail.com","testPassword")

	if err!=nil{
		fmt.Println("Error Inserting user:",err)
		return err
	}

	rowsAffected,err:=res.RowsAffected()

	if err!=nil{
		fmt.Println("Error Occured:",err)
		return err
	}else if rowsAffected==0{
		fmt.Println("No rows affected, user not created")
		return nil
	}

	fmt.Println("User created successfully. Rows Affected:",rowsAffected)

	return nil
}

func (u *UserRepositoryImpl) GetById() (*models.User,error) {
	fmt.Println("Repo layer called")
	
	// Step1: Preparing Query
	query:=`Select id, username, email, password, created_at, updated_at
			From users
			Where id=?`
	
	// Step2: Execute the query
	row:=u.db.QueryRow(query,1)

	// Step3: Process the result
	user:=&models.User{}

	err:=row.Scan(&user.Id,&user.Username,&user.Email,&user.Password,&user.CreatedAt,&user.UpdatedAt)

	if err!=nil {
		if err==sql.ErrNoRows{
			fmt.Println("No rows found with given id")
			return nil,err
		}else{
			fmt.Println("Error scanning user",err)
			return nil,err
		}
	}

	fmt.Println("User found:",user)

	return user,nil
}