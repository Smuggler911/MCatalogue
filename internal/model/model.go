package model

type Car struct {
	Id      int    `json:"id"`
	RegNum  string `json:"regNum"`
	Mark    string `json:"mark "`
	Model   string `json:"model"`
	Year    int    `json:"year,omitempty"`
	OwnerId int    `json:"owner_id"`
	Owner   People `json:"owner"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}
