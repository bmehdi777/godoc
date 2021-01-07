package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Documentation struct {
	gorm.Model
	Title      string
	Command    string
	Definition string
}

func main() {
	db, err := gorm.Open(sqlite.Open("doc.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to conenct to the database.")
	}
	db.AutoMigrate(&Documentation{})
	handleArg(db)
}

func handleArg(db *gorm.DB) {
	arrayArg := [...]string{"-l\t- List every documentations you stored.", "-h\t- List every args possible.", "-a\t- Add a command", "-r\t- Remove a command"}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-l":
			showCmd(db)
		case "-h":
			for _, a := range arrayArg {
				color.Yellow("%v\n", a)
			}
		case "-a":
			promptAdd(db)
		case "-e":
			promptEdit(db)
		case "-r":
			promptRemove(db)
		}
	} else {
		color.Red("Use doc -h to see every arg.")
	}
}

func showCmd(db *gorm.DB) {
	var res []Documentation
	db.Find(&res)
	if len(res) > 0 {
		for _, r := range res {
			color.Yellow("%v :\n", r.Title)
			color.Green("Command : %v\n", r.Command)
			color.White("%v\n\n", r.Definition)
		}
	} else {
		color.Red("No doc to show")
	}
}

func promptAdd(db *gorm.DB) {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("Value too short")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Title",
		Validate: validate,
	}
	resTitleCmd, err := prompt.Run()
	if err != nil {
		color.Red("Prompt failed : %v\n", err)
		return
	}

	prompt = promptui.Prompt{
		Label:    "Command",
		Validate: validate,
	}
	resCmd, err := prompt.Run()
	if err != nil {
		color.Red("Prompt failed :: %v\n", err)
		return
	}

	prompt = promptui.Prompt{
		Label:    "Description",
		Validate: validate,
	}
	resDescCmd, err := prompt.Run()
	if err != nil {
		color.Red("Prompt failed : %v\n", err)
		return
	}

	prompt = promptui.Prompt{
		Label:     "Confirm",
		IsConfirm: true,
	}
	confirm, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if strings.ToLower(confirm) == "y" {
		doc := Documentation{Title: resTitleCmd, Command: resCmd, Definition: resDescCmd}
		db.Create(&doc)
		color.Green("Created.\nYou can see it by using 'doc -l'")
	} else {
		color.Red("Creation canceled")
	}
}
func promptRemove(db *gorm.DB) {
	var res []Documentation
	db.Find(&res)
	if len(res) > 0 {
		displayRes := make([]string, 0, len(res))
		resDict := make(map[int]uint)
		for i, r := range res {
			displayRes = append(displayRes, r.Command+" - "+r.Title+" - "+r.Definition)
			resDict[i] = r.ID
		}
		prompt := promptui.Select{
			Label: "Select Day",
			Items: displayRes,
		}

		index, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		db.Delete(&Documentation{}, resDict[index])
		color.Green("You succesfully removed :%v %v", index, result)
	} else {
		color.Red("No doc to remove")
	}

}

func promptEdit(db *gorm.DB) {
	var res []Documentation
	db.Find(&res)
	if len(res) > 0 {

	} else {
		color.Red("No doc to edit")
	}
}
