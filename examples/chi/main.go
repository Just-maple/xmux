package main

import (
	"log"
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/Just-maple/xmux/examples/common/business"
)

func main() {
	controller := NewController()
	userService := business.NewUserService()

	userGroup := xmux.DefineGroup(func(r xmux.Router, svc *business.UserService) {
		xmux.Register(r, http.MethodPost, "/users", svc.CreateUser)
		xmux.Register(r, http.MethodGet, "/users", svc.ListUsers)
		xmux.Register(r, http.MethodGet, "/user", svc.GetUser)
		xmux.Register(r, http.MethodPut, "/users", svc.UpdateUser)
		xmux.Register(r, http.MethodDelete, "/users", svc.DeleteUser)
	})

	err := userGroup.Bind(controller, func(ptr any) error {
		switch p := ptr.(type) {
		case **business.UserService:
			*p = userService
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Chi server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", controller))
}
