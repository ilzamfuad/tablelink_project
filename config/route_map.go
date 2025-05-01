package config

type RouteMapping struct {
	Route  string
	Method string
}

var GrpcToRestfulMapping = map[string]RouteMapping{
	"/UserService/GetAllUsers": {
		Route:  "/users",
		Method: "GET",
	},
	"/UserService/CreateUser": {
		Route:  "/users/user",
		Method: "POST",
	},
	"/UserService/UpdateUser": {
		Route:  "/users/user",
		Method: "PUT",
	},
	"/UserService/DeleteUser": {
		Route:  "/users/user/{user_id}",
		Method: "DELETE",
	},
}
