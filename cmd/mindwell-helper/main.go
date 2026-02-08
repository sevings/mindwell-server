package main

import (
	"bufio"
	"go.uber.org/zap"
	"log"
	"os"

	"github.com/sevings/mindwell-server/helper"
	"github.com/sevings/mindwell-server/utils"
)

const admArg = "adm"
const cleanupArg = "cleanup"
const helpArg = "help"
const mailArg = "mail"
const logArg = "log"
const surveyArg = "survey"
const webappArg = "webapp"

func printHelp() {
	log.Printf(
		`
Usage: mindwell-helper [option] [flags]

Options are:
%s		- set grandfathers in adm and sent emails to them.
%s		- cleanup orphaned image files.
		  Flags:
		    --unused    Check for unused database images (not in entries, >6 months old)
		    --delete    Actually delete files (requires confirmation)
		    --verbose   Show detailed progress
%s		- send email reminders.
%s		- send email survey.
%s		- import user request log.
%s		- create the official Mindwell web app.
%s		- print this help message.
`, admArg, cleanupArg, mailArg, surveyArg, logArg, webappArg, helpArg)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == helpArg {
		printHelp()
		return
	}

	cfg := utils.LoadConfig("configs/server")

	baseURL, _ := cfg.String("server.base_url")
	support, _ := cfg.String("server.support")
	moderator, _ := cfg.String("server.moderator")

	if len(support) == 0 || len(moderator) == 0 || len(baseURL) == 0 {
		log.Println("Check config consistency")
	}

	zapLog, err := zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		log.Println(err)
	}

	emailLog := zapLog.With(zap.String("type", "email"))

	mail := &utils.Postman{
		BaseUrl:   baseURL,
		Support:   support,
		Moderator: moderator,
		Logger:    emailLog,
	}

	smtpHost, _ := cfg.String("email.host")
	smtpPort, _ := cfg.Int("email.port")
	smtpUsername, _ := cfg.String("email.username")
	smtpPassword, _ := cfg.String("email.password")
	smtpHelo, _ := cfg.String("email.helo")
	err = mail.Start(smtpHost, smtpPort, smtpUsername, smtpPassword, smtpHelo)
	if err != nil {
		log.Println(err.Error())
	}

	surveyUrl, _ := cfg.String("email.survey")

	db := utils.OpenDatabase(cfg)
	tx := utils.NewAutoTx(db)
	defer tx.Finish()

	args := os.Args[1:]
	for _, arg := range args {
		// Skip flags (arguments starting with --)
		if len(arg) >= 2 && arg[0:2] == "--" {
			continue
		}

		switch arg {
		case admArg:
			helper.UpdateAdm(tx, mail)
		case cleanupArg:
			helper.CleanupOrphanedImages(tx, cfg)
		case mailArg:
			helper.SendReminders(tx, mail)
		case surveyArg:
			helper.SendSurvey(tx, mail, surveyUrl)
		case logArg:
			helper.ImportUserLog(tx)
		case webappArg:
			helper.CreateWebApp(tx, cfg)
		case helpArg:
			printHelp()
		default:
			log.Printf("Unknown argument: %s\n", arg)
		}

		log.Println("(press Enter to continue)")

		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
	}
}
