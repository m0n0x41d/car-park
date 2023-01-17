package telegram

const msgHelp = `This is CarParkDev helper bot. In order to make requests to CarParkDev service you should login with manager credentials by sending me command:
"/login <USERNAME> <PASSWORD>".

Use command /logout if you need to login with another credentials pair.

Use command /mycreds to check your current CarPark credentials you logged in with.

Use command /ping to check is CarPark serive available and your credentials are valid.

Use command /help_reports to see help on commands generating reports.

Use /help to see this message again.
`

const msgHelpReports = `

Currently the milage report is the only one report available to fetch.
Of course you need to be logged in with valid CarPark manager credentials.

To check this use /mycreds or /ping command.

For milage report you need to specify a several request parameters:

1. Your manager ID (ask your team lead if you don't know this ID);
2. VehicleID for which you want to generate report;
3. From date (including) in format YYYY-MM-DD;
4. To date (including) in format YYYY-MM-DD;
5. Time Period for report results. Valid options is: ByDay, ByWEek, or ByMonth;

Example of valid command:

/milage 17 554432 2023-01-01 2023-12-31 DAILY

All of parameters are required.

If request is valid you will get answer, one message for each report entry with
time range starting timestamp and corresponding value. 

For example, if you ask for Daily report, you will get messages with every first ride of every
day with one or more rides, and sum of milage for all rides in this day.

If you ask for weekly report — timestamp will be a timestamp of first ride in each week, etc.
`

//TODO: Add help about reports fetch command when will be implemented.

const msgHello = "Hello Friend. \n\n" + msgHelp

const (
	msgUnknownCommand       = "I do not know this command. See /help"
	msgNotLoggedIn          = `You are not logged in. Use "/login <USERNAME> <PASSWORD>"`
	msgCredentialsExists    = `You are already logged in. If you need to logout use "/logout" command`
	msgCredentialsSaved     = `You are successfully logged in. See /help_reports to check what you can fetch.`
	msgInvalidCredentials   = `Invalid credentials for CarParkDev, try /login again`
	msgInvalidReportRequest = `Invalid report request, see /help_reports`

	msgPingOk    = `CarPark is up, your credentials is valid`
	msgPingOk401 = `CarPark is up, but your credentials is invalid, try /login again`
	msgPingNotOk = `CarPark is unaccessible`
)
