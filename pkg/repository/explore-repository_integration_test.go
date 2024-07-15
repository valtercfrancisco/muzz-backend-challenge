package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"muzz-backend-challenge/pkg/repository"
)

func setupDB(t *testing.T) (*sql.DB, func()) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println("Connection String:", connStr)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "failed to connect to database")

	// Wait for a while to ensure PostgreSQL is up
	time.Sleep(5 * time.Second)
	err = db.Ping()
	require.NoError(t, err, "failed to ping database")
	fmt.Println("Connected to the database")

	// Create tables if they do not exist
	createTables := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS likes (
			id SERIAL PRIMARY KEY,
			actor_user_id UUID NOT NULL,
			recipient_user_id UUID NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(actor_user_id, recipient_user_id),
			FOREIGN KEY (actor_user_id) REFERENCES users(user_id),
			FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS decisions (
			actor_user_id UUID NOT NULL,
			recipient_user_id UUID NOT NULL,
			liked_recipient BOOLEAN NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(actor_user_id, recipient_user_id),
			FOREIGN KEY (actor_user_id) REFERENCES users(user_id),
			FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
		);`,
	}

	for _, query := range createTables {
		_, err := db.Exec(query)
		require.NoError(t, err, "failed to create table")
	}

	closeFunc := func() {
		db.Close()
	}

	return db, closeFunc
}

func seedTestData(t *testing.T, db *sql.DB, recipientUserID, user1ID, user2ID uuid.UUID) {
	insertUsers(t, db, recipientUserID, user1ID, user2ID)

	likeInsertQueries := []string{
		fmt.Sprintf("INSERT INTO likes (actor_user_id, recipient_user_id, created_at) VALUES ('%s', '%s', NOW()) ON CONFLICT DO NOTHING;", user1ID, recipientUserID),
		fmt.Sprintf("INSERT INTO likes (actor_user_id, recipient_user_id, created_at) VALUES ('%s', '%s', NOW()) ON CONFLICT DO NOTHING;", user2ID, recipientUserID),
		fmt.Sprintf("INSERT INTO likes (actor_user_id, recipient_user_id, created_at) VALUES ('%s', '%s', NOW()) ON CONFLICT DO NOTHING;", recipientUserID, user1ID),
	}

	for _, query := range likeInsertQueries {
		fmt.Printf("Inserting like with query: %s\n", query)
		_, err := db.Exec(query)
		require.NoError(t, err, "failed to insert like")
	}
	fmt.Println("Test data seeded successfully")
}

func insertUsers(t *testing.T, db *sql.DB, userIDs ...uuid.UUID) {
	for _, userID := range userIDs {
		query := fmt.Sprintf("INSERT INTO users (user_id, username, created_at) VALUES ('%s', 'user-%s', NOW()) ON CONFLICT (user_id) DO NOTHING;", userID, userID)
		fmt.Printf("Inserting user with query: %s\n", query)
		_, err := db.Exec(query)
		require.NoError(t, err, "failed to insert user")
	}
	fmt.Println("Users inserted successfully")
}

func cleanupTestData(t *testing.T, db *sql.DB, userIDs ...uuid.UUID) {
	for _, userID := range userIDs {
		queries := []string{
			fmt.Sprintf("DELETE FROM likes WHERE actor_user_id = '%s' OR recipient_user_id = '%s';", userID, userID),
			fmt.Sprintf("DELETE FROM decisions WHERE actor_user_id = '%s' OR recipient_user_id = '%s';", userID, userID),
			fmt.Sprintf("DELETE FROM users WHERE user_id = '%s';", userID),
		}

		for _, query := range queries {
			fmt.Printf("Cleaning up data with query: %s\n", query)
			_, err := db.Exec(query)
			require.NoError(t, err, "failed to cleanup data")
		}
	}
	fmt.Println("Test data cleaned up successfully")
}

func TestIntegrationGetLikedYou(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	recipientUserID := uuid.New()
	user1ID := uuid.New()
	user2ID := uuid.New()

	cleanupTestData(t, db, recipientUserID, user1ID)
	defer cleanupTestData(t, db, recipientUserID, user1ID)
	seedTestData(t, db, recipientUserID, user1ID, user2ID)

	repo := repository.NewExploreRepository(db)
	likers, err := repo.GetLikedYou(ctx, recipientUserID.String(), 10, 0)
	require.NoError(t, err)
	assert.Len(t, likers, 2)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestIntegrationGetNewLikedYou(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	recipientUserID := uuid.New()
	user1ID := uuid.New()
	user2ID := uuid.New()

	cleanupTestData(t, db, recipientUserID, user1ID)
	defer cleanupTestData(t, db, recipientUserID, user1ID)
	seedTestData(t, db, recipientUserID, user1ID, user2ID)

	repo := repository.NewExploreRepository(db)
	likers, err := repo.GetNewLikedYou(ctx, recipientUserID.String(), 10, 0)
	require.NoError(t, err)
	assert.Len(t, likers, 1) // user2 has not liked test-recipient back

	err = tx.Commit()
	require.NoError(t, err)
}

func TestIntegrationCountLikes(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	recipientUserID := uuid.New()
	user1ID := uuid.New()
	user2ID := uuid.New()

	cleanupTestData(t, db, recipientUserID, user1ID)
	defer cleanupTestData(t, db, recipientUserID, user1ID)
	seedTestData(t, db, recipientUserID, user1ID, user2ID)

	repo := repository.NewExploreRepository(db)
	count, err := repo.CountLikes(ctx, recipientUserID.String())
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestIntegrationInsertDecision(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	actorUserID := uuid.New()
	recipientUserID := uuid.New()

	cleanupTestData(t, db, recipientUserID, actorUserID)
	defer cleanupTestData(t, db, recipientUserID, actorUserID)

	// Insert users
	insertUsers(t, db, actorUserID, recipientUserID)
	likedRecipient := true

	repo := repository.NewExploreRepository(db)
	err = repo.InsertDecision(ctx, tx, actorUserID.String(), recipientUserID.String(), likedRecipient)
	require.NoError(t, err)

	// Verify insertion
	var liked bool
	query := "SELECT liked_recipient FROM decisions WHERE actor_user_id = $1 AND recipient_user_id = $2"
	err = tx.QueryRowContext(ctx, query, actorUserID, recipientUserID).Scan(&liked)
	require.NoError(t, err)
	assert.Equal(t, likedRecipient, liked)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestIntegrationInsertLike(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	actorUserID := uuid.New()
	recipientUserID := uuid.New()

	cleanupTestData(t, db, recipientUserID, actorUserID)
	defer cleanupTestData(t, db, recipientUserID, actorUserID)

	// Insert users
	insertUsers(t, db, actorUserID, recipientUserID)

	repo := repository.NewExploreRepository(db)
	err = repo.InsertLike(ctx, tx, actorUserID.String(), recipientUserID.String())
	require.NoError(t, err)

	// Verify insertion
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM likes WHERE actor_user_id = $1 AND recipient_user_id = $2)"
	err = tx.QueryRowContext(ctx, query, actorUserID, recipientUserID).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestIntegrationDeleteLike(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	recipientUserID := uuid.New()
	user1ID := uuid.New()

	cleanupTestData(t, db, recipientUserID, user1ID)
	defer cleanupTestData(t, db, recipientUserID, user1ID)
	seedTestData(t, db, recipientUserID, user1ID, uuid.New())

	repo := repository.NewExploreRepository(db)
	err = repo.DeleteLike(ctx, tx, user1ID.String(), recipientUserID.String())
	require.NoError(t, err)

	// Verify deletion
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM likes WHERE actor_user_id = $1 AND recipient_user_id = $2)"
	err = tx.QueryRowContext(ctx, query, user1ID, recipientUserID).Scan(&exists)
	require.NoError(t, err)
	assert.False(t, exists)

	err = tx.Commit()
	require.NoError(t, err)
}

func TestIntegrationCheckMutualLike(t *testing.T) {
	db, closeFunc := setupDB(t)
	defer closeFunc()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	recipientUserID := uuid.New()
	user1ID := uuid.New()

	cleanupTestData(t, db, recipientUserID, user1ID)
	defer cleanupTestData(t, db, recipientUserID, user1ID)
	seedTestData(t, db, recipientUserID, user1ID, uuid.New())

	repo := repository.NewExploreRepository(db)
	exists, err := repo.CheckMutualLike(ctx, tx, recipientUserID.String(), user1ID.String())
	require.NoError(t, err)
	assert.True(t, exists)

	err = tx.Commit()
	require.NoError(t, err)
}
