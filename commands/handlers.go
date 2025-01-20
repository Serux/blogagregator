package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/serux/blogagregator/internal/config"
	"github.com/serux/blogagregator/internal/database"
	"github.com/serux/blogagregator/rss"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	CommandsHandlers map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandsHandlers[name] = f
}
func (c *Commands) Run(s *State, cmd Command) error {
	fun, ok := c.CommandsHandlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not handled")
	}
	return fun(s, cmd)
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, c Command) error {
		user, err := s.Db.GetOneUserByName(context.Background(), s.Config.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, c, user)
	}

}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("format: unfollow url")
	}
	params := database.DeleteFeedFollowsUserIDAndURLParams{
		UserID: user.ID,
		Url:    cmd.Arguments[0],
	}
	err := s.Db.DeleteFeedFollowsUserIDAndURL(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, v := range feedFollows {
		fmt.Println(v.Feedname)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {

	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("format: follow url")
	}

	feed, err := s.Db.GetFeedByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return err
	}

	fmt.Println(feedFollow)
	return nil
}

func HandlerGetAllFeeds(s *State, cmd Command) error {

	ret, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(ret)
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("format: addfeed name url")
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}
	feed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	feedfollowparams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), feedfollowparams)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("format: agg time (1s,1m,1h...)")
	}
	time_between_reqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("Collectiong feeds every ", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func HandlerReset(s *State, cmd Command) error {

	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Println("Cannot reset users:")
		return err
	}

	fmt.Println("Reset successful")

	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("Cannot get users:")
		return err
	}

	for _, user := range users {
		current := ""
		if user.Name == s.Config.Current_user_name {
			current = "(current)"
		}
		fmt.Printf("* %v %v\n", user.Name, current)
	}

	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username required")
	}
	_, err := s.Db.GetOneUserByName(context.Background(), cmd.Arguments[0])
	if err != nil {
		fmt.Println("Cannot find user:")
		return err
	}
	s.Config.SetUser(cmd.Arguments[0])
	fmt.Println("Username Set")

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("name required")
	}

	us, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Arguments[0],
		},
	)
	if err != nil {
		fmt.Println("User exists: ")
		return err
	}

	s.Config.SetUser(cmd.Arguments[0])
	fmt.Println("User created:", us)

	return nil
}

func scrapeFeeds(s *State) error {
	nextFeed, err := s.Db.GetNextFeedToFecth(context.Background())
	if err != nil {
		fmt.Println("No feeds: ")
		return err
	}
	err = s.Db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		fmt.Println("Err marking feed: ")
		return err
	}

	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	//fmt.Println("URL: ", feed.Channel.Link)
	//fmt.Println("Title: ", feed.Channel.Title)
	fmt.Println("Description: ", feed.Channel.Description)
	fmt.Println()
	for _, v := range feed.Channel.Item {
		//fmt.Println("Date: ", v.PubDate)
		//fmt.Println("Link: ", v.Link)
		fmt.Println("Title: ", v.Title)
		//fmt.Println("Description: ", v.Description)
		//fmt.Println()
	}

	return nil
}
