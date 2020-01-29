package dto

type ProductList struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	UpdateTime  int64  `json:"update_time"`
	CreateTime  int64  `json:"create_time"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Status      int16  `json:"status"`
	Star        int64  `json:"star"`
}
