package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type LogItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []LogItem
}

func (u User) getUserActivity() string {
	out := fmt.Sprintf("ID: %d | EMAIL: %s\nActivity log:\n", u.id, u.email)
	for index, item := range u.logs {
		out += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp)
	}
	return out
}

var actions []string = []string{
	"LOGGED IN",
	"LOGGED OFF",
	"UPDATE RECORD",
	"DELETE RECORD",
	"CREATE RECORD",
}

func generateUsers(count int) []User {
	users := make([]User, count)
	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@gmail.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
	}

	return users
}

func generateLogs(count int) []LogItem {
	logs := make([]LogItem, count)

	for i := 0; i < count; i++ {
		logs[i] = LogItem{
			timestamp: time.Now(),
			action:    actions[rand.Intn(len(actions)-1)],
		}
	}

	return logs
}

func saveUserActivity(user User, wg *sync.WaitGroup) error {
	fmt.Printf("WRITING FILE FOR USER ID: %d\n\n", user.id)
	filename := fmt.Sprintf("logs/uuid_%d.txt", user.id)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, error := file.WriteString(user.getUserActivity())
	if error != nil {
		return err
	}
	wg.Done()
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())

	wg := &sync.WaitGroup{}

	t := time.Now()

	users := generateUsers(1000)

	for _, user := range users {
		wg.Add(1)
		go saveUserActivity(user, wg)
	}
	wg.Wait()
	fmt.Println("TIME ELAPSED: ", time.Since(t).String())
}
