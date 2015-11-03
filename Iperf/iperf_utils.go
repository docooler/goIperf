package Iperf

import (
		"fmt"
		"os"
		)


func HandleError(err error, needExit int, addInfo string){
	if err != nil {
		fmt.Println(addInfo)
		fmt.Println(err.Error())
		if needExit != 0 {
			os.Exit(2)
		}
	}
}