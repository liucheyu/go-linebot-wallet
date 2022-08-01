package handler

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/liucheyu/go-linebot-wallet/pkg/linebot/service"
)

var (
	AmountRegexp, _        = regexp.Compile("^\\d+")
	insertItemYesRegexp, _ = regexp.Compile("[是|yes]")
	insertItemNoRegexp, _  = regexp.Compile("[否|yes]")
	UserService            service.UserActionService
	BaseDataCache          service.BaseDataCache
	Bot                    *linebot.Client

	handlerMap map[string]func(context *Context)
	nextMap    map[string]string
)

type Context struct {
	HandlerName string
	UserID      string
	Token       string
	RequestText string
}

func init() {
	handlerMap = map[string]func(context *Context){}
	nextMap = map[string]string{}

	addHandler("amount", askInsertAmount)
	addHandler("itemName", askInsertItemName)
	addHandler("checkItemName", checkInsertItemIsCurrect)
	addHandler("itemType", askSelectItemType)
	addHandler("payMethod", askSelectPayMethod)

	addHandler("notMatch", insertNotCurrect)

	setNextHandler("amount", "itemName")
	setNextHandler("itemName", "checkItemName")
	setNextHandler("checkItemName", "itemType")
	setNextHandler("itemType", "payMethod")
}

func GetNextHandler(userID string) string {
	userMap, err := UserService.GetUserCacheMap(userID)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return nextMap[userMap["handler"]]
}

func ExcuteHandler(context *Context) {
	handlerMap[context.HandlerName](context)
}

func MappingFirstStat(requestText string) string {
	if AmountRegexp.Match([]byte(requestText)) {
		return "amount"
	}

	return "notMatch"
}

func addHandler(handlerName string, handler func(context *Context)) {
	handlerMap[handlerName] = handler
}

func setNextHandler(handlerName string, nextHandlerName string) {
	nextMap[handlerName] = nextHandlerName
}

func insertNotCurrect(context *Context) {
	Bot.ReplyMessage(context.Token, linebot.NewTextMessage("輸入錯誤，請重新輸入金額或點選選單。")).Do()
}

func askInsertAmount(context *Context) {
	if !AmountRegexp.Match([]byte(context.RequestText)) {
		insertNotCurrect(context)
		return
	}

	cacheMap := map[string]string{"handler": context.HandlerName, "amount": context.RequestText}
	UserService.CacheUserDataMap(context.UserID, cacheMap, 300)
	Bot.ReplyMessage(context.Token, linebot.NewTextMessage("請輸入記帳名稱")).Do()
}

func askInsertItemName(context *Context) {
	userMap, err := UserService.GetUserCacheMap(context.UserID)

	if err != nil {
		fmt.Println(err)
		return
	}

	itemTypesMap := BaseDataCache.GetItemTypesMap()
	typeID := func() int {
		for itemTypeID, itemTypeName := range itemTypesMap {
			if itemTypeName == context.RequestText {
				res, _ := strconv.Atoi(itemTypeID)
				return res
			}
		}
		return 0
	}()

	if typeID == 0 {
		Bot.ReplyMessage(context.Token, linebot.NewTextMessage("資料錯誤")).Do()
		return
	}

	userMap["itemName"] = context.RequestText
	userMap["handler"] = context.HandlerName
	UserService.CacheUserDataMap(context.UserID, userMap, 300)

	Bot.ReplyMessage(context.Token, linebot.NewTextMessage(fmt.Sprintf("您輸入的為： %s, 是否正確？[輸入 是,yes,否,no]", context.RequestText))).Do()
}

func checkInsertItemIsCurrect(context *Context) {
	userMap, err := UserService.GetUserCacheMap(context.UserID)

	if err != nil {
		fmt.Println(err)
		return
	}

	var isNo bool

	if insertItemYesRegexp.Match([]byte(context.RequestText)) {
		isNo = false
		userMap["handler"] = "checkItemName"
	}

	if insertItemNoRegexp.Match([]byte(context.RequestText)) {
		isNo = true
		userMap["handler"] = "amount"
	}

	UserService.CacheUserDataMap(context.UserID, userMap, 300)

	if isNo {
		return
	}

	template1 := linebot.NewButtonsTemplate(
		"", "項目類型", "請選擇項目類型",
		linebot.NewPostbackAction("食", "1", "", "食"),
		linebot.NewPostbackAction("衣", "2", "", "衣"),
		linebot.NewPostbackAction("住", "3", "", "住"),
	)

	template2 := linebot.NewButtonsTemplate(
		"", "項目類型", "請選擇項目類型",
		linebot.NewPostbackAction("行", "4", "", "行"),
		linebot.NewPostbackAction("育", "5", "", "育"),
		linebot.NewPostbackAction("樂", "6", "", "樂"),
	)

	Bot.PushMessage(context.UserID,
		linebot.NewTemplateMessage("項目類型", template1))
	Bot.PushMessage(context.UserID,
		linebot.NewTemplateMessage("項目類型", template2))
}

func askSelectItemType(context *Context) {
	userMap, err := UserService.GetUserCacheMap(context.UserID)

	if err != nil {
		fmt.Println(err)
		return
	}

	userMap["itemType"] = context.RequestText
	userMap["handler"] = context.HandlerName

	UserService.CacheUserDataMap(context.UserID, userMap, 300)

	//1:cash 2:creditCard 3:bankTransfer 4:nfc 5:qrcode 6:mobilePay 7:eTicket
	template1 := linebot.NewButtonsTemplate(
		"", "支付類型", "請選擇支付類型",
		linebot.NewPostbackAction("現金", "1", "", "現金"),
		linebot.NewPostbackAction("信用卡", "2", "", "信用卡"),
		linebot.NewPostbackAction("轉帳", "3", "", "轉帳"),
		linebot.NewPostbackAction("NFC", "4", "", "NFC"),
	)

	template2 := linebot.NewButtonsTemplate(
		"", "支付類型", "請選擇支付類型",
		linebot.NewPostbackAction("QR Code", "5", "", "QR Code"),
		linebot.NewPostbackAction("行動支付", "6", "", "行動支付"),
		linebot.NewPostbackAction("電子支付", "7", "", "電子支付"),
	)

	Bot.PushMessage(context.UserID,
		linebot.NewTemplateMessage("支付類型", template1)).Do()
	Bot.PushMessage(context.UserID,
		linebot.NewTemplateMessage("支付類型", template2)).Do()
}

func askSelectPayMethod(context *Context) {
	userMap, err := UserService.GetUserCacheMap(context.UserID)

	if err != nil {
		fmt.Println(err)
	}

	userMap["payMethod"] = context.RequestText
	affected, err := UserService.SaveUserDataMapToDB(userMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	if affected <= 0 {
		Bot.ReplyMessage(context.Token, linebot.NewTextMessage("資料錯誤")).Do()
		UserService.DeleteUserDataMapCache(context.UserID)
		return
	}

	Bot.ReplyMessage(context.Token, linebot.NewTextMessage("儲存成功")).Do()

	UserService.DeleteUserDataMapCache(context.UserID)
}
