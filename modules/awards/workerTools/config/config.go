package config

import (
	"os"

	"github.com/SteMak/ani-helper/workerTools/bankirapi"
)

// Usage var of basic syntax
const Usage = "Общий синтаксис:\n\t-запрос <причина> -> <заслужившие награду пользователи> <сумма для выдачи>, <заслужившие другую награду пользователи> <сумма для выдачи>\n" +
	"Пример:\n\t-запрос Победители ивента \"Шляпа\" -> @someone @elseone 100, @anotherone 200"

var (
	// Token discord token
	Token string
	// BankirapiToken bankirapi token
	BankirapiToken string

	// API an interface for bankir bot
	API *bankirapi.API

	// GdHouseID guildID
	GdHouseID string

	// ChForRequestID channelID
	ChForRequestID string
	// ChForLogsID channelID
	ChForLogsID string
	// ChForBumpSiupID for monitoring
	ChForBumpSiupID string

	// UsConfirmatorID userID
	UsConfirmatorID string
	// UsSiupID siup user
	UsSiupID string
	// UsBumpID bump user
	UsBumpID string

	// RoRequestMakerID roleID
	RoRequestMakerID string

	// PostgresURI uri of postgress
	PostgresURI string

	// Responces beautiful log
	Responces = []string{
		"Вы сделали %s сервера и Тихий Ужас вручил Вам %s",
		"Вы героически поймали %s и Вас наградили %s",
		"за %s сервера Жмяк отдала Вам свои печеньки и вы получили %s",
		"после кибератаки вы подняли сервер своим %sом и получили %s",
		"пожертвовав жизнью на войне за %s, Вас посмертно наградили %s",
		"смакуючи ельфійською абракадаброю (%s), ви начарували %s",
		"сыр съел сырный сырник %sая сервер, Вам заплатили моральную компенсацию в %s",
		"Вы помогли Ведьмаку с %sом, за что Вам заплатили __ЧЕКАННЫМИ__ %s",
		"Вы сделали %s сервера и Скромный Модератор вручил Вам %s",
		"Вы съели свою поджелудочную во время %sа и нашли в ней %s",
		"Вы собрали с тысячи людей по АниКоину и %sнули сервер. Вы получили %s",
		"Вы попытались выписать бан Кнопычу, но сделали %s и получили %s",
		"на вечерних посиделках с Хикаро вы сражались за %s. Хикаро Вас наградил %s",
		"Вы сохранили последние пять минут и угрожающе сделали %s. Глюк расстрогался и отдал Вам %s",
		"Вы наблюдали за программированием Стёмы и Меро, но не забыли сделать %s и получили %s",
		"Меро уходил спать, но %sнул за Вас и вы получили %s",
		"Нев отвлёк всех разговорами об отсутствии холодильника и вы %sнули сервер. Холодильник дал вам %s",
		"Кемпер поднимал актив, чтобы кемперить сакуру, так что вы тихо %sнули сервер и получили %s",
		"Вы вычислили Маргинала и он рассказал вам секрет %sанья. Получено %s",
		"у Боннуса провис интернет, так что вы беспрепятствеено %sнули сервер и забрали %s",
		"Эксля заснул, что-то бормоча во сне: \"Z-z-z... %s Z-z-z... %s Z-z-z...\"",
		"Эспада зажала Вас в тиски, но вы успели сделать %s и они Вас зауважали и дали %s",
		"Маю-Маю снова уснула в войсе, вы вдохновились её ворочанием и %sнули, заработав %s",
		"ɔıloqoou написал %s перевёрнутыти буквами, поэтому вы неспеша забрали %s",
		"Фузу мирно рисовала в войсе, а Вы сделали %s и собрали %s",
	}

	// FairyReplacement for good parse
	FairyReplacement = [][2]string{
		[2]string{"  ", " "},
		[2]string{",,", ","},
		[2]string{", ,", ","},
		[2]string{"><", "> <"},
		[2]string{">1", "> 1"},
		[2]string{">2", "> 2"},
		[2]string{">3", "> 3"},
		[2]string{">4", "> 4"},
		[2]string{">5", "> 5"},
		[2]string{">6", "> 6"},
		[2]string{">7", "> 7"},
		[2]string{">8", "> 8"},
		[2]string{">9", "> 9"},
		[2]string{">0", "> 0"},
	}
)

// Init inits main vars for the project
func Init() {
	Token = os.Getenv("TOKEN")
	BankirAPIToken = os.Getenv("BANKIRAPI_TOKEN")

	API = bankirapi.New(BankirAPIToken)

	GdHouseID = os.Getenv("GD_HOUSE_ID")

	ChForRequestID = os.Getenv("CH_FOR_REQUEST_ID")
	ChForLogsID = os.Getenv("CH_FOR_LOGS_ID")
	ChForBumpSiupID = os.Getenv("CH_FOR_BUMP_SIUP")

	UsConfirmatorID = os.Getenv("US_CONFIRMATOR_ID")
	UsSiupID = os.Getenv("US_SIUP_ID")
	UsBumpID = os.Getenv("US_BUMP_ID")

	RoRequestMakerID = os.Getenv("RO_REQUEST_MAKER_ID")

	PostgresURI = os.Getenv("DATABASE_URL")
}
