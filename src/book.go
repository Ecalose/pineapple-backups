package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/epub"
	"github.com/VeronicaAlexia/pineapple-backups/src/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/hbooker"
	"github.com/VeronicaAlexia/pineapple-backups/src/https"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	"os"
	"path"
	"strings"
)

type BookInits struct {
	BookID      string
	ShowBook    bool
	Locks       *config.GoLimit
	EpubSetting *epub.Epub
}

func (books *BookInits) InitEpubFile() {
	AddImage := true                                                // add image to epub file
	books.EpubSetting = epub.NewEpub(config.Current.Book.NovelName) // set epub setting and add section
	books.EpubSetting.SetAuthor(config.Current.Book.AuthorName)     // set author
	if !config.Exist(config.Current.CoverPath) {
		if reader := https.GetCover(config.Current.Book.NovelCover); reader == nil {
			fmt.Println("download cover failed!")
			AddImage = false
		} else {
			_ = os.WriteFile(config.Current.CoverPath, reader, 0666)
		}
	}
	if AddImage {
		_, _ = books.EpubSetting.AddImage(config.Current.CoverPath, "")
		books.EpubSetting.SetCover(strings.ReplaceAll(config.Current.CoverPath, "cover", "../images"), "")
	}

}

func SettingBooks(book_id string) Catalogue {
	var err error
	var result _struct.Books
	switch config.Vars.AppType {
	case "sfacg":
		result, err = boluobao.GET_BOOK_INFORMATION(book_id)
	case "cat":
		result, err = hbooker.GET_BOOK_INFORMATION(book_id)
	}
	if err == nil {
		config.Current.Book = result
		config.Current.ConfigPath = path.Join(config.Vars.ConfigName, config.Current.Book.NovelName)
		config.Current.OutputPath = path.Join(config.Vars.OutputName, config.Current.Book.NovelName+".txt")
		config.Current.CoverPath = path.Join("cover", config.Current.Book.NovelName+".jpg")
		books := BookInits{BookID: book_id, Locks: nil, ShowBook: true}
		return books.BookDetailed()
	} else {
		return Catalogue{Test: false, BookMessage: fmt.Sprintf("book_id:%v is invalid:%v", book_id, err)}
	}

}

func (books *BookInits) BookDetailed() Catalogue {
	books.InitEpubFile()
	briefIntroduction := fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
		config.Current.Book.NovelName, config.Current.Book.NovelID, config.Current.Book.AuthorName, config.Current.Book.CharCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	config.Write(config.Current.OutputPath, briefIntroduction, "w")
	return Catalogue{Test: true, EpubSetting: books.EpubSetting}
}
