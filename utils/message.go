package db

import (
	"fmt"
)

const (
	insertMessageStr  = "insert into cug_map_messages_tpl(receiver,sender,message,send_time) values(?,?,?,?);"
	queryMyMessageStr = "select receiver,sender,message,send_time from cug_map_messages_tpl where receiver = ?;"
)

type MessageType struct {
	SenderId   string
	ReceiverId string
	SendTime   string
	Message    string
}

type MessageTypeWithName struct {
	SenderId   string
	SenderName string
	ReceiverId string
	SendTime   string
	Message    string
}

func (this *MessageType) AddMessage() (err error) {
	stmt, err := Db.Prepare(insertMessageStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(this.ReceiverId, this.SenderId, this.Message, this.SendTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("rows", rows)
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("id", id)
	return nil
}
func (this *MessageType) GetMyMessage() (MessageList []*MessageTypeWithName, err error) {
	stmt, err := Db.Prepare(queryMyMessageStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(this.ReceiverId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for rows.Next() {
		MyMessage := &MessageType{}
		err = rows.Scan(&MyMessage.ReceiverId, &MyMessage.SenderId, &MyMessage.Message, &MyMessage.SendTime)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		newUser := &User{
			StudentId: MyMessage.SenderId,
		}
		username := newUser.QueryUserName()
		MyMessageWithName := &MessageTypeWithName{
			SenderId:   MyMessage.SenderId,
			SenderName: username,
			ReceiverId: MyMessage.ReceiverId,
			SendTime:   MyMessage.SendTime,
			Message:    MyMessage.Message,
		}
		MessageList = append(MessageList, MyMessageWithName)
	}
	return
}
