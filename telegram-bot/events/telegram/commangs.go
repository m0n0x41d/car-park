package telegram

import (
	"car-park/telegram-bot/lib/er"
	"car-park/telegram-bot/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	LoginCmd        = "/login"
	LogoutCmd       = "/logout"
	ShowCredsCmd    = "/mycreds"
	HelpCmd         = "/help"
	HelpReportsCmd  = "/help_reports"
	StartCmd        = "/start"
	PingCmd         = "/ping"
	MilageReportCmd = "/milage"
)

var (
	CARPARK_PORT = os.Getenv("CARPARK_PORT")
	CARPARK_HOST = os.Getenv("CARPARK_HOST")
)

type timeSlice []time.Time

func (s timeSlice) Less(i int, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)          { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int               { return len(s) }

func (p *Processor) processCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	textSlice := splitText(text)
	maybeCmd := textSlice[0]

	log.Printf("got new command %s from %s", text, userName)

	switch maybeCmd {
	case LoginCmd:
		return p.loginTry(chatID, text, userName)
	case ShowCredsCmd:
		return p.showMyCreds(chatID, userName)
	case LogoutCmd:
		return p.logout(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID, userName)
	case HelpReportsCmd:
		return p.sendHelpReports(chatID, userName)
	case StartCmd:
		return p.sendHello(chatID, userName)
	case PingCmd:
		return p.pingCarPark(chatID, userName)
	case MilageReportCmd:
		return p.milageReport(chatID, text, userName)
	default:

	}

	return nil
}

func (p *Processor) milageReport(chatID int, text string, userName string) (err error) {
	r := regexp.MustCompile("[^\\s]+")
	textSlice := r.FindAllString(text, -1)

	fmt.Println(textSlice)

	sliceLength := len(textSlice)
	if sliceLength != 6 {
		return p.tg.SendMessage(chatID, msgInvalidReportRequest)
	}

	creds, err := p.storage.GetCredentials(userName)
	if err != nil {
		p.tg.SendMessage(chatID, msgNotLoggedIn)
		return er.Wrap("can't find credentials of user", err)
	}

	requestUrl := buildRequestURL(creds.CarParkLogin, creds.CarParkPassword, CARPARK_HOST, CARPARK_PORT,
		textSlice[1], textSlice[2], textSlice[3], textSlice[4], "MilageReport", textSlice[5])

	var Report Report
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err)
		return er.Wrap("[ERROR]: can't get report ", err)
	}

	if resp.StatusCode == 401 {
		return p.tg.SendMessage(chatID, msgInvalidCredentials)
	} else if resp.StatusCode != 200 {
		return p.tg.SendMessage(chatID, msgInvalidReportRequest)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	errUnmarshal := json.Unmarshal(body, &Report)
	if errUnmarshal != nil {
		p.tg.SendMessage(chatID, "Can't process report, try again")
		return er.Wrap("can't unmarshal report", err)
	}

	if Report.ReportType == "" {
		return p.tg.SendMessage(chatID, msgInvalidReportRequest)
	}

	utcMap := make(map[time.Time]interface{})

	for k, v := range Report.Results {
		utcMap[k.UTC()] = v
	}

	utcKeys := make(timeSlice, 0, len(Report.Results))
	for k := range utcMap {
		utcKeys = append(utcKeys, k)
	}
	sort.Sort(utcKeys)

	for _, k := range utcKeys {
		message := fmt.Sprintf("For time range started at %s milage is %.2f km", k, utcMap[k])
		p.tg.SendMessage(chatID, message)
	}

	return nil
}

func buildRequestURL(login, password, host, port, managerID, vehicleID, fromDate, toDate, reportType, timePeriod string) string {
	URL := fmt.Sprintf("http://%s:%s@%s:%s/api/manager/%s/vehicles/%s/reports?notBefore=%s&notAfter=%s&report_type=%s&time_period=%s",
		login, password, host, port, managerID, vehicleID, fromDate, toDate, reportType, timePeriod)
	return URL
}

func (p *Processor) pingCarPark(chatID int, userName string) (err error) {

	creds, err := p.storage.GetCredentials(userName)
	if err != nil {
		p.tg.SendMessage(chatID, msgNotLoggedIn)
		return er.Wrap("can't find credentials of user", err)
	}

	respCode := pingCarParkWithCreds(CARPARK_HOST, CARPARK_PORT, creds.CarParkLogin, creds.CarParkPassword)

	switch respCode {
	case 200:
		p.tg.SendMessage(chatID, msgPingOk)
		return nil
	case 401:
		p.tg.SendMessage(chatID, msgPingOk401)
		return fmt.Errorf("[WARNING]: CarPark credentials unauthorized for user %s", userName)
	default:
		p.tg.SendMessage(chatID, msgPingNotOk)
		log.Printf("[ERROR]: CarPark pinged with %s response code for user %s", strconv.Itoa(respCode), userName)
		return fmt.Errorf("[WARNING]: CarPark possibly unavailable")
	}
}

func (p *Processor) loginTry(chatID int, text string, userName string) (err error) {
	r := regexp.MustCompile("[^\\s]+")
	textSlice := r.FindAllString(text, -1)

	fmt.Println(textSlice)

	sliceLength := len(textSlice)
	if sliceLength != 3 {
		return p.tg.SendMessage(chatID, msgInvalidCredentials)
	}

	creds := &storage.Credentials{
		Username:        userName,
		CarParkLogin:    textSlice[1],
		CarParkPassword: textSlice[2],
	}

	respCode := pingCarParkWithCreds(CARPARK_HOST, CARPARK_PORT, creds.CarParkLogin, creds.CarParkPassword)

	isExist, err := p.storage.IsExistsCredentials(userName)
	if err != nil {
		return er.Wrap("can't check credentials existence while saving credentials", err)
	}

	if isExist {
		return p.tg.SendMessage(chatID, msgCredentialsExists)
	}

	if respCode == 401 {
		return p.tg.SendMessage(chatID, msgInvalidCredentials)
	}

	if err := p.storage.SaveCredentials(creds); err != nil {
		return er.Wrap("can't save credentials", err)
	}

	if err := p.tg.SendMessage(chatID, msgCredentialsSaved); err != nil {
		return er.Wrap("can't send message about successfull credentials save by Processor", err)
	}

	return nil
}

func (p *Processor) showMyCreds(chatID int, userName string) (err error) {
	creds, err := p.storage.GetCredentials(userName)
	if err != nil {
		p.tg.SendMessage(chatID, msgNotLoggedIn)
		return er.Wrap("can't find credentials of user", err)
	}
	if err := p.tg.SendMessage(chatID, fmt.Sprintf("You logged in with LOGIN = %s, and PASSWORD = %s", creds.CarParkLogin, creds.CarParkPassword)); err != nil {
		return er.Wrap("can't send user credentials to user", err)
	}

	return nil
}

func (p *Processor) logout(chatID int, userName string) (err error) {
	creds, err := p.storage.GetCredentials(userName)
	if err != nil {
		p.tg.SendMessage(chatID, msgNotLoggedIn)
		return er.Wrap("can't find credentials of user", err)
	}

	if err := p.storage.RemoveCredentials(creds); err != nil {
		return er.Wrap("can't log out user", err)
	} else {
		p.tg.SendMessage(chatID, "You are logged out successfully, bye.")
	}

	return nil
}

// func buildRequestURL(host, port, managerID, vehicleID, fromDate, toDate, reportType, timePeriod string) string {
// 	URL := fmt.Sprintf("%s:%s/api/manager/%s/vehicles/%s/reports?notBefore=%s&notAfter=%s&report_type=%s&time_period=%s",
// 		host, port, managerID, vehicleID, fromDate, toDate, reportType, timePeriod)
// 	return URL
// }

func pingCarParkWithCreds(host string, port string, managerLogin string, managerPassword string) int {

	resp, err := http.Get("http://" + managerLogin + ":" + managerPassword + "@" + host + ":" + port)
	if err != nil {
		log.Println(err)
		return -1
	}
	if resp.StatusCode == 404 {
		log.Println("404 from CarPark.")
		return 404
	}

	if resp.StatusCode == 401 {
		log.Println("401 from CarPark. Invalid Credentials")
		return 401
	}

	if resp.StatusCode == 200 {
		return 200
	}

	return -1
}

func (p *Processor) sendHelp(chatID int, userName string) (err error) {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHelpReports(chatID int, userName string) (err error) {
	return p.tg.SendMessage(chatID, msgHelpReports)
}

func (p *Processor) sendHello(chatID int, userName string) (err error) {
	return p.tg.SendMessage(chatID, msgHello)
}

func splitText(text string) []string {
	r := regexp.MustCompile("[^\\s]+")
	textSlice := r.FindAllString(text, -1)
	return textSlice
}
