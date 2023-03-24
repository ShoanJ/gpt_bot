package lark

import (
	"context"
	"encoding/json"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gpt_bot/biz/conf"
	"gpt_bot/biz/utils"
)

var (
	larkCli = lark.NewClient(conf.GetConf().Lark.AppID, conf.GetConf().Lark.AppSecret)
)

func ReplyMsg(ctx context.Context, openID, msgID, content string) error {
	reqBuilder := larkim.NewReplyMessageReqBuilder()
	reqBuilder.MessageId(msgID)
	contentJson, _ := json.Marshal(map[string]string{
		"text": content,
	})
	reqBuilder.Body(&larkim.ReplyMessageReqBody{
		Content: utils.StringPtr(string(contentJson)),
		MsgType: utils.StringPtr("text"),
	})

	resp, err := larkCli.Im.Message.Reply(ctx, reqBuilder.Build())
	if err != nil {
		return err
	}
	if resp.StatusCode != 0 && resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d, body: %s", resp.StatusCode, string(resp.RawBody))
	}
	return nil
}
