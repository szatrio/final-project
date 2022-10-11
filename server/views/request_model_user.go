package views

import "time"

type Request_Register struct {
	Id_Number int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	Create_At time.Time
	Update_At time.Time
}

type Request_Photos struct {
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	Photo_Url string `json:"photo_url`
}

type Request_Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
