package db

import (
	"html"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Pi = 3.1415926535
)

type UserWithDistance struct {
	StudentId string  `json:"student_id"`
	Username  string  `json:"username"`
	Distance  float64 `json:"distance"`
	Position  string  `json:"position"`
	Signature string  `json:"signature"`
}

func Getdistance(position string, user *User) (userWithD *UserWithDistance) {
	user2Position := user.Position
	var str string
	arr := strings.Split(position, ";")
	str = arr[0]
	user1Lng, _ := strconv.ParseFloat(str, 32)
	str = arr[1]
	user1Lat, _ := strconv.ParseFloat(str, 32)
	arr = strings.Split(user2Position, ";")
	str = arr[0]
	user2Lng, _ := strconv.ParseFloat(str, 32)
	str = arr[1]
	user2Lat, _ := strconv.ParseFloat(str, 32)
	user1Lng = user1Lng * Pi / 180
	user1Lat = user1Lat * Pi / 180
	user2Lng = user2Lng * Pi / 180
	user2Lat = user2Lat * Pi / 180
	hsinX := math.Sin((user1Lng - user2Lng) / 2)
	hsinY := math.Sin((user1Lat - user2Lat) / 2)
	h := hsinY*hsinY + math.Cos(user1Lat)*math.Cos(user2Lat)*hsinX*hsinX
	distance := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h)) * 6367000 / 1000
	signature := html.EscapeString(user.Signature)
	userWithD = &UserWithDistance{
		Username:  user.Username,
		StudentId: user.StudentId,
		Position:  user.Position,
		Distance:  distance,
		Signature: signature,
	}
	return
}

func getTargetList(userList []*UserWithDistance) (targetUserList []*UserWithDistance) {
	//targetUserList = make([]*UserWithDistance,0)
	//fmt.Println("userList", userList)
	arr := getRandList()
	seed := rand.NewSource(time.Now().Unix()) //同前面一样的种子
	randNum := rand.New(seed)
	for i := 0; i < 3+randNum.Intn(4); i++ {
		targetUserList = append(targetUserList, findLaskK(userList, arr[i]))
	}
	return
}

func findLaskK(userList []*UserWithDistance, k int) *UserWithDistance {
	var walk func(userList []*UserWithDistance, left int, right int, k int) *UserWithDistance
	walk = func(userList []*UserWithDistance, left int, right int, k int) *UserWithDistance {
		if left >= right || len(userList) <= 1 {
			return userList[left]
		}
		pivot := findPivot(userList, left, right)
		if pivot+1 == k {
			return userList[pivot]
		} else if pivot+1 < k {
			return walk(userList, pivot+1, right, k)
		} else {
			return walk(userList, left, pivot-1, k)
		}

	}
	return walk(userList, 0, len(userList)-1, k)

}
func findPivot(userList []*UserWithDistance, left int, right int) int {
	pivot := userList[left].Distance
	for left < right {
		for left < right && pivot <= userList[right].Distance {
			right--
		}
		swap(userList, left, right)
		for left < right && pivot >= userList[left].Distance {
			left++
		}
		swap(userList, left, right)
	}
	return left
}
func swap(userList []*UserWithDistance, left int, right int) {
	if left == right {
		return
	}
	temp := userList[left]
	userList[left] = userList[right]
	userList[right] = temp
}
func getRandList() (arr []int) {
	arr = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	seed := rand.NewSource(time.Now().Unix()) //同前面一样的种子
	randNum := rand.New(seed)
	for i := len(arr) - 1; i > 5; i-- {
		number := randNum.Intn(i)
		arr = append(arr[:number], arr[number+1:]...)
	}

	return
}
