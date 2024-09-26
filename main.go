package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Post struct {
	UserId int
	Id     int
	Title  string
	Body   string
}

func getPost(postId int, recivedPosts *map[int]Post, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get("https://winry.khashaev.ru/posts/" + strconv.Itoa(postId))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var post Post

	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil { // заглушка на пост который не смогли считать
		errPost := Post{
			Id:     postId,
			UserId: -1,
			Title:  "error in JSON decoding",
			Body:   "error in JSON decoding",
		}
		mu.Lock()
		(*recivedPosts)[postId] = errPost
		mu.Unlock()
		return
	}

	mu.Lock() //чтобы нельзя было два раза писать в один postId
	(*recivedPosts)[postId] = post
	mu.Unlock()
}

func readPosts(sliceId []int) map[int]Post {
	posts := make(map[int]Post)

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, val := range sliceId {
		wg.Add(1)
		go getPost(val, &posts, mu, wg)
	}
	wg.Wait()
	return posts
}

func main() {
	sliceId := make([]int, 0)

	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&sliceId)
	if err != nil {
		fmt.Println(err)
	}

	posts := readPosts(sliceId)

	for _, key := range sliceId {
		fmt.Printf("key: %d\n%v\n\n", key, posts[key].Body)
	}
}
