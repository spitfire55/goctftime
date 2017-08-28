package main

import (
	"os"

	"github.com/zabawaba99/firego"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/log"
	"fmt"
	"strconv"
)

var fb *firego.Firebase

func connect(ctx context.Context) {

	hc, err := google.DefaultClient(ctx,
		"https://www.googleapi.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
	fb = firego.New(os.Getenv("FIREBASE_BASE"), hc)
}

func saveAllRankings(teamRankings interface{}, ctx context.Context) {
	if err := fb.Child("Rankings").Set(teamRankings); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func saveCurrentRankings(teamRankings interface{}, ctx context.Context) {
	if err := fb.Child("Rankings/2017").Set(teamRankings); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func saveAllTeams(teams interface{}, ctx context.Context) {
	if err:= fb.Child("Teams").Set(teams); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func saveNewTeam(team interface{}, ctx context.Context) {
	highestNode := fmt.Sprintf("%d", fb.Child("TeamHighestNode"))
	if err := fb.Child("Teams/" + highestNode).Set(team); err != nil {
		log.Errorf(ctx, err.Error())
	}
	highestNodeInt, err := strconv.Atoi(highestNode)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
	if err := fb.Child("TeamHighestNode").Set( highestNodeInt+ 1); err != nil {
		log.Errorf(ctx, err.Error())
	}

}
