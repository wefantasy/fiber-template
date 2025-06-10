package main

import (
	"app/server"
)

//	@title			app api
//	@version		1.0
//	@description	app 的接口文档
//	@termsOfService	https://github.com/wefantasy/fiber-template

//	@contact.name	Fantasy
//	@contact.url	https://github.com/wefantasy/fiber-template

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/api/v1
func main() {
	s, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	if err := s.Run(); err != nil {
		panic(err)
	}
}
