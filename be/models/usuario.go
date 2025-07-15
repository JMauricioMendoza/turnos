package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	ID             int    `json:"id"`
	Usuario        string `json:"usuario" validate:"min=6,max=25"`
	Password       string `json:"password,omitempty" validate:"min=6,max=25"`
	NombreCompleto string `json:"nombre_completo" validate:"min=1,max=100"`
	ActividadID    int    `json:"actividad_id"`
	RolID          int    `json:"rol_id"`
	Mesa           int    `json:"mesa" validate:"gte=0"`
}

func (u *Usuario) ValidarUsuario(Usuario) error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *Usuario) HashearPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *Usuario) VerificarPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
