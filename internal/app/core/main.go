package core

import (
	"bufio"
	"context"
	"fmt"
	"homework-5/internal/app/database"
	"homework-5/internal/app/group"
	"homework-5/internal/app/student"
	"log"
	"os"
	"strings"
)

type GroupRepository interface {
	Add(ctx context.Context, group *group.Group) (uint64, error)                 //create
	GetById(ctx context.Context, id uint64) (*group.Group, error)                //read
	UpdateById(ctx context.Context, id uint64, group *group.Group) (bool, error) // update
	Remove(ctx context.Context, id uint64) (bool, error)                         //delete
}

type StudentRepository interface {
	Add(ctx context.Context, student *student.Student) (uint64, error)                 //create
	GetById(ctx context.Context, id uint64) (*student.Student, error)                  //read
	UpdateById(ctx context.Context, id uint64, student *student.Student) (bool, error) // update
	Remove(ctx context.Context, id uint64) (bool, error)                               //delete
}

func Run() {
	ctx := context.Background()

	client, err := database.NewDb(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	studentRepo := student.NewStudentsRepository(client)
	groupRepo := group.NewGroupsRepository(client)

	cons := NewConsole(ctx, studentRepo, groupRepo)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.Replace(line, "\r\n", "", -1)
		err = cons.Action(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}
