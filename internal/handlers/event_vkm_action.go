package handlers

import (
	"context"
	"lark-vkm/pkg/openai"
	"lark-vkm/pkg/qdrantkit"
	"lark-vkm/pkg/utils"
	"log"
	"strconv"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

// 知识库查询
type VkmAction struct {
	info   *ActionInfo
	client *qdrantkit.Qdrant
}

func (a *VkmAction) Execute(info *ActionInfo) bool {
	a.info = info
	a.client = qdrantkit.New(
		info.handler.config.QdrantHost, info.handler.config.QdrantCollection,
	)

	if msg, foundPrefix := utils.EitherCutPrefix(info.info.qParsed,
		"/faq ", "知识库 "); foundPrefix {
		info.handler.sessionCache.Clear(*info.info.sessionId)
		info.handler.sessionCache.SetMsg(*info.info.sessionId, a.loadEmbeddings(msg))
		info.info.qParsed = msg
		return true
	}
	return true
}

func (a *VkmAction) loadEmbeddings(msg string) []openai.Messages {
	// 计算向量
	response := a.info.handler.vectorCache.IfProcessed(msg)
	if response == nil {
		var err error
		response, err = a.info.handler.gpt.Embeddings(msg)
		if err != nil {
			log.Println("err:", err)
			return nil
		}
		a.info.handler.vectorCache.SetEmbedings(msg, response)
	}

	params := make(map[string]interface{})
	params["exact"] = false
	params["hnsw_ef"] = 128

	sr := qdrantkit.PointSearchRequest{
		Params:      params,
		Vector:      response.Data[0].Embedding,
		Limit:       3,
		WithPayload: true,
	}
	// 查询相似
	res, err := a.client.SearchPoints(sr)
	if err != nil {
		log.Println(err)
		return nil
	}

	// 组装本地数据
	localData := ""
	for i, v := range res {
		re := v.Payload.(map[string]interface{})
		//pp.Println(re)
		localData += "\n"
		localData += strconv.Itoa(i)
		localData += "."
		localData += re["Title"].(string)
		localData += ":"
		localData += re["Content"].(string)
	}

	messages := make([]openai.Messages, 0)

	q := "使用以下段落来回答问题，如果段落内容不相关就返回未查到相关信息："
	q += localData

	system := openai.Messages{
		Role:    "system",
		Content: "你是一个知识库问答接待员，你的回答需要根据提供的段落进行准确回答。",
	}
	assistant := openai.Messages{
		Role:    "assistant",
		Content: q,
	}

	messages = append(messages, system)
	messages = append(messages, assistant)

	return messages
}

// 知识库管理
type VkmOperationtAction struct {
	info   *ActionInfo
	client *qdrantkit.Qdrant
}

func (a *VkmOperationtAction) Execute(info *ActionInfo) bool {
	a.info = info
	a.client = qdrantkit.New(
		info.handler.config.QdrantHost, info.handler.config.QdrantCollection,
	)

	_, foundPrefix := utils.EitherCutPrefix(info.info.qParsed, "/faqmgr ", "知识库管理 ")
	if !foundPrefix {
		return true
	}

	// TODO: add vector database CRUD
	// 文件，链接，etc...
	operators := []string{}
	sendVkmOperationtInstructionCard(
		*info.ctx, info.info.sessionId, info.info.msgId, operators,
	)

	return true
}

func sendVkmOperationtInstructionCard(ctx context.Context,
	sessionId *string, msgId *string, operations []string) {
	newCard, _ := newSendCard(
		withHeader("🥷  已进入知识库管理", larkcard.TemplateIndigo),
		withVkmOperationtActionBtn(sessionId, operations...),
		withNote("请注意，以前操作将会修改知识数据库，且无法撤消已经执行操作"))
	replyCard(ctx, msgId, newCard)
}

func withVkmOperationtActionBtn(sessionID *string, operations ...string) larkcard.
MessageCardElement {
	var menuOptions []MenuOption

	for _, operation := range operations {
		menuOptions = append(menuOptions, MenuOption{
			label: operation,
			value: operation,
		})
	}
	cancelMenu := newMenu("选择需要执行的操作",
		map[string]interface{}{
			"value":     "0",
			"kind":      VkmOperationChooseKind,
			"sessionId": *sessionID,
			"msgId":     *sessionID,
		},
		menuOptions...,
	)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{cancelMenu}).
		Layout(larkcard.MessageCardActionLayoutFlow.Ptr()).
		Build()

	return actions
}
