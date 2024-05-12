package main

type User struct {
	Id int;
	Username string;
	HashedPassword string;
}

func NewUserAuthenticator() UserAuth {
	return UserAuth{}
}



