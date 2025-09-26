package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MarDoA/gator/internal/database"
	"github.com/google/uuid"
)

func handleBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		limit, _ = strconv.Atoi(cmd.args[0])
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Println("================================================================")
		if post.Title.Valid {
			fmt.Println(post.Title.String)
		}
		if post.Description.Valid {
			fmt.Println(post.Description.String)
		}
		fmt.Println("link: " + post.Url)
	}
	return nil

}
func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Println("feed unfollowed")
	return nil
}

func handlerListFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	if len(follows) == 0 {
		fmt.Println("no feeds followed")
		return nil
	}
	for _, f := range follows {
		fmt.Println(f.FeedName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: follow <URL>")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	printFeedFollow(ff)
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: agg <duration ex: 1m>")
	}
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		fmt.Println("Collecting feeds every " + cmd.args[0])
		err = scrapFeeds(s)
		if err != nil {
			return err
		}

	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: addfeed <feed name> <feed url>")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("a feed was added")
	printFeed(feed, user)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf(" * URL:        %s\n", feed.Url)
		fmt.Printf(" * UserName:   %s\n", user.Name)
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserName:        %s\n", user.Name)
}

func printFeedFollow(ffr database.CreateFeedFollowRow) {
	fmt.Printf("* ID:            %s\n", ffr.ID)
	fmt.Printf("* Created:       %v\n", ffr.CreatedAt)
	fmt.Printf("* Updated:       %v\n", ffr.UpdatedAt)
	fmt.Printf("* User Name:          %s\n", ffr.UserName)
	fmt.Printf("* Feed Name:           %s\n", ffr.FeedName)
}
