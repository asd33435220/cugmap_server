package db

import "fmt"

const (
	AddCommentStr = "insert into cug_map_place_comment_tpl(place_code,commentator,comment_message,score,comment_time,comment_time_str,likes) values(?,?,?,?,?,?,?);"
	GetCommentStr = "select commentator,comment_message,score,comment_time,comment_time_str,likes from cug_map_place_comment_tpl where place_code = ?"
	getPlaceCommentNumber = "select comment_number,score from cug_map_place_tpl where place_code = ?"
	updatePlaceScore = "update cug_map_place_tpl set comment_number = ?,score = ? where place_code = ? "
	)

type Comment struct {
	CommentId string
	PlaceCode string
	Commentator string
	CommentatorInfo User
	CommentMessage string
	Score int
	CommentTime string
	CommentTimeStr string
	Likes int
}

func(this *Comment)AddComment()(err error){
	tx ,err := Db.Begin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_,err = tx.Exec(AddCommentStr,this.PlaceCode,this.Commentator,this.CommentMessage,this.Score,this.CommentTime,this.CommentTimeStr,0)
	if err != nil {
		return
	}
	var score float64
	var number float64
	row:=tx.QueryRow(getPlaceCommentNumber,this.PlaceCode)
	err = row.Scan(&number,&score)
	if err != nil {
		return
	}
	newScore := (score * (number+1) + float64(this.Score))/(number+2)
	if newScore !=score {
		_,err = tx.Exec(updatePlaceScore,int(number+1),newScore,this.PlaceCode)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return

}
func GetComment(placeCode string)(commentList []*Comment,err error){
	commentList = make([]*Comment,0)
	stmt,err := Db.Prepare(GetCommentStr)
	if err != nil {
		fmt.Println("err",err)
		return
	}
	defer stmt.Close()

	rows,err := stmt.Query(placeCode)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
	for rows.Next() {
		newComment := Comment{
			PlaceCode:  placeCode,
		}
		err = rows.Scan(&newComment.Commentator,&newComment.CommentMessage,&newComment.Score,&newComment.CommentTime,&newComment.CommentTimeStr,&newComment.Likes)
		if err != nil {
			return
		}
		user := &User{
			StudentId: newComment.Commentator,
		}
		user.QueryAllInfo()
		newComment.CommentatorInfo = *user
		commentList = append(commentList,&newComment)
	}
	return
}