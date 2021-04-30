package db

import (
	Rdb "../redis"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Place struct {
	Name           string
	Address        string
	Longitude      float64
	Latitude       float64
	Image1Url      string
	Image2Url      string
	Type           int
	Score          float64
	CommentNumber  int
	Number         string
	PlaceCode      string
	Founder        string
	FounderComment string
}

const (
	AddPlaceStr         = "insert into cug_map_place_tpl(name,founder,founder_comment,score,type,phone_number,address,lng,lat,image1_url,image2_url) values(?,?,?,?,?,?,?,?,?,?,?)"
	QueryOnePlaceStr    = "select place_code,name,founder,founder_comment,score,type,phone_number,address,lng,lat,image1_url,image2_url,comment_number from cug_map_place_tpl where place_code = ?"
	QueryAllPlaceStr    = "select place_code,name,founder,founder_comment,score,type,phone_number,address,lng,lat,image1_url,image2_url,comment_number from cug_map_place_tpl where lng between ? and ? and lat between ? and ?"
	CountPlaceNumberStr = "select COUNT(place_code) from cug_map_place_tpl where lng between ? and ? and lat between ? and ?"
)

func (this *Place) AddPlace() (err error) {
	stmt, err := Db.Prepare(AddPlaceStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(this.Name, this.Founder, this.FounderComment, this.Score, this.Type, this.Number, this.Address, this.Longitude, this.Latitude, this.Image1Url, this.Image2Url)
	if err != nil {
		return err
	}
	return nil
}
func GetPlace(user *User) (placeInfoList []*Place, err error) {
	number := 0
	Lrange := 10.0
	placeInfoList = make([]*Place, 0)
	stmt2, err2 := Db.Prepare(QueryAllPlaceStr)
	defer stmt2.Close()
	if err2 != nil {
		fmt.Println("err2", err2.Error())
		return
	}
	stmt, err := Db.Prepare(CountPlaceNumberStr)
	defer stmt.Close()
	if err != nil {
		fmt.Println("err1", err.Error())
		return
	}
	row := stmt.QueryRow(-180, 180, -90, 90)
	row.Scan(&number)
	for number > 100 {
		row = stmt.QueryRow(user.Longitude-Lrange, user.Longitude+Lrange, user.Latitude-Lrange, user.Latitude+Lrange)
		row.Scan(&number)
		Lrange = Lrange / 1.2
	}
	Lrange = Lrange * 1.2
	var rows *sql.Rows

	if Lrange > 9 {
		rows, err = stmt2.Query(-180, 180, -90, 90)
	} else {
		rows, err = stmt2.Query(user.Longitude-Lrange, user.Longitude+Lrange, user.Latitude-Lrange, user.Latitude+Lrange)
	}
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	count := 0
	defer rows.Close()
	for rows.Next() {
		count++
		if count > 100 {
			break
		}
		place := Place{}
		rows.Scan(&place.PlaceCode, &place.Name, &place.Founder, &place.FounderComment, &place.Score, &place.Type, &place.Number, &place.Address, &place.Longitude, &place.Latitude, &place.Image1Url, &place.Image2Url, &place.CommentNumber)
		place.Founder = user.Username
		placeInfoList = append(placeInfoList, &place)
	}
	return
}
func GetOnePlace(placeCode string) (place *Place, err error) {
	place = &Place{
		Name:           "",
		Address:        "",
		Longitude:      0,
		Latitude:       0,
		Image1Url:      "",
		Image2Url:      "",
		Type:           0,
		Score:          0,
		CommentNumber:  0,
		Number:         "",
		PlaceCode:      "",
		Founder:        "",
		FounderComment: "",
	}
	rdb := Rdb.Rdb
	val, err := rdb.Get(context.TODO(), "P"+placeCode).Result()
	if err == nil {
		err = json.Unmarshal([]byte(val), place)
		if err != nil {

		} else {
			return place, nil
		}
	}
	stmt, err := Db.Prepare(QueryOnePlaceStr)
	if err != nil {
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(placeCode)
	err = row.Scan(&place.PlaceCode, &place.Name, &place.Founder, &place.FounderComment, &place.Score, &place.Type, &place.Number, &place.Address, &place.Longitude, &place.Latitude, &place.Image1Url, &place.Image2Url, &place.CommentNumber)
	if err != nil {
		fmt.Println("hello", err.Error())
		return
	}
	placeByte, err := json.Marshal(place)
	placeStr := string(placeByte)
	if err != nil {
		fmt.Println("marshal error", err.Error())
		return
	}
	rdb.Set(context.TODO(), "P"+placeCode, placeStr, time.Second*600)
	return
}
