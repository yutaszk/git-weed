package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/motemen/go-gitconfig"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)


func Action(_ *cli.Context) error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	headHash, err := r.Head()
	if err != nil {
		return err
	}
	headCommit, err := r.CommitObject(headHash.Hash())
	if err != nil {
		return err
	}

	committedAt, err := editTimestamp(headCommit)
	if err != nil {
		return err
	}

	userName, err := gitconfig.Default.GetString("user.name")
	if err != nil {
		return err
	}
	userEmail, err := gitconfig.Default.GetString("user.email")
	if err != nil {
		return err
	}

	_, err = w.Commit(commitMessage, &git.CommitOptions{
		Committer: &object.Signature{
			Name:  userName,
			Email: userEmail,
			When:  committedAt,
		},
		Author: &object.Signature{
			Name:  userName,
			Email: userEmail,
			When:  committedAt,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func editTimestamp(head *object.Commit) (time.Time, error) {
	committedAt := time.Now()
	if random {
		randMax := time.Now().Unix() - head.Author.When.Unix()
		if randMax < 0 {
			return committedAt, errors.New("HEAD is future commit.")
		}
		committedAt = time.Unix(time.Now().Unix() - rand.Int63n(randMax), 0)
	}

	if year == 0 {
		year = committedAt.Year()
	}
	if month == 0 {
		month = int(committedAt.Month())
	}
	if day == 0 {
		day = committedAt.Day()
	}
	if hour == 0 {
		hour = committedAt.Hour()
	}
	if minute == 0 {
		minute = committedAt.Minute()
	}
	if second == 0 {
		second = committedAt.Second()
	}
	committedAt = time.Date(year, time.Month(month), day, hour, minute, second, 0, committedAt.Location())

	committedAt = committedAt.AddDate(-yearsAgo, -monthsAgo, -daysAgo)
	committedAt = committedAt.Add(-time.Hour * time.Duration(hoursAgo))
	committedAt = committedAt.Add(-time.Minute * time.Duration(minutesAgo))
	committedAt = committedAt.Add(-time.Second * time.Duration(secondsAgo))

	if strict && committedAt.After(time.Now()) {
		return time.Now(), errors.New("Not allow future commit: " + committedAt.String())
	}

	return committedAt, nil
}
