package main

import (
	"context"
	"github.com/kobayashicomputing/gator/internal/database"
)

/***************************************************************************
**
**	This middleware handles checking for the current user in the database 
**	and getting the appropriate user ID
**
**	It takes a 'middleware aware' function and returns an anonymous 
**	function that calls it. The anonymous function is a 'normal' 
**	handler for this app.
**
***************************************************************************/
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.ptrConfig.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
