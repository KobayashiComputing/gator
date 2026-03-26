package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kobayashicomputing/gator/internal/config"
	"github.com/kobayashicomputing/gator/internal/database"
)

//
//=====[ These handlers use the 'middlewareLoggedIn' functionality ]=====
//

// handlerBrowse - browse the latest feeds in the given user's follow list
func handlerBrowse(s *state, cmd command, user database.User) error {
	// handle any parameters passed in...
	pLimit := 2

	if len(cmd.args) > 1 {
		return errors.New("Gator: (browse command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	if len(cmd.args) == 1 {
		theLimit, err := strconv.Atoi(cmd.args[0])
		if err == nil {
			pLimit = theLimit
		} else {
			pLimit = theLimit
		}
	}

	fmt.Println("Going to show at most " + strconv.Itoa(pLimit) + " posts...")

	// make sure the user is logged in and valid
	user, err := s.db.GetUser(context.Background(), s.ptrConfig.CurrentUserName)
	if err != nil {
		return errors.New("Gator: (browse command) error - user '" + s.ptrConfig.CurrentUserName + "' not found in database...")
	}

	// create the arg param struct to use
	arg := database.GetPostsForUserParams{
		UsersID: user.ID,
		Limit: int32(pLimit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), arg)
	if err != nil {
		return errors.New("Gator: (browse command) error - feed follow record for '" + s.ptrConfig.CurrentUserName + "' and '" + cmd.args[0] + "' cound not be deleted...")
	}

	for _, post := range(posts) {
		fmt.Println(post.Url + ": " + post.Title)
	}


	return nil
}

// handlerUnfollow - unfollow a given feed for the current user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return errors.New("Gator: (unfollow command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	user, err := s.db.GetUser(context.Background(), s.ptrConfig.CurrentUserName)
	if err != nil {
		return errors.New("Gator: (unfollow command) error - user '" + s.ptrConfig.CurrentUserName + "' not found in database...")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return errors.New("Gator: (unfollow command) error - feed '" + cmd.args[0] + "' not found in database...")
	}

	// create the arg param struct to use
	arg := database.DeleteFeedFollowParams{
		UsersID: user.ID,
		FeedsID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), arg)
	if err != nil {
		return errors.New("Gator: (unfollow command) error - feed follow record for '" + s.ptrConfig.CurrentUserName + "' and '" + cmd.args[0] + "' cound not be deleted...")
	}

	return nil 

}

// handlerFollowing - print all of the names of the feeds that the current user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return errors.New("Gator: (following command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	// get a list of the feeds that the current user is following
	feedList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("Gator: (following command) warning - no feeds followed by '" + s.ptrConfig.CurrentUserName + "'...")
	}

	// fmt.Println(feedList)
	for _, feed := range(feedList) {
		fmt.Println(feed.FeedName)
	}

	return nil 
}

// handlerFollow - create a new feed_follow record for the given URL and current user
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return errors.New("Gator: (follow command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	// get the URL of the desired feed from the function parameters
	url := cmd.args[0]
	// get the UUID for this feed, err if it does not exist
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Println("Gator: (follow command) error - feed for '" + url + "' not found in database...")
		os.Exit(1)
	}

	// create the arg param struct to use
	arg := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UsersID: user.ID,
		FeedsID: feed.ID,
	}

	row, err := s.db.CreateFeedFollow(context.Background(), arg)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "duplicate key value") {
			fmt.Println("Gator: (follow command) warning - '" + s.ptrConfig.CurrentUserName + "' is already following feed '" + url + "'...")
		} else {
			fmt.Println("Gator: (follow command) error - could not create 'feed_follows' record for ...")
		}
	} else {
		fmt.Println("New Row: ")
		fmt.Println(row)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Usage: gator addfeed <name> <url>")
	}

	// get the params
	feedURL := cmd.args[1]
	feedName := cmd.args[0]
	
	arg := database.CreateFeedParams{
		ID: uuid.New(), 
		CreatedAt: time.Now(), 
		UpdatedAt: time.Now(), 
		Url: feedURL, 
		Name: feedName, UsersID: user.ID,
	}

	newFeed, err := s.db.CreateFeed(
		context.Background(), arg)
	
	if err != nil {
		return errors.New("Gator: (addfeed command) could not add feed named '" + feedName + "' at URL '" + feedURL + "'...")
	}

	// test the new 'mark_feed_fetched' query by marking the new feed fetched...
	newFeed, err = s.db.MarkFeedFetched(context.Background(), newFeed.ID)

	fmt.Println(newFeed)

	// add a feed_follows record for this feed and the current user
	err = handlerFollow(s, command{"follow", cmd.args[1:]}, user)
	
	return nil
}

//
//=====[ These handlers DO NOT use the 'middlewareLoggedIn' functionality...]
//
// handlerFeeds - list all feeds (name, url, and user who added it)
func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("Gator: (feeds command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	feedList, err := s.db.GetFeedList(context.Background())
	if err != nil {
		return errors.New("Gator: (feeds command) error - not successful...")
	}

	userName := ""
	for _, f := range feedList {
		dbUser, err := s.db.GetUserByID(context.Background(), f.UsersID)
		if err != nil {
			userName = "unknown"
		} else {
			userName = dbUser.Name
		}
		fmt.Println(f.Name, f.Url, userName)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Usage: gator agg <time_between_reqs> (format such as 1s 1m 1h 1m30s)")
	}

	time_between_reqs := cmd.args[0]
	fmt.Println("Collecting feeds every " + time_between_reqs + "... ")
	fmt.Println("")

	// run the scrapeFeeds funtion every 'time_between_reqs' period...
	scrapeDelay, err := time.ParseDuration(time_between_reqs)
	// fmt.Println("scrapeDelay is " + scrapeDelay.String())
	if err != nil {
		fmt.Println("...could not convert '" + time_between_reqs + "' to a time duration...")
		return err
	} else {
		ticker := time.NewTicker(scrapeDelay)
		// ticker := time.NewTicker(time.Second * 15)
		for ; ; <-ticker.C {
			scrapeFeeds(s)
		}	
	}
	
	// return nil 
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Usage: gator login <username>")
	}

	userName := cmd.args[0]

	dbUser, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		fmt.Println("Gator: (login command) error - user '" + userName + "' not found in database...")
		os.Exit(1)
	}

	if config.SetUserName(*s.ptrConfig, dbUser.Name) != nil {
		return errors.New("Gator: (login command) Could not set username to '" + dbUser.Name + "': ")
	}

	fmt.Println("Current user set to '" + dbUser.Name + "'...")


	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Usage: gator register <username>")
	}

	userName := cmd.args[0]

	// check to see if the user already exists...
	existingUser, err := s.db.GetUser(context.Background(), userName)
	if err == nil {
		if existingUser.Name == userName {
			fmt.Println("Gator: A user with the name '" + existingUser.Name + "' already exists... exiting...")
		} else {
			fmt.Println("Gator: There was a problem with the database... exiting...")
		}
		os.Exit(1)
	}

	// if not, create the new user...
	dbUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: userName})
	if err != nil {
		fmt.Println("Gator: cannot create a new user with name '" + userName + "'... exiting...")
		os.Exit(1)
	}

	if config.SetUserName(*s.ptrConfig, dbUser.Name) != nil {
		return errors.New("Gator: (register command) Could not set username to '" + dbUser.Name + "': ")
	}

	fmt.Println("Gator: user '" + dbUser.Name + "' created and set as current user in '.gatorconfig.json' file...")
	// fmt.Println("Gator: (debug) newly created user is... ")
	// fmt.Print(dbUser)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("Gator: (users command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	dbUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return errors.New("Gator: (reset command) error - not successful...")
	}

	// get the current username
	currentUserName := s.ptrConfig.CurrentUserName

	for _, u := range dbUsers {
		if u.Name == currentUserName {
			fmt.Println("* " + u.Name + " (current)")
		} else {
			fmt.Println("* " + u.Name)
		}
	}

	fmt.Println("Gator: (reset command) database reset successful...")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("Gator: (reset command) error - extraneous input found on command line '" + cmd.args[0] + "'...")
	}

	err := s.db.Reset(context.Background())
	if err != nil {
		return errors.New("Gator: (reset command) error - not successful...")
	}

	if config.SetUserName(*s.ptrConfig, "") != nil {
		return errors.New("Gator: (register command) Could not set username to ''... ")
	}

	fmt.Println("Gator: (reset command) database reset successful...")
	return nil
}
