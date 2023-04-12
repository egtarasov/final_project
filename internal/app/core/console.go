package core

import (
	"context"
	"errors"
	"fmt"
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
			Add [id | name | amount_of_students | year] insert group in db
			Update [id | name | amount_of_students | year] update group in db
			Delete [id] delete group from db
	spell word 
		spell [word], which contains only of latin letters, by letters
	`
)

var (
	InvalidInput    = errors.New("invalid input")
	ProcessingError = errors.New("processing error")
)

type Console struct {
	ctx         context.Context
	studentRepo StudentRepository
	groupRepo   GroupRepository
}

func NewConsole(ctx context.Context, studentRepo StudentRepository, groupRepo GroupRepository) *Console {
	return &Console{
		ctx:         ctx,
		studentRepo: studentRepo,
		groupRepo:   groupRepo}
}

type consoleCommand interface {
	Process(ctx context.Context, params []string) error
}

func (c *Console) Action(command string) error {
	params := strings.Split(command, " ")
	if command == "help" {
		fmt.Println(helpMessage)
	}
	var consoleCommand consoleCommand

	switch params[0] {
	case "students":
		consoleCommand = NewStudentCommand(c.studentRepo)
	case "groups":
		consoleCommand = NewStudentCommand(c.studentRepo)
	case "spell":
		consoleCommand = NewSpellCommand()
	default:
		return InvalidInput
	}

	return consoleCommand.Process(c.ctx, params[1:])
}

func getId(params []string) (uint64, error) {
	return strconv.ParseUint(params[1], 10, 64)
}
