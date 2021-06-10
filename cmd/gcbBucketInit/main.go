package main

import (
	"context"
	"github.com/BRBussy/toolios/internal/pkg/logs"
	"github.com/BRBussy/toolios/pkg/storage"
	"github.com/rs/zerolog/log"
	"time"
)

const bucketName = "sol-ar-world"

func main() {
	logs.Setup(nil)

	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelFunc()

	if _, err := storage.NewGCSProvider(
		ctx,
		bucketName,
		map[string]bool{
			"http://localhost:3000": true,
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error constructing storage provider")
	}
}
