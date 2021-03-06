package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	dynamodbStore "github.com/fallncrlss/dictionary-app-backend/src/pkg/store/dynamodb"
)

func createWordTable(ctx context.Context, logger echo.Logger, db *dynamodb.Client) error {
	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("language"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("language"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("Word"),
	}
	if err := dynamodbStore.CreateTable(ctx, logger, db, params); err != nil {
		return errors.Wrap(err, "createWordTable failed")
	}

	return nil
}

func runDynamoMigrations(ctx context.Context, logger echo.Logger, db *dynamodb.Client) error {
	logger.Debug("Running Migrations...")

	if err := createWordTable(ctx, logger, db); err != nil {
		return err
	}

	logger.Debug("Migrations ran successfully!")

	return nil
}
