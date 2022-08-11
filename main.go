package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sf/cfg"
	"sf/multi"
	"sf/src"
	"sf/structural"
	"strings"
)

func shellBookDownload(downloadId any) {
	switch downloadId.(type) {
	case string:
		start := src.BookInits{BookID: downloadId.(string), Index: 0, Locks: nil, ShowBook: true}
		catalogues := start.DownloadBookInit() // get book catalogues
		if catalogues.TestBookResult {
			catalogues.InitCatalogue()
		}
	case []string:
		Locks := multi.NewGoLimit(7)
		for BookIndex, BookId := range downloadId.([]string) {
			Locks.Add()
			start := src.BookInits{BookID: BookId, Index: BookIndex, Locks: Locks, ShowBook: true}
			catalogues := start.DownloadBookInit() // get book catalogues
			if catalogues.TestBookResult {
				catalogues.InitCatalogue()
			}
		}
		Locks.WaitZero() // wait for all goroutines to finish
	}
	os.Exit(0) // exit the program if no error
}

func shellLoginAccount(inputs []string) bool {
	if cfg.Vars.AppType == "sfacg" {
		if len(inputs) >= 3 {
			src.LoginAccount(inputs[1], inputs[2], 0)
		} else {
			fmt.Println("you must input account and password, like: sf account password")
		}
	} else if cfg.Vars.AppType == "cat" {
		if ok := src.InputAccountToken(); !ok {
			fmt.Println("you must input account and token.")
		}
	}
	return true
}

func shellBookMain(inputs []string) {
	if len(inputs) == 2 {
		if cfg.Vars.AppType == "cat" {
			if len(inputs[1]) == 9 { // test if the input is hbooker book id
				shellBookDownload(inputs[1])
			} else {
				fmt.Println("hbooker bookid is 9 characters, please input again:")
			}
		} else {
			shellBookDownload(inputs[1])
		}
	} else {
		fmt.Println("input book id or url, like:download <bookid/url>")
	}
}

func shellSearchBookMain(inputs []string) {
	if len(inputs) == 2 {
		if NovelId := src.SearchBook(inputs[1]); NovelId != "" {
			shellBookDownload(NovelId)
		} else {
			fmt.Println("No found search book, please input again:")
		}
	} else {
		fmt.Println("input book id or url, like:download <bookid/url>")
	}
}

func ParseCommandLine() structural.Command {
	download := flag.String("download", "", "input book id or url")
	account := flag.String("account", "", "input account")
	password := flag.String("password", "", "input password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: search keyword")
	thread := flag.Int("max", 0, "input thread number, like: thread 1")
	showConfig := flag.Bool("show", false, "show config, like: show config")
	flag.Parse()
	if *thread > 0 && *thread < 64 {
		cfg.Vars.MaxThreadNumber = *thread
	}
	if *account != "" && *password != "" {
		shellConsole([]string{"login", *account, *password})
	} else {
		cfg.Vars.AppType = *appType
		src.TestAppTypeAndAccount()
	}
	return structural.Command{Download: *download, Search: *search, ShowConfig: *showConfig}
}

func shellConsole(inputs []string) {
	switch inputs[0] {
	case "a", "app":
		cfg.Vars.AppType = inputs[1]
		src.TestAppTypeAndAccount()
	case "q", "quit":
		os.Exit(0)
	case "h", "help":
		fmt.Println("help:")
	case "show", "test":
		cfg.FormatJson(cfg.ReadConfig(""))
	case "book", "download":
		shellBookMain(inputs)
	case "s", "search":
		shellSearchBookMain(inputs)
	case "l", "login":
		shellLoginAccount(inputs)
	default:
		fmt.Println("command not found,please input help to see the command list:", inputs[0])
	}

}
func init() {
	cfg.ConfigInit()
	cfg.Vars.SaveFile = "save"
	cfg.Vars.ConfigFile = "cache"
	cfg.Vars.AppType = "sfacg"
	cfg.Vars.MaxThreadNumber = 32
	cfg.Vars.MaxRetry = 5 // retry times when failed
}

func main() {
	if len(os.Args) <= 1 {
		src.TestAppTypeAndAccount()
		for {
			spaceRe, _ := regexp.Compile(`\s+`)
			inputs := spaceRe.Split(strings.TrimSpace(cfg.Input(">")), -1)
			if len(inputs) > 1 {
				shellConsole(inputs)
			} else if inputs[0] != "" {
				fmt.Println("you must input command, like: sf command")
			}
			os.Exit(1)
		}
	} else {
		var CommandLine []string
		ArgsCommandLine := ParseCommandLine()
		if ArgsCommandLine.ShowConfig {
			CommandLine = []string{"show", "config"}
		}
		if ArgsCommandLine.Download != "" {
			CommandLine = []string{"download", ArgsCommandLine.Download}
		}
		if ArgsCommandLine.Search != "" {
			CommandLine = []string{"search", ArgsCommandLine.Search}
		}
		if len(CommandLine) > 0 {
			shellConsole(CommandLine)
		}
	}
}
