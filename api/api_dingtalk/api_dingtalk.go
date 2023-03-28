package api_dingtalk

import (
	"GPT-Integrate/config"
	"GPT-Integrate/gpt"
	"GPT-Integrate/platform/dingtalk"
	lg "GPT-Integrate/utils/log"
	"GPT-Integrate/utils/sqlite"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	GPT_TIMEOUT = 180
)

func gpt0_3(msg string) gpt.GptCompletion {
	parm := gpt.GptCompletion{
		Key: config.GetConfigGpt().Gpt_keys[0],
		Field: openai.CompletionRequest{
			Model:            openai.GPT3TextDavinci003,
			Prompt:           msg,
			MaxTokens:        config.GetConfigGpt().Max_tokens,
			Temperature:      config.GetConfigGpt().Temperature,
			TopP:             config.GetConfigGpt().Top_p,
			N:                config.GetConfigGpt().N,
			PresencePenalty:  config.GetConfigGpt().Presence_penalty,
			FrequencyPenalty: config.GetConfigGpt().Frequency_penalty,
		},
	}
	return parm
}
func gpt35_4(msg string, assistant []string, system string) gpt.GptChatCompletion {
	var messages []openai.ChatCompletionMessage
	if system != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: system,
		})
	}
	if len(assistant) > 0 {
		for index, v := range assistant {
			if index%2 == 0 {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: v,
				})
			} else {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: v,
				})
			}
		}
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})

	param := gpt.GptChatCompletion{
		Key: config.GetConfigGpt().Gpt_keys[0],
		Field: openai.ChatCompletionRequest{
			Model:            openai.GPT3Dot5Turbo,
			Messages:         messages,
			MaxTokens:        config.GetConfigGpt().Max_tokens,
			Temperature:      config.GetConfigGpt().Temperature,
			TopP:             config.GetConfigGpt().Top_p,
			N:                config.GetConfigGpt().N,
			PresencePenalty:  config.GetConfigGpt().Presence_penalty,
			FrequencyPenalty: config.GetConfigGpt().Frequency_penalty,
		},
	}
	return param
}

func GetGptData(msg string, old []string, system string) string {
	str := ""
	switch config.GetConfigGpt().Model {
	case openai.GPT432K0314:
		fallthrough
	case openai.GPT432K:
		fallthrough
	case openai.GPT40314:
		fallthrough
	case openai.GPT4:
		fallthrough
	case openai.GPT3Dot5Turbo0301:
		fallthrough
	case openai.GPT3Dot5Turbo:
		str = gpt.GptRequest(gpt35_4(msg, old, system), GPT_TIMEOUT)
	case openai.GPT3TextDavinci003:
		fallthrough
	case openai.GPT3TextDavinci002:
		fallthrough
	case openai.GPT3TextCurie001:
		fallthrough
	case openai.GPT3TextBabbage001:
		fallthrough
	case openai.GPT3TextAda001:
		fallthrough
	case openai.GPT3TextDavinci001:
		fallthrough
	case openai.GPT3DavinciInstructBeta:
		fallthrough
	case openai.GPT3Davinci:
		fallthrough
	case openai.GPT3CurieInstructBeta:
		fallthrough
	case openai.GPT3Curie:
		fallthrough
	case openai.GPT3Ada:
		fallthrough
	case openai.GPT3Babbage:
		str = gpt.GptRequest(gpt0_3(msg), GPT_TIMEOUT)

	}

	return str
}

func GetInsAndQuestrion(ding *dingtalk.DingTalk, ques string) error {
	var Err error
	ins := ""
	question := ques
	assistant := []string{}
	system := ""
	res := ""
	// txtType := config.MSG_KEY_TXT
	rep := regexp.MustCompile(`^\s*/(\w+)\s+(.*)`)
	sp := rep.FindStringSubmatch(question)
	if len(sp) != 0 {
		ins = sp[1]
		// 将ins转换为小写
		ins = strings.ToLower(ins)
		question = sp[2]
	}
	switch ins {
	case "h":
		res = config.HELP
		// txtType = config.MSG_KEY_MD
	case "b":
		data, err := gpt.GetBalance(config.GetConfigGpt().Gpt_keys[0], GPT_TIMEOUT)
		if err != nil {
			Err = err
			break
		}
		res = fmt.Sprintf("当前已经使用：%v 剩余：%v", data.TotalUsed, data.TotalAvailable)
	case "c":
		msql := sqlite.NewSql()
		user := sqlite.User{
			Name:         ding.SenderNick,
			DingUserName: ding.SenderNick,
			DingUserId:   ding.SenderStaffId,
		}
		msql.AddUserCheckDingId(user)
		cx, err := msql.GetContextByUidAndSize(user.Id, config.TOKEN_MAX) //config.DATABASE_CONTEXT_MAX msql.GetContextByUid(user.Id)//
		if err != nil {
			Err = err
			break
		}
		for _, v := range cx {
			assistant = append(assistant, v.Question)
			assistant = append(assistant, v.Answer)
		}
		res = GetGptData(question, assistant, system)

		go func() {
			if res != "" {
				c := sqlite.Context{
					Uid:      user.Id,
					Question: question,
					Answer:   res,
				}
				msql.AddContext(c)
			}
		}()
	case "cc":
		msql := sqlite.NewSql()
		user, err := msql.GetUserByDingId(ding.SenderStaffId)
		if err != nil {
			Err = err
			res = config.ERROR + config.CLEAR_ERROR
			break
		}
		err = msql.DeleteContextByUid(user.Id)
		if err != nil {
			Err = err
			res = config.ERROR + config.CLEAR_ERROR
			break
		}
		res = config.CLEAR_SUCCESS

	// case "m":
	// 	question = config.MD_INS + question
	// 	res = GetGptData(question, assistant, system)
	// txtType = config.MSG_KEY_MD
	case "i":
		question = config.MD_IMG_INS + question
		res = GetGptData(question, assistant, system)
		// txtType = config.MSG_KEY_MD
	default:
		res = GetGptData(question, assistant, system)
	}
	lg.Infof("gpt返回的消息:%s", res)
	if res == "" {
		res = config.ERROR + config.VISIT_TIMEOUT
	}
	// if txtType == config.MSG_KEY_TXT {
	// 	err = ding.SendText(res)
	// } else if txtType == config.MSG_KEY_MD {
	// 	err = ding.SendMarkDown(res)
	// }
	err := ding.SendMarkDown(res)
	lg.Infof("----------------END------------------")
	if err != nil {
		return errors.New(Err.Error() + "\n" + err.Error())
	}
	return nil
}
