package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Prescription struct {
	Id            int    `json:"id"`
	Patient_Name  string `json:"patient_name"`
	Medicine_name string `json:"medicine_name"`
	Dosage_form   string `json:"dosage_form"`
	Quantity      int    `json:"quantity"`
}

type Add_Medicine struct {
	Medicine_Name  string `json:"medicine_name"`
	Dosage_form    string `json:"dosage_form"`
	Stock_Quantity int    `json:"stock_quantity"`
}
