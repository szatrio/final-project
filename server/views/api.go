package views

import "time"

// Success Response
type Response struct {
	Message string      `json:"message" example:"GET_SUCCESS"`
	Status  int         `json:"status" example:"201"`
	Payload interface{} `json:"payload,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// /////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////// Success Response For USER Table ////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////
type Data_Register struct {
	Age       int    `json:"age"`
	Email     string `json:"email"`
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Update_At time.Time
	Create_At time.Time
}

type Resp_Register_Success struct {
	Message string        `json:"message" example:"GET_SUCCESS"`
	Status  int           `json:"status" example:"201"`
	Data    Data_Register `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}
type Resp_Login struct {
	Message string `json:"message" example:"GET_SUCCESS"`
	Status  int    `json:"status" example:"201"`
	Data    Token  `json:"data"`
}

type Put struct {
	Age       int    `json:"age"`
	Email     string `json:"email"`
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Update_At time.Time
}
type Resp_Put struct {
	Message string `json:"message" example:"GET_SUCCESS"`
	Status  int    `json:"status" example:"201"`
	Data    Put    `json:"data"`
}

type Data_Delete struct {
	Message string `json:"message"`
}
type Resp_Delete struct {
	Message string      `json:"message" example:"GET_SUCCESS"`
	Status  int         `json:"status" example:"201"`
	Data    Data_Delete `json:"data"`
}

// /////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////// Success Response For PHOTO Table ////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////
type Data_Photo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	Photo_Url string `json:"photo_url"`
	User_Id   int    `json:"user_id"`
	Create_At time.Time
}

type Resp_Post_Photo struct {
	Message string     `json:"message" example:"GET_SUCCESS"`
	Status  int        `json:"status" example:"201"`
	Data    Data_Photo `json:"data"`
}
