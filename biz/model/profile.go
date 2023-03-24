package model

import (
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const UserProfileCollectionName = "UserProfile"

type UserProfile struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"UserID"`
	ApiKey       string             `bson:"ApiKey"`
	ChatContext  []*ChatContent     `bson:"ChatContext"`
	LastChatTime *int64             `bson:"LastChatTime"`
	Version      int64              `bson:"Version"`
}

func NewUserProfile() *UserProfile {
	return &UserProfile{}
}

func (p *UserProfile) CollectionName() string {
	return "UserProfile"
}

func (p *UserProfile) AppendChatContext(chats []*ChatContent, lengthLimit int, createdAt int64) *UserProfile {
	p.ChatContext = append(p.ChatContext, chats...)

	var cutIndex int
	if len(p.ChatContext) > lengthLimit {
		cutIndex = len(p.ChatContext) - lengthLimit
	}
	for ; cutIndex < len(p.ChatContext); cutIndex++ {
		if p.ChatContext[cutIndex].CreatedAt > createdAt {
			break
		}
	}
	p.ChatContext = p.ChatContext[cutIndex:]
	return p
}
func (p *UserProfile) ToUpdateBsonD() bson.D {
	return bson.D{
		{"$set", bson.D{
			{"ApiKey", p.ApiKey},
			{"ChatContext", p.ChatContext},
			{"LastChatTime", p.LastChatTime},
			{"Version", p.Version + 1},
		}},
	}
}

type ChatContent struct {
	Role      string `bson:"Role"`
	Content   string `bson:"Content"`
	CreatedAt int64  `bson:"CreatedAt"`
}

func NewChatContent(msg openai.ChatCompletionMessage) *ChatContent {
	return &ChatContent{
		Role:      msg.Role,
		Content:   msg.Content,
		CreatedAt: time.Now().Unix(),
	}
}

func (c *ChatContent) ToChatCompletionMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    c.Role,
		Content: c.Content,
	}
}
