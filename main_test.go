package main

import "testing"

func TestReadPost(t *testing.T) {
	var cases = [][]int{
		{1, 2, 3, 4, 5},
		{},
		{7, 4, 2, 1},
		{15, 4, 1, 7, 50, 3},
		{1, 2, 500, 3, 4},
	}

	for caseNum, postIds := range cases {
		posts := readPosts(postIds)
		for _, val := range postIds {
			if posts[val].Id != val {
				t.Error("case: ", caseNum,
					"\n\tresult: ", posts[val].Id,
					"\n\texpected: ", val)
			}
		}
	}
}
