package main

import "app/app"

//	@title			app api
//	@version		1.0
//	@description	app 的接口文档
//	@termsOfService	https://www.ifantasy.net

//	@contact.name	Fantasy
//	@contact.url	https://ifantasy.net
//	@contact.email	root@ifantasy.net

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:80
// @BasePath	/api/v1
func main() {
	s, err := app.NewServer()
	if err != nil {
		panic(err)
	}
	if err := s.Run(); err != nil {
		panic(err)
	}
}
