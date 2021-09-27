package repository

import (
	"context"
	"database/sql"
	"github.com/INEFFABLE-games/positionsUserService/internal/models"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUserRepository_Insert(t *testing.T) {
	conn, err := sql.Open("postgres", "port=5432 host=localhost user=postgres password=12345 dbname=positions sslmode=disable")
	if err != nil {
		log.Errorf("unable to connect with postgres %v", err)
	}

	r := NewUserRepository(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = r.Insert(ctx, models.User{
		Uid:      "23635utdykdtyuk4j",
		Login:    "AnthonyBest",
		Password: "12345",
	})

	require.Nil(t, err)
}

func TestUserRepository_UpdateBalance(t *testing.T) {
	conn, err := sql.Open("postgres", "port=5432 host=localhost user=postgres password=12345 dbname=positions sslmode=disable")
	if err != nil {
		log.Errorf("unable to connect with postgres %v", err)
	}

	r := NewUserRepository(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = r.UpdateBalance(ctx, "23635utdykdtyuk4j", 1000)

	require.Nil(t, err)
}

func TestUserRepository_GetBalance(t *testing.T) {
	conn, err := sql.Open("postgres", "port=5432 host=localhost user=postgres password=12345 dbname=positions sslmode=disable")
	if err != nil {
		log.Errorf("unable to connect with postgres %v", err)
	}

	r := NewUserRepository(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	balance, err := r.GetBalance(ctx, "23635utdykdtyuk4j")

	require.Nil(t, err)
	require.Equal(t, int64(1000), balance)
}
