package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/kobayashicomputing/gator/internal/config"
	"github.com/kobayashicomputing/gator/internal/database"
)

type state struct {
	ptrConfig *config.Config
	db *database.Queries
}

func main() {
	dbArgs := os.Args
	if len(os.Args) < 2 {
		fmt.Println("Gator: No command supplied... exiting...")
		os.Exit(1)
	}

	myConfig, err := config.ReadConfigFile()
	if err != nil {
		println("Error getting file '.gatorconfig.json' from user's home directory...")
		os.Exit(2)
	}

	// at this point we have the info from the config file, which includes our 
	// database url which contains the connection string...
	dbURL := myConfig.Db_URL
	
	// open a connection to the database
	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)



	s := state{&myConfig, dbQueries}
	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	cmdName := dbArgs[1]
	cmdArgs := dbArgs[2:]

	err = cmds.run(&s, command{cmdName, cmdArgs})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func scrapeFeeds(s *state) error {
	// get the next feed to fetch...
	f, err := s.db.GetNextFeedToFetchSingle(context.Background())
	if err != nil {
		fmt.Println("Error getting the single feed to fetch...")
	} else {
		s.db.MarkFeedFetched(context.Background(), f.ID)
		fmt.Println("Feed '" + f.Name + "' (" + f.Url + ") last fetched at: " + f.LastFetchedAt.Time.Format("2006-01-02 15:04:05"))
	}

	// fetch the feed...
	ptrFeedData, err := fetchFeed(context.Background(), f.Url)
	if err != nil {
		fmt.Println(err)
		return errors.New("Gator: (scrape) could not fetch feed data for '" + f.Name + "' (at " + f.Url + ")")
	}

	// save the data to the 'posts' table...
	for _, item := range(ptrFeedData.Channel.Item) {
		pTime, err := parsePubTime(item.PubDate)
		if err != nil {
			pTime = time.Now()
		}

		arg := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: item.Description,
			PublishedAt: pTime,
			FeedsID: f.ID,
		}

		_, err = s.db.CreatePost(context.Background(), arg)
		if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
			fmt.Println("'scrapeFeeds' could not create row in 'post' table: " + err.Error())
			return err
		}
	}

	return nil
}

// parsePubTime tries multiple layouts to parse a blog post's publication time string
// I got this function from Bing.com's AI, "Copilot"
func parsePubTime(pubTimeStr string) (time.Time, error) {
	// Common time layouts for blog posts
	layouts := []string{
		time.RFC3339,           // "2026-03-20T14:30:00Z"
		"2006-01-02 15:04:05",  // "2026-03-20 14:30:00"
		"2006-01-02",           // "2026-03-20"
		time.RFC1123Z,          // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,           // "Mon, 02 Jan 2006 15:04:05 MST"
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, pubTimeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse pubTime: %s", pubTimeStr)
}
