package model

type Architecture struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Shellcodes []Shellcode `json:"shellcodes"`
}

type Shellcode struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	DatePublished string `json:"date_published"`
	Data          string `json:"data"`
}
