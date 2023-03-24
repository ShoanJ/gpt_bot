package user

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gpt_bot/biz/constant"
	"gpt_bot/biz/model"
	"gpt_bot/biz/mongo_cli"
	"gpt_bot/biz/utils"
	"time"
)

const (
	// RemainTime 最早保存半小时内的上下文
	RemainTime = 1800
	// ChatContentsMaxLen 最多保存的对话条数
	ChatContentsMaxLen = 3
)

func AppendChatsToChatCtx(ctx context.Context, chats []*model.ChatContent, fromUser bool) (*model.UserProfile, error) {
	userProfile, err := FindUserProfileByUserID(ctx)
	if err != nil {
		return nil, err
	}
	if userProfile == nil {
		userProfile = &model.UserProfile{}
		err = InsertUserProfile(ctx)
		if err != nil {
			return nil, err
		}
	}
	userProfile = userProfile.AppendChatContext(chats, ChatContentsMaxLen, time.Now().Unix()-RemainTime)
	if fromUser {
		now := time.Now().Unix()
		if userProfile.LastChatTime != nil && now < *userProfile.LastChatTime+30 {
			return nil, constant.ChatNotReplyYetError
		}
		userProfile.LastChatTime = utils.Int64Ptr(now)
	} else {
		userProfile.LastChatTime = nil
	}

	updated, err := UpdateUserProfileByUserID(ctx, userProfile)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, constant.ConcurrencyWriteDBError
	}
	return userProfile, nil
}

func FindUserProfileByUserID(ctx context.Context) (*model.UserProfile, error) {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		logrus.Errorf("GetUserIDFromCtx err: %s", err.Error())
	}
	filter := bson.M{
		"UserID": userID,
	}

	result := mongo_cli.Client.Collection(model.UserProfileCollectionName).FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		logrus.Errorf("FindOne err: %s", result.Err().Error())
		return nil, result.Err()
	}

	userProfile := &model.UserProfile{}
	err = result.Decode(userProfile)
	if err != nil {
		logrus.Errorf("Decode err: %s", err.Error())
		return nil, err
	}
	return userProfile, nil
}

func UpdateUserProfileByUserID(ctx context.Context, userProfile *model.UserProfile) (bool, error) {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		logrus.Errorf("GetUserIDFromCtx err: %s", err.Error())
	}
	filter := bson.M{
		"UserID":  userID,
		"Version": userProfile.Version,
	}
	result, err := mongo_cli.Client.Collection(model.UserProfileCollectionName).UpdateOne(ctx, filter, userProfile.ToUpdateBsonD())
	if err != nil {
		logrus.Errorf("UpdateOne err: %s", err.Error())
		return false, err
	}
	return result.MatchedCount > 0, nil
}

func InsertUserProfile(ctx context.Context) error {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		logrus.Errorf("GetUserIDFromCtx err: %s", err.Error())
	}
	_, err = mongo_cli.Client.Collection(model.UserProfileCollectionName).InsertOne(ctx, bson.D{
		{"UserID", userID},
		{"Version", 0},
	})
	if err != nil {
		logrus.Errorf("InsertOne err: %s", err.Error())
	}
	return err
}
