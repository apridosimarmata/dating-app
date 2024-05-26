package activity

import (
	"context"
	"dating-app/domain/activity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type activityRepository struct {
	client   *mongo.Client
	activity *mongo.Collection
	likes    *mongo.Collection
}

func NewActivityRepository(client *mongo.Client) activity.ActivityRepository {
	return &activityRepository{
		activity: client.Database("dating-app").Collection("activity"),
		likes:    client.Database("dating-app").Collection("likes"),
	}
}

func (repository *activityRepository) GetUserActivityCountByDate(ctx context.Context, userUid string, date string) (res int, err error) {

	activityKey := fmt.Sprintf("activities.%s", date)
	filter := bson.M{
		"user_uid":  userUid,
		activityKey: 1,
	}
	result := repository.activity.FindOne(ctx, filter, nil)
	if err = result.Err(); err != nil {
		return 0, err
	}

	var _res map[string]interface{}

	if err := result.Decode(_res); err != nil {
		return 0, err
	}

	return len(_res), nil
}

func (repository *activityRepository) GetUserActivitiesByDate(ctx context.Context, userUid string, date string) (res map[string]string, err error) {

	activityKey := fmt.Sprintf("activities.%s", date)
	filter := bson.M{
		"user_uid":  userUid,
		activityKey: 1,
	}
	result := repository.activity.FindOne(ctx, filter, nil)
	if err = result.Err(); err != nil {
		return nil, err
	}

	var _res map[string]interface{}
	if err := result.Decode(_res); err != nil {
		return nil, err
	}

	return res, nil
}

func (repository *activityRepository) UpdateUserActivitiesByDate(ctx context.Context, userUid string, updatedActivities map[string]string, like map[string]interface{}, date string) (err error) {
	var session mongo.Session
	if session, err = repository.activity.Database().Client().StartSession(); err != nil {
		return err
	}

	if err = session.StartTransaction(); err != nil {
		return err
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		activityKey := fmt.Sprintf("activities.%s", date)
		filter := bson.M{
			"user_uid": userUid,
		}
		result := repository.activity.FindOneAndUpdate(
			ctx,
			filter,
			bson.D{
				{
					Key: "$set",
					Value: bson.M{
						activityKey: updatedActivities,
					},
				},
			},
		)
		if err = result.Err(); err != nil {
			return err
		}

		if like != nil {
			_, err := repository.likes.InsertOne(ctx, like)
			if err != nil {
				return err
			}
		}

		/*** commit ***/
		if err = session.CommitTransaction(sc); err != nil {
			return err
		}
		/*** end ***/

		return nil
	}); err != nil {
		return err
	}

	return nil
}
