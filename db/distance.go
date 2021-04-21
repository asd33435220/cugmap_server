package db

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const (
	Pi                = 3.1415926535
	queryAllUserInfo2 = "select student_id,username,lng,lat,signature from cug_map_users_tpl where student_id != ?;"
)

type UserWithDistance struct {
	StudentId string  `json:"student_id"`
	Username  string  `json:"username"`
	Distance  float64 `json:"distance"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Signature string  `json:"signature"`
}

func Getdistance(lng1, lat1, lng2, lat2 float64) (distance float64) {
	user1Lng := lng1 * Pi / 180
	user1Lat := lat1 * Pi / 180
	user2Lng := lng2 * Pi / 180
	user2Lat := lat2 * Pi / 180
	hsinX := math.Sin((user1Lng - user2Lng) / 2)
	hsinY := math.Sin((user1Lat - user2Lat) / 2)
	h := hsinY*hsinY + math.Cos(user1Lat)*math.Cos(user2Lat)*hsinX*hsinX
	distance = 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h)) * 6367000 / 1000
	return
}

func getTargetList(userList []*UserWithDistance) (targetUserList []*UserWithDistance) {
	for i := 0; i < 30; i++ {
		targetUserList = append(targetUserList, findLaskK(userList, 1))
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
			userList = append(userList[:pivot], userList[pivot+1:]...)
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
func GetAllUserInfo2(lng, lat float64, myId string) (userList []*UserWithDistance, err error) {
	userList = make([]*UserWithDistance, 0)
	start1 := time.Now() // 获取当前时间
	stmt2, err2 := Db.Prepare(queryAllUserInfo2)
	if err2 != nil {
		fmt.Println("err2", err2.Error())
		return
	}
	defer stmt2.Close()
	rows, err := stmt2.Query(myId)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := UserWithDistance{}
		rows.Scan(&user.StudentId, &user.Username, &user.Longitude, &user.Latitude, &user.Signature)
		if user.Longitude != 0 && user.Latitude != 0 && user.StudentId != myId {

			userList = append(userList, &user)
		}
	}
	elapsed1 := time.Since(start1)
	fmt.Println("数据库执行完成耗时：", elapsed1)
	start2 := time.Now() // 获取当前时间
	for _, user := range userList {
		user.Distance = Getdistance(user.Longitude, user.Latitude, lng, lat)
	}
	elapsed2 := time.Since(start2)
	fmt.Println("坐标计算函数执行完成耗时：", elapsed2)
	start3 := time.Now() // 获取当前时间
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].Distance < userList[j].Distance
	})
	userList = userList[:50]
	elapsed3 := time.Since(start3)
	fmt.Println("排序函数执行完成耗时：", elapsed3)
	return
}
