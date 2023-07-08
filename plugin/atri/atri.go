/*
Package atri 本文件基于 https://github.com/Kyomotoi/ATRI
为 Golang 移植版，语料、素材均来自上述项目
本项目遵守 AGPL v3 协议进行开源
*/
package atri

import (
	"encoding/base64"
	"math/rand"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
)

type datagetter func(string, bool) ([]byte, error)

func (dgtr datagetter) randImage(file ...string) message.MessageSegment {
	data, err := dgtr(file[rand.Intn(len(file))], true)
	if err != nil {
		return message.Text("ERROR: ", err)
	}
	return message.ImageBytes(data)
}

func (dgtr datagetter) randRecord(file ...string) message.MessageSegment {
	data, err := dgtr(file[rand.Intn(len(file))], true)
	if err != nil {
		return message.Text("ERROR: ", err)
	}
	return message.Record("base64://" + base64.StdEncoding.EncodeToString(data))
}

func randText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

// isAtriSleeping 凌晨0点到6点，ATRI 在睡觉，不回应任何请求
func isAtriSleeping(*zero.Ctx) bool {
	if now := time.Now().Hour(); now >= 1 && now < 6 {
		return false
	}
	return true
}

func init() { // 插件主体
	engine := control.Register("atri", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "atri人格文本回复",
		Help: "本插件基于 ATRI ，为 Golang 移植版\n" +
			"- ATRI醒醒\n- ATRI睡吧\n- 萝卜子\n- 喜欢 | 爱你 | 爱 | suki | daisuki | すき | 好き | 贴贴 | 老婆 | 亲一个 | mua\n" +
			"- 草你妈 | 操你妈 | 脑瘫 | 废柴 | fw | 废物 | 战斗 | 爬 | 爪巴 | sb | SB | 傻B\n- 早安 | 早哇 | 早上好 | ohayo | 哦哈哟 | お早う | 早好 | 早 | 早早早\n" +
			"- 中午好 | 午安 | 午好\n- 晚安 | oyasuminasai | おやすみなさい | 晚好 | 晚上好\n- 高性能 | 太棒了 | すごい | sugoi | 斯国一 | よかった\n" +
			"- 没事 | 没关系 | 大丈夫 | 还好 | 不要紧 | 没出大问题 | 没伤到哪\n- 好吗 | 是吗 | 行不行 | 能不能 | 可不可以\n- 啊这\n- 我好了\n- ？ | ? | ¿\n" +
			"- 离谱\n- 答应我",
		PublicDataFolder: "Atri",
		OnEnable: func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("休眠中…请勿打扰。"))
		},
		OnDisable: func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("Zzz……Zzz……"))
		},
	})
	engine.UsePreHandler(isAtriSleeping)
	var dgtr datagetter = engine.GetLazyData
	engine.OnFullMatch("萝卜子").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			switch rand.Intn(2) {
			case 0:
				ctx.SendChain(randText("萝卜子是对机器人的蔑称哦。", "请称呼我为TREE。"))
			case 1:
				ctx.SendChain(dgtr.randRecord("RocketPunch.amr"))
			}
		})
	engine.OnFullMatchGroup([]string{"喜欢", "爱你", "爱", "suki", "daisuki", "すき", "好き", "贴贴", "老婆", "亲一个", "mua"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(dgtr.randImage("SUKI.jpg", "SUKI1.jpg", "SUKI2.png"))
		})
	engine.OnKeywordGroup([]string{"草你妈", "操你妈", "脑瘫", "废柴", "fw", "five", "废物", "战斗", "爬", "爪巴", "sb", "SB", "傻B"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(dgtr.randImage("FN.jpg", "WQ.jpg", "WQ1.jpg"))
		})
	engine.OnFullMatchGroup([]string{"早安", "早哇", "早上好", "ohayo", "哦哈哟", "お早う", "早好", "早", "早早早"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now().Hour()
			switch {
			case now < 6: // 凌晨
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"zzzz......",
					"zzzzzzzz......",
					"zzz...zzz....",
					"zzz....zzz..",
				))
			case now >= 6 && now < 9:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"早安。",
					"早上好，已为您备份好昨日的工作数据。",
				))
			case now >= 9 && now < 18:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"现在的时间不符合人类定义的早上哦，您熬夜了吗？",
					"您是现在才起床吗？较多临床数据指出早睡早起的人能活更久。",
				))
			case now >= 18 && now < 24:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"现在应该说晚安更合适哦。",
					"晚上好。",
				))
			}
		})
	engine.OnFullMatchGroup([]string{"中午好", "午安", "午好"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now().Hour()
			if now > 11 && now < 15 { // 中午
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"午安:->",
					"鉴于您昨日的作息安排，推荐您进行午休，TREE会在必要的时候唤醒您。",
					"中午好。",
				))
			}
		})
	engine.OnFullMatchGroup([]string{"晚安", "oyasuminasai", "おやすみなさい", "晚好", "晚上好"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now().Hour()
			switch {
			case now < 6: // 凌晨
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"zzzz......",
					"zzzzzzzz......",
				))
			case now >= 6 && now < 11:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"当心猝死。",
					"？",
					"这边建议您立刻去睡觉。",
				))
			case now >= 11 && now < 15:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"午安:->",
					"鉴于您昨日的作息安排，推荐您进行午休，TREE会在必要的时候唤醒您。",
					"中午好。",
				))
			case now >= 15 && now < 19:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"您还有未完成的工作，现在还不能休息哦。",
					"现在还是白天哦，睡觉还太早了。",
				))
			case now >= 19 && now < 24:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), randText(
					"晚安，系统祝您一夜好梦。",
					"已为您定好明早的闹钟，祝安眠。",
				))
			}
		})
	engine.OnKeywordGroup([]string{"高性能", "太棒了", "すごい", "sugoi", "斯国一", "よかった"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText(
				"因为我是高性能且高智能的AI。",
				"毕竟这是造物主们的智慧结晶呢。",
			))
		})
	engine.OnKeywordGroup([]string{"没事", "没关系", "大丈夫", "还好", "不要紧", "没出大问题", "没伤到哪"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText(
				"没关系，TREE正在协助解决。",
				"不用担心，系统会协助您解决。",
			))
		})

	engine.OnKeywordGroup([]string{"好吗", "是吗", "行不行", "能不能", "可不可以"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if rand.Intn(2) == 0 {
				ctx.SendChain(dgtr.randImage("YES.png", "NO.jpg"))
			}
		})
	engine.OnKeywordGroup([]string{"啊这"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if rand.Intn(2) == 0 {
				ctx.SendChain(dgtr.randImage("AZ.jpg", "AZ1.jpg"))
			}
		})
	engine.OnKeywordGroup([]string{"我好了"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), randText("？", "我坏了。"))
		})
	engine.OnFullMatchGroup([]string{"？", "?", "¿"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			switch rand.Intn(5) {
			case 0:
				ctx.SendChain(randText("?", "？", "嗯？", "ん？"))
			case 1, 2:
				ctx.SendChain(dgtr.randImage("WH.jpg", "WH1.jpg", "WH2.jpg", "WH3.jpg"))
			}
		})
	engine.OnKeyword("离谱").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			switch rand.Intn(5) {
			case 0:
				ctx.SendChain(randText("?", "？", "嗯？", "(。´・ω・)ん?", "ん？"))
			case 1, 2:
				ctx.SendChain(dgtr.randImage("WH.jpg"))
			}
		})
	engine.OnKeyword("答应我", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText("我无法回应你的请求"))
		})
}
