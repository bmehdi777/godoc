package main

import (
	"errors"
	"fmt"
	"os"

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
	// db.Create(&Documentation{Title: "test", Command: "ls -l", Definition: "Afficher tous les fichiers dans le dossier courant"})
	// var d Documentation
	// db.First(&d, 1)
	// fmt.Printf("%v %v %v\n", d.Title, d.Command, d.Definition)
	handleArg(db)
}

func handleArg(db *gorm.DB) {
	arrayArg := [...]string{"-l\t- List every documentations you stored.", "-h\t- List every args possible.", "-a\t- Add a command", "-r\t- Remove a command"}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-l":
			fmt.Printf("Afficher toutes les commandes\n")
		case "-h":
			for _, a := range arrayArg {
				color.Yellow("%v\n", a)
			}
		case "-a":
			promptAdd()
		case "-e":
			promptEdit()
		case "-r":
			promptRemove(db)
		}
	} else {
		color.Red("Use doc -h to see every arg.")
	}
}

func promptAdd() {
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
	color.Green("You choosed : %v, %v, %v", resTitleCmd, resCmd, resDescCmd)
}
func promptRemove(db *gorm.DB) {
	var res []Documentation
	db.Find(&res)
	displayRes := make([]string, 0, len(res))
	for _, r := range res {
		displayRes = append(displayRes, r.Command+" - "+r.Title+" - "+r.Definition)
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
	color.Green("You choosed :%v %v", index, result)
}

func promptEdit() {

}
