package main

// #cgo LDFLAGS: -O3 -lsqlite3
// #include "Simple-ATM-With-C/SHA-256/sha256.c"
// #include "Simple-ATM-With-C/SHA-256/sha256.h"
// #include "Simple-ATM-With-C/ATM/atm.c"
// #include "Simple-ATM-With-C/ATM/atm.h"
import "C"

import (
	"fmt"
	"os"
)

func main() {
	if C.initDB() != C.OK {
		os.Exit(1)
	}

	defer C.closeDB()

	var user C.struct_BankUser
	var CardId string
	var choose int
	var value int

	for {
		fmt.Print("Enter Your CardID: ")
		fmt.Scanln(&CardId)

		if CardId == "exit" {
			return
		}

		if len(CardId) > 10 {
			C.handelError(C.NOT_VALID_CARD_ID)
			continue
		}

		password := C.getpass(C.CString("Password: "))
		if C.handelError(C.makeLoginUser(&user, C.CString(CardId), password)) != C.OK {
			continue
		}

		C.clearPrint(C.accountInfo(&user))

		for {
			fmt.Println("Choose One Action!")
			fmt.Println("1 - Take Money From ATM")
			fmt.Println("2 - Show Account Info")
			fmt.Println("3 - Send Money For Another")
			fmt.Println("4 - Exit From ATM")
			fmt.Print("Enter Number: ")
			fmt.Scanln(&choose)
			fmt.Println("\033c")

			if choose == 1 {
				fmt.Print("Enter Value: ")
				fmt.Scanln(&value)
				if C.handelError(C.getMoney(&user, C.int64_t(value))) != C.OK {
					continue
				}

				C.clearPrint(C.CString(
					fmt.Sprintf(
						"Get Your Money (%d)\n"+
							"Fee Is (%d)\n"+
							"Your money now: %d",
						value, C.FEE, user.Money),
				))
			} else if choose == 2 {
				C.clearPrint(C.accountInfo(&user))
			} else if choose == 3 {
				var cardId string
				fmt.Print("Enter Value: ")
				fmt.Scanln(&value)
				if C.handelError(C.getMoney(&user, C.int64_t(value))) != C.OK {
					continue
				}

				fmt.Print("Enter Your CardID: ")
				fmt.Scanln(&cardId)

				if len(cardId) > 10 {
					C.handelError(C.NOT_VALID_CARD_ID)
					continue
				}

				if C.handelError(C.sendMoney(&user, C.CString(cardId), C.int64_t(value))) != C.OK {
					continue
				}

				C.clearPrint(C.CString(
					fmt.Sprintf(
						"Send Money (%d)\n"+
							"For (%s)\n"+
							"Fee Is (%d)\n"+
							"Your money now: %d",
						value, cardId, C.FEE, user.Money),
				))
			} else if choose == 4 {
				C.clearPrint(C.CString("See You!"))
				break
			} else {
				C.clearPrint(C.CString("Choose Not Valid!\n"))
			}
		}
	}
}
