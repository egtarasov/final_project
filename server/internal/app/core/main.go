package core

import (
	"bufio"
	"context"
	"fmt"
	"homework-5/server/internal/app/database"
	group2 "homework-5/server/internal/app/group"
	student2 "homework-5/server/internal/app/student"
	"log"
	"os"
	"strings"
)

func Run() {
	ctx := context.Background()

	client, err := database.NewDb(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	studentRepo := student2.NewStudentsRepository(client)
	groupRepo := group2.NewGroupsRepository(client)

	cons := NewConsole(ctx, studentRepo, groupRepo)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.Replace(line, "\r\n", "", -1)
		response, err := cons.Action(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(response)
	}
}
