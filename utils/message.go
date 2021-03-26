package db

import (
	"fmt"
)

const (
	insertMessageStr = "insert into cug_map_message_tpl2(receiver_id,sender_id,message,send_time,send_time_str) values(?,?,?,?,?);"
	queryReceiverStr = "select sender_id,receiver_id,message,send_time_str,send_time,is_read from cug_map_message_tpl2 where receiver_id = ?;"
	querySenderStr   = "select sender_id,receiver_id,message,send_time_str,send_time,is_read from cug_map_message_tpl2 where sender_id = ?;"

	//queryMessageStr = "select receiver,sender,message,send_time from cug_map_messages_tpl where sender = ?;"
	updateMessageStr = "update cug_map_message_tpl2 set is_read = true where sender_id = ? and receiver_id = ?;"
)

type MessageType struct {
	SenderId    string
	ReceiverId  string
	SendTime    string
	Message     string
	SendTimeStr string
	IsRead      bool
}

type MessageTypeWithName struct {
	SenderId     string
	SenderName   string
	ReceiverId   string
	ReceiverName string
	SendTime     string
	Message      string
	SendTimeStr  string
	IsRead       bool
}

func (this *MessageType) UpdateMessage() (err error) {
	stmt, err := Db.Prepare(updateMessageStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(this.SenderId, this.ReceiverId)
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

func (this *MessageType) AddMessage() (err error) {
	stmt, err := Db.Prepare(insertMessageStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(this.ReceiverId, this.SenderId, this.Message, this.SendTime, this.SendTimeStr)
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
func (this *MessageType) GetMyMessage() (ReceiverMessageList []*MessageTypeWithName, SenderMessageList []*MessageTypeWithName, err error) {
	stmt, err := Db.Prepare(queryReceiverStr)
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
		err = rows.Scan(&MyMessage.SenderId, &MyMessage.ReceiverId, &MyMessage.Message, &MyMessage.SendTimeStr, &MyMessage.SendTime, &MyMessage.IsRead)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		sender := &User{
			StudentId: MyMessage.SenderId,
		}
		senderName := sender.QueryUserName()
		receiver := &User{
			StudentId: MyMessage.ReceiverId,
		}
		receiverName := receiver.QueryUserName()
		var MyMessageWithName *MessageTypeWithName

		MyMessageWithName = &MessageTypeWithName{
			SenderId:    MyMessage.SenderId,
			SenderName:  senderName,
			ReceiverId:  MyMessage.ReceiverId,
			ReceiverName:receiverName,
			SendTime:    MyMessage.SendTime,
			Message:     MyMessage.Message,
			SendTimeStr: MyMessage.SendTimeStr,
			IsRead:      MyMessage.IsRead,
		}

		ReceiverMessageList = append(ReceiverMessageList, MyMessageWithName)
	}
	stmt, err = Db.Prepare(querySenderStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	rows, err = stmt.Query(this.ReceiverId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for rows.Next() {
		MyMessage := &MessageType{}
		err = rows.Scan(&MyMessage.SenderId, &MyMessage.ReceiverId, &MyMessage.Message, &MyMessage.SendTimeStr, &MyMessage.SendTime, &MyMessage.IsRead)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		sender := &User{
			StudentId: MyMessage.SenderId,
		}
		senderName := sender.QueryUserName()
		receiver := &User{
			StudentId: MyMessage.ReceiverId,
		}
		receiverName := receiver.QueryUserName()
		MyMessageWithName := &MessageTypeWithName{
			SenderId:    MyMessage.SenderId,
			SenderName:  senderName,
			ReceiverId:  MyMessage.ReceiverId,
			ReceiverName:receiverName,
			SendTime:    MyMessage.SendTime,
			Message:     MyMessage.Message,
			SendTimeStr: MyMessage.SendTimeStr,
			IsRead:      MyMessage.IsRead,
		}

		SenderMessageList = append(SenderMessageList, MyMessageWithName)
	}

	return
}
