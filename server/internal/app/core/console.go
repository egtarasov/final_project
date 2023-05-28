package core

import (
	"context"
	"errors"
	"fmt"
	"homework-5/server/internal/app/group"
	"homework-5/server/internal/app/student"
	"strconv"
	"strings"
)

const (
	helpMessage = `Available commands:
	students
			Get [id] retrieve student from db with Id = id
			Add [id | first_name | middle_name | second_name | gpa | attendance_rate | groupId] insert student in db
			Update [id | first_name | middle_name | second_name | gpa | attendance_rate | groupId] update student in db
			Delete [id] delete student from db
	groups
			Get [id] retrieve group from db with Id = id
			Add [id | name | year] insert group in db
			Update [id | name | year] update group in db
			Delete [id] delete group from db
	spell word 
		spell [word], which contains only of latin letters, by letters
	gofmt ../filename.txt
		format file`
)

var (
	InvalidInput    = errors.New("invalid input")
	ProcessingError = errors.New("processing error")
)

type Console struct {
	ctx         context.Context
	studentRepo student.Repository
	groupRepo   group.Repository
}

func NewConsole(ctx context.Context, studentRepo student.Repository, groupRepo group.Repository) *Console {
	return &Console{
		ctx:         ctx,
		studentRepo: studentRepo,
		groupRepo:   groupRepo}
}

type consoleCommand interface {
	Process(ctx context.Context, params []string) (string, error)
}

func (c *Console) Action(command string) (string, error) {
	params := strings.Split(command, " ")
	if command == "help" {
		fmt.Println(helpMessage)
		return "", nil
	}
	var consoleCommand consoleCommand

	switch params[0] {
	case "students":
		consoleCommand = NewStudentCommand(c.studentRepo)
	case "groups":
		consoleCommand = NewGroupCommand(c.groupRepo)
	case "spell":
		consoleCommand = NewSpellCommand()
	case "gofmt":
		consoleCommand = NewDummyGoFmt()
	default:
		return "", InvalidInput
	}

	return consoleCommand.Process(c.ctx, params[1:])
}

func getId(params []string) (int64, error) {
	return strconv.ParseInt(params[1], 10, 64)
}
