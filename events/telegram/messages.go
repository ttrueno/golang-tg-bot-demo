package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me a link.(it must start from host, like 'http://' and others)

To get random one(from your list ofc), type command /rnd.
WARNING: page will be deleted automatically after fetching, so u must read it!`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnkownCommand = "Unkown command 😵‍💫😵‍💫😵‍💫"
	msgNoSavedPage   = "You have no saved links 👺👺👺"
	msgSaved         = "Saved! 😈😈😈"
	msgAlreadyExists = "You already saved this link ☠️☠️☠️"
)
