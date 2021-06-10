package storage

import (
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	googleStorage "cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GCSProvider struct {
	bucketName        string
	gcsClient         *googleStorage.Client
	bucketHandle      *googleStorage.BucketHandle
	credentialsClient *credentials.IamCredentialsClient
}

func NewGCSProvider(
	ctx context.Context,
	bucketName string,
	allowedDomains map[string]bool,
) (*GCSProvider, error) {
	// construct new google storage client
	client, err := googleStorage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create new google cloud storage client")
		return nil, fmt.Errorf(
			"failed to create new google cloud storage client: %w",
			err,
		)
	}

	// get a handle on the required bucket
	bucketHandle := client.Bucket(bucketName)

	// get bucket attributes to confirm connection
	bucketAttrs, err := bucketHandle.Attrs(context.Background())
	if err != nil {
		msg := fmt.Sprintf("failed to get handle to google cloud storage bucket '%s'", bucketName)
		log.Error().Err(err).Msg(msg)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}
	log.Debug().Msg(fmt.Sprintf(
		"Got a handle to google cloud storage bucket %s, created at %s, located in %s with storage class %s",
		bucketAttrs.Name, bucketAttrs.Created, bucketAttrs.Location, bucketAttrs.StorageClass,
	))

	// set up cors on the bucket

	// prepare cors origins
	corsOrigins := make([]string, 0)
	for origin := range allowedDomains {
		corsOrigins = append(
			corsOrigins,
			origin,
		)
	}

	// set cors
	if _, err := bucketHandle.Update(
		context.Background(),
		googleStorage.BucketAttrsToUpdate{
			CORS: []googleStorage.CORS{
				{
					Methods:         []string{http.MethodGet, http.MethodPut},
					Origins:         corsOrigins,
					ResponseHeaders: []string{"Content-Type"},
				},
			},
		},
	); err != nil {
		log.Error().Err(err).Msg("unable to set bucket cors policy")
		return nil, fmt.Errorf(
			"unable to set bucket cors policy: %w",
			err,
		)
	}

	// get a credentials client to be used for signing the urls
	credentialsClient, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not get IAM credentials client")
		return nil, fmt.Errorf(
			"could not get IAM credentials client: %w",
			err,
		)
	}

	return &GCSProvider{
		bucketName:        bucketName,
		gcsClient:         client,
		bucketHandle:      bucketHandle,
		credentialsClient: credentialsClient,
	}, nil
}

func (G *GCSProvider) Store(ctx context.Context, request StoreRequest) (*StoreResponse, error) {
	panic("implement me")
}
