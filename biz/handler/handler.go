package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"gpt_bot/biz/conf"
	"gpt_bot/biz/constant"
	"gpt_bot/biz/gpt"
	"gpt_bot/biz/lark"
	"gpt_bot/biz/model"
	"gpt_bot/biz/user"
	"gpt_bot/biz/utils"
	"runtime/debug"
	"time"
)

var (
	eventDispatcher = dispatcher.NewEventDispatcher(conf.GetConf().Lark.VerificationToken, conf.GetConf().Lark.EncryptKey)
)

func init() {
	eventDispatcher.OnP2MessageReceiveV1(msgRecvEvent)
}

func ReceiveLarkEvent(ctx *gin.Context) {
	eventHandler := httpserverext.NewEventHandlerFunc(eventDispatcher)
	eventHandler(ctx.Writer, ctx.Request)
}

func ReceiveLarkEventFacade(ctx *gin.Context) {
	eventHandler := httpserverext.NewEventHandlerFunc(eventDispatcher)
	eventHandler(ctx.Writer, ctx.Request)
}

func msgRecvEvent(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	if event == nil || event.Event == nil || event.Event.Sender == nil || event.Event.Sender.SenderId == nil ||
		event.Event.Sender.SenderId.OpenId == nil || event.Event.Message == nil || event.Event.Message.Content == nil ||
		event.Event.Message.MessageId == nil {
		eventBytes, _ := json.Marshal(event)
		logrus.Errorf("something is nil, event: %s", string(eventBytes))
		return nil
	}
	ctx = utils.SetUserIDToCtx(ctx, *event.Event.Sender.SenderId.OpenId)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logrus.Errorf("err: %v \nstack: %s", e, string(debug.Stack()))
			}
		}()
		reply, err := getReplyMsg(ctx, event)
		if err != nil {
			return
		}
		err = lark.ReplyMsg(ctx, *event.Event.Sender.SenderId.OpenId, *event.Event.Message.MessageId, reply)
		if err != nil {
			logrus.Errorf("lark.ReplyMsg err: %s", err.Error())
			return
		}
		logrus.Infof("reply success")
	}()

	return nil
}

func getReplyMsg(ctx context.Context, event *larkim.P2MessageReceiveV1) (string, error) {
	chatContent := &model.ChatContent{
		Role:      "user",
		Content:   *event.Event.Message.Content,
		CreatedAt: time.Now().Unix(),
	}
	userProfile, err := user.AppendChatsToChatCtx(ctx, []*model.ChatContent{
		chatContent,
	}, true)
	if err != nil {
		return err2Msg(err)
	}

	replyMsg, err := gpt.Chat(ctx, *event.Event.Sender.SenderId.OpenId, userProfile.ChatContext)
	if err != nil {
		logrus.Errorf("gpt.Chat err: %s", err.Error())
		return "", err
	}

	replyContent := &model.ChatContent{
		Role:      replyMsg.Role,
		Content:   replyMsg.Content,
		CreatedAt: time.Now().Unix(),
	}
	_, err = user.AppendChatsToChatCtx(ctx, []*model.ChatContent{
		replyContent,
	}, false)
	if err != nil {
		return err2Msg(err)
	}

	return replyMsg.Content, err
}

func err2Msg(e error) (string, error) {
	switch e {
	case constant.ChatNotReplyYetError, constant.ConcurrencyWriteDBError:
		return e.Error(), nil
	}
	return "", e
}
