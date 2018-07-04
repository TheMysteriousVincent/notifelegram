package telegram

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/telegram-bot-api.v4"
)

const Version = "v0.1"

type Handler struct {
	bot           *tgbotapi.BotAPI
	sqlConnection *sql.DB
	gitlabClient  *gitlab.Client
}

func NewHandler(bot *tgbotapi.BotAPI, sqlConnection *sql.DB, gitlabClient *gitlab.Client) *Handler {
	return &Handler{
		bot:           bot,
		sqlConnection: sqlConnection,
		gitlabClient:  gitlabClient,
	}
}

func (h *Handler) HandleEnableCommits(msg *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, h.enableCommitsText(msg)))
}

func (h *Handler) enableCommitsText(msg *tgbotapi.Message) string {
	var nid int
	if err := h.sqlConnection.QueryRow(
		"INSERT INTO commits (chatId) SELECT $1 WHERE NOT EXISTS (SELECT commitId FROM commits WHERE chatId = $2) RETURNING commitId",
		msg.Chat.ID,
		msg.Chat.ID,
	).Scan(&nid); err != nil {
		if err == sql.ErrNoRows {
			return "Commits are already enabled."
		}
		return err.Error()
	}

	return "Enabled commit events."
}

func (h *Handler) HandleDisableCommits(msg *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, h.disableCommitsText(msg)))
}

func (h *Handler) disableCommitsText(msg *tgbotapi.Message) string {
	res, err := h.sqlConnection.Exec(
		"DELETE FROM commits WHERE chatId = $1",
		msg.Chat.ID,
	)
	if err != nil {
		return err.Error()
	}

	r, err := res.RowsAffected()
	if err != nil {
		return err.Error()
	}

	if r != 1 {
		return "Commits are already disabled."
	}

	return "Disabled commit events."
}

func (h *Handler) HandleAddMentions(msg *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, h.addMentions(msg)))
}

func (h *Handler) addMentions(msg *tgbotapi.Message) string {
	username := msg.CommandArguments()

	users, _, err := h.gitlabClient.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		log.Fatal(err.Error())
		return "User does not exist"
	}

	fmt.Println(users)

	var uid int
	for _, u := range users {
		if u.Username == username {
			uid = u.ID
			break
		}
	}

	if uid <= 0 {
		return "User does not exist"
	}

	var mentionID int
	if err := h.sqlConnection.QueryRow(
		"INSERT INTO mentions (chatId, gitlabUsername) SELECT $1, $2 WHERE NOT EXISTS (SELECT mentionId FROM mentions WHERE gitlabUsername = $3 AND chatId = $4) RETURNING mentionId",
		msg.Chat.ID,
		username,
		username,
		msg.Chat.ID,
	).Scan(&mentionID); err != nil {
		if err == sql.ErrNoRows {
			return "Mentions of that user are already subscribed."
		}
		return err.Error()
	}

	return "Successfully subscribed to user mentions."
}

func (h *Handler) HandleRemoveMentions(msg *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, h.removeMentionsText(msg)))
}

func (h *Handler) removeMentionsText(msg *tgbotapi.Message) string {
	username := msg.CommandArguments()

	users, _, err := h.gitlabClient.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		log.Fatal(err.Error())
		return "User does not exist"
	}

	fmt.Println(users)

	var uid int
	for _, u := range users {
		if u.Username == username {
			uid = u.ID
			break
		}
	}

	if uid <= 0 {
		return "User does not exist"
	}

	res, err := h.sqlConnection.Exec(
		"DELETE FROM mentions WHERE chatId = $1 AND gitlabUsername = $2",
		msg.Chat.ID,
		username,
	)
	if err != nil {
		return err.Error()
	}

	r, err := res.RowsAffected()
	if err != nil {
		return err.Error()
	}

	if r != 1 {
		return "Mentions of that user are already unsubscribed."
	}

	return "Successfully unsubscribed from that user mentions."
}

func (h *Handler) HandleListMentions(msg *tgbotapi.Message) {
	nMsg := tgbotapi.NewMessage(msg.Chat.ID, h.listMentionsText(msg))
	nMsg.ParseMode = "Markdown"
	h.bot.Send(nMsg)
}

type Mention struct {
	Username string
	DaysAgo  int
}

func (h *Handler) listMentionsText(msg *tgbotapi.Message) string {
	rows, err := h.sqlConnection.Query(
		"SELECT timestamp_add, gitlabUsername FROM mentions WHERE chatId = $1",
		msg.Chat.ID,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return err.Error()
		}
	}

	var arr []Mention
	for rows.Next() {
		var gitlabUsername, timestamp_add string

		err := rows.Scan(
			&timestamp_add,
			&gitlabUsername,
		)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		t, err := time.Parse(time.RFC3339Nano, timestamp_add)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		arr = append(arr, Mention{
			Username: gitlabUsername,
			DaysAgo:  int(math.Max(time.Since(t).Hours()/24, 1)),
		})
	}

	var tmpBuf bytes.Buffer
	if err := ParsedTemplateListMentions.Execute(&tmpBuf, arr); err != nil {
		return err.Error()
	}

	return tmpBuf.String()
}

func (h *Handler) HandleCommitsEnabled(msg *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, h.commitsText(msg)))
}

func (h *Handler) commitsText(msg *tgbotapi.Message) string {
	res, err := h.sqlConnection.Exec(
		"SELECT commitId FROM commits WHERE chatId = $1",
		msg.Chat.ID,
	)
	if err != nil {
		return "Commits are not subscribed."
	}

	r, err := res.RowsAffected()
	if err != nil {
		return "Commits are not subscribed."
	}
	if r != 1 {
		return "Commits are not subscribed."
	}

	return "Commits are subscribed."
}

func (h *Handler) HandleVersion(msg *tgbotapi.Message) {
	h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Version: %s", Version)))
}

func (h *Handler) HandleHelp(msg *tgbotapi.Message) {
	var tmpBuf bytes.Buffer
	nMsg := tgbotapi.NewMessage(msg.Chat.ID, "")
	err := ParsedTemplateHelp.Execute(&tmpBuf, msg.Chat)
	if err != nil {
		nMsg.Text = err.Error()
	} else {
		nMsg.Text = tmpBuf.String()
	}
	nMsg.ParseMode = "Markdown"
	h.bot.Send(nMsg)
}
