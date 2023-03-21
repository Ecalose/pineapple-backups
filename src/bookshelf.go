package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/Template"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	HbookerAPI "github.com/VeronicaAlexia/HbookerAPI/ciweimao/bookshelf"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
)

type Bookshelf struct {
	ShelfIndex    int
	ShelfBookData *[]Template.ShelfData
	ShelfBook     map[string]int
}

func NewChoiceBookshelf() *Bookshelf {
	var bs Bookshelf
	bs.ShelfBookData = boluobao.API.BookShelf.NovelBookShelf()
	if bs.ShelfBookData != nil {
		if len(*bs.ShelfBookData) == 1 {
			fmt.Println("检测到只有一个书架，无需选择书架")
			bs.ShelfIndex = 0
			//bs.ShowBookshelf(response)

		} else {
			fmt.Println("检测到多个书架，需要选择书架")
			for index, value := range *bs.ShelfBookData {
				fmt.Println("书架号:", index, "\t\t\t书架名:", value.Name)
			}
			bs.ShelfIndex = tools.InputInt(">", len(*bs.ShelfBookData))
			//bs.ShowBookshelf(response)
		}
		return &bs

	}
	return nil
}

func (bs *Bookshelf) NewSfacgBookshelf() {
	for i, value := range *bs.ShelfBookData {
		if i == bs.ShelfIndex {
			if value.Expand.Novels != nil {
				bs.ShelfBook = make(map[string]int)
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"序号", "小说名", "小说ID"})
				for index, bookInfo := range value.Expand.Novels {
					table.Append([]string{strconv.Itoa(index), bookInfo.NovelName, strconv.Itoa(bookInfo.NovelID)})
					//bs.ShelfBook = append(bs.ShelfBook, map[int]int{index: bookInfo.NovelID})
					bs.ShelfBook[strconv.Itoa(index)] = bookInfo.NovelID
				}
				table.Render()
			} else {
				fmt.Println(value.Name, "书架中没有书籍!")
			}
		}
	}

}

func hbooker_bookshelf() (map[int][]map[string]string, error) {
	bookshelf_info := make(map[int][]map[string]string)
	shelf := HbookerAPI.GET_BOOK_SHELF_INFORMATION()

	if shelf.Code != "100000" {
		return nil, fmt.Errorf(shelf.Tip.(string))
	}
	for index, value := range shelf.Data.ShelfList {
		fmt.Println("书架号:", index, "\t\t\t书架名:", value.ShelfName)
		bookshelf_index := HbookerAPI.GET_BOOK_SHELF_INDEXES_INFORMATION(value.ShelfID)
		if bookshelf_index.Code != "100000" {
			return nil, fmt.Errorf(bookshelf_index.Tip.(string))
		}
		var bookshelf_info_list []map[string]string
		for _, book := range bookshelf_index.Data.BookList {
			bookshelf_info_list = append(bookshelf_info_list,
				map[string]string{"novel_name": book.BookInfo.BookName, "novel_id": book.BookInfo.BookID},
			)
		}
		bookshelf_info[index] = bookshelf_info_list
	}
	return bookshelf_info, nil
}

func request_bookshelf_book_list() (map[int][]map[string]string, error) {
	if config.Vars.AppType == "sfacg" {
		//return sfacg_bookshelf()
		return nil, nil
	} else if config.Vars.AppType == "cat" {
		return hbooker_bookshelf()
	} else {
		return nil, fmt.Errorf("app type error")
	}
}

func InitBookShelf() ([]int, []map[string]string) {
	bookshelf_book_list, response_err := request_bookshelf_book_list()
	if response_err != nil || bookshelf_book_list == nil {
		var test_login_status bool
		fmt.Println("config.Vars.AppType ", config.Vars.AppType)
		fmt.Println("BookShelf Error:", response_err)
		if config.Vars.AppType == "sfacg" {
			test_login_status = AutoAccount()
		} else if config.Vars.AppType == "cat" {
			test_login_status = InputAccountToken()
		}
		if !test_login_status && config.Vars.AppType == "sfacg" {
			fmt.Println("please login your sfacg account and password!")
			account := tools.InputStr("please input your account:")
			password := tools.InputStr("please input your password:")
			LoginAccount(strings.TrimSpace(account), strings.TrimSpace(password), 0)
		}
		return InitBookShelf()

	}

	return select_bookcase(bookshelf_book_list)

}

func select_bookcase(bookshelf_book_list map[int][]map[string]string) ([]int, []map[string]string) {
	var bookshelf_index int
	if len(bookshelf_book_list) == 1 {
		fmt.Println("you only have one bookshelf, default loading bookshelf index:1")
		bookshelf_index = 0
	} else {
		fmt.Println("please input bookshelf index:")
		bookshelf_index = tools.InputInt(">", len(bookshelf_book_list))
	}
	book_shelf_bookcase := bookshelf_book_list[bookshelf_index]
	var bookshelf_book_index []int
	for book_index, book := range book_shelf_bookcase {
		fmt.Println("index:", book_index, "\t\tid:", book["novel_id"], "\t\tname:", book["novel_name"])
		bookshelf_book_index = append(bookshelf_book_index, book_index)
	}
	return bookshelf_book_index, book_shelf_bookcase
}

//func (bs *Bookshelf) ShowBookshelf(shelf *[]Template.ShelfData) {
//	for index, value := range *shelf {
//		fmt.Println("书架号:", index, "\t\t\t书架名:", value.Name)
//	}
//
//}
//
//func (bs *Bookshelf) ChoiceBookshelf(BookInfoData []Template.BookInfoData) *Template.BookInfoData {
//	for _, book := range BookInfoData {
//		fmt.Println("小说名:", book.NovelName, "\t\t\t小说ID:", book.NovelID)
//	}
//	choice := tools.InputStr(">")
//	if strings.Contains(choice, "d ") {
//		book_id := strings.Replace(choice, "d ", "", 1)
//		res := boluobao.API.Book.NovelInfo(book_id)
//		if res != nil {
//			fmt.Println(res.NovelName, res.NovelID)
//			return res
//		}
//
//	}
//	if strings.ToLower(choice) == "y" {
//		BookIndex := tools.InputInt(">", len(BookInfoData))
//		boluobao.API.Book.NovelInfo(strconv.Itoa(BookInfoData[BookIndex].NovelID))
//	} else if strings.ToLower(choice) == "a" {
//		for _, book := range BookInfoData {
//			boluobao.API.Book.NovelInfo(strconv.Itoa(book.NovelID))
//		}
//	} else {
//		fmt.Println("已退出书架下载")
//	}
//	return nil
//}

//func sfacg_bookshelf() (map[int][]map[string]string, error) {
//	response := boluobao.API.BookShelf.NovelBookShelf()
//	bookshelf_info := make(map[int][]map[string]string)
//	if response == nil {
//		return nil, fmt.Errorf("get bookshelf error")
//	}
//	for index, value := range *response {
//		fmt.Println("书架号:", index, "\t\t\t书架名:", value.Name)
//		var bookshelf_info_list []map[string]string
//		for _, book := range value.Expand.Novels {
//			bookshelf_info_list = append(bookshelf_info_list,
//				map[string]string{"novel_name": book.NovelName, "novel_id": strconv.Itoa(book.NovelID)},
//			)
//		}
//		bookshelf_info[index] = bookshelf_info_list
//	}
//	return bookshelf_info, nil
//}
