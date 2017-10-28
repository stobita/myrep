package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	current_dir, _ := os.Getwd()
	app := cli.NewApp()
	app.Name = "MyGrep"
	app.Usage = "Provide Personal Grep!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Value: current_dir,
			Usage: "search target",
		},
	}
	app.Action = func(c *cli.Context) error {
		println("------ Start MyGrep ------")
		search_word := c.Args().Get(0)
		file_path := c.String("target")
		_, err := os.Stat(file_path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		println("I find [" + search_word + "] at " + file_path)
		err = filepath.Walk(file_path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				channel := make(chan int)
				go func(resc chan<- int) {
					myScan(path, search_word)
					defer close(resc)
				}(channel)
				for {
					ok := <-channel
					if ok == 0 {
						break
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}
	app.Run(os.Args)
}

func myWalk() {

}

func myScan(file_path string, search_word string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()
	file_index := 0
	match_flg := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		file_index++
		if strings.Index(scanner.Text(), search_word) > -1 {
			if !match_flg {
				//一度配列とかMapとかに入れた方がいい
				fmt.Println(file_path)
			}
			match_flg = true
			fmt.Print("\t" + strconv.Itoa(file_index) + "\t")
			if len(scanner.Text()) > 60 {
				fmt.Println(scanner.Text()[0:55] + ".....")
			} else {
				fmt.Println(scanner.Text())
			}
		}
	}
	return nil
}
