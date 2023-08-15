package models

import "fmt"

type ListAll struct {
	Owner       UserBasic    `json:"owner"`
	ContactList []UserBasic  `json:"contactList"`
	OnLineList  []UserBasic  `json:"onLineList"`
	GroupList   []GroupBasic `json:"groupList"`
}

func ListInOnePage(userId string) (ListAll, error) {
	listInOne := ListAll{
		Owner: FindUserByUID(userId),
	}
	for _, each := range listInOne.ContactList {
		if IsOnline(each.UID) {
			listInOne.OnLineList = append(listInOne.OnLineList, each)
		}
	}
	listInOne.GroupList = GroupsList(userId)
	fmt.Println(GroupList(userId))
	return listInOne, nil
}
