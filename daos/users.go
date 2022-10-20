package daos

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log"
	"time"
)

type User struct {
	UserID    string    `dynamodbav:"userID"`
	Login     string    `dynamodbav:"login"`
	Email     string    `dynamodbav:"email"`
	CreatedAt time.Time `dynamodbav:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt"`
}

type UsersDAO interface {
	CreateUser(ctx context.Context, login, email string) (User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type usersDAOImpl struct {
	DynamoDBClient *dynamodb.Client

	tableName string
}

func NewUsersDAO(dynamoDBClient *dynamodb.Client) UsersDAO {
	return &usersDAOImpl{
		DynamoDBClient: dynamoDBClient,
		tableName:      "users",
	}
}

func (dao usersDAOImpl) CreateUser(ctx context.Context, login, email string) (User, error) {
	newUser := User{
		UserID:    uuid.New().String(),
		Login:     login,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	putItem, err := attributevalue.MarshalMap(newUser)
	if err != nil {
		return User{}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: &dao.tableName,
		Item:      putItem,
	}

	_, err = dao.DynamoDBClient.PutItem(ctx, input)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (dao usersDAOImpl) GetUserByID(ctx context.Context, id string) (*User, error) {
	getItemOutput, err := dao.DynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"userID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(dao.tableName),
	})
	if err != nil {
		return &User{}, err
	}

	if getItemOutput.Item == nil {
		return nil, nil
	}

	user := &User{}
	err = attributevalue.UnmarshalMap(getItemOutput.Item, user)
	if err != nil {
		log.Panicf("unmarshal failed, %v", err)
	}

	return user, nil
}
