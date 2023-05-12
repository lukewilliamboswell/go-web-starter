package user

import (
	"database/sql"
	"fmt"
	"log"
)

type UserMSSQL struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {

	// Create the user table if it doesn't exist
	createUserTable(db)

	return &UserMSSQL{db}
}

func createUserTable(db *sql.DB) {
	const query = `
	IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'users')
	BEGIN
		CREATE TABLE users (
			id INT IDENTITY(1,1) PRIMARY KEY,
			principal_id NVARCHAR(36) UNIQUE NOT NULL,
			principal_name NVARCHAR(255) NOT NULL,
			principal_provider NVARCHAR(255) NOT NULL,
			access_level INT NOT NULL DEFAULT 0
		);	
	END
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Print("failed to create user table: %w", err)
	}
}

func (repo *UserMSSQL) GetUsers() ([]User, error) {
	const query = `
	SELECT 
		principal_id, 
		principal_name, 
		principal_provider,
		access_level
	FROM
		users
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users = make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.PrincipalId,
			&user.PrincipalName,
			&user.PrincipalProvider,
			&user.AccessLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user data: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to process user rows: %w", err)
	}

	return users, nil
}

func (repo *UserMSSQL) GetUser(principalID string) (User, error) {

	var user User

	// Look up the user in the database
	const query = `
	SELECT 
		principal_id, 
		principal_name, 
		principal_provider,
		access_level
	FROM 
		users 
	WHERE 
		principal_id = @principal_id
	`
	err := repo.db.QueryRow(
		query,
		sql.Named("principal_id", principalID),
	).Scan(
		&user.PrincipalId,
		&user.PrincipalName,
		&user.PrincipalProvider,
		&user.AccessLevel,
	)
	if err == sql.ErrNoRows {
		// No user found
		return User{}, nil
	}
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user for principal ID %s: %w", principalID, err)
	}

	// Return the user
	return user, nil

}

// InsertUser inserts a new user into the database
func (repo *UserMSSQL) InsertUser(user User) error {
	const query = `
	IF NOT EXISTS (SELECT 1 FROM users WHERE principal_id = @principal_id)
	BEGIN
		INSERT INTO users (principal_id, principal_name, principal_provider)
		VALUES (@principal_id, @principal_name, @principal_provider)
	END
	`
	_, err := repo.db.Exec(
		query,
		sql.Named("principal_id", user.PrincipalId),
		sql.Named("principal_name", user.PrincipalName),
		sql.Named("principal_provider", user.PrincipalProvider),
	)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
