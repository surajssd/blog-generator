package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	date := flags.String("date", "", "Provide the date of the event in the form 'DAY MON DD' e.g. Sat Oct 27")
	flags.Parse(os.Args[1:])

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" || *date == "" {
		log.Fatal("Consumer key/secret and Access token/secret required. " +
			"Provide flags: --consumer-key --consumer-secret --access-token " +
			"--access-secret --date")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	log.Printf("All clients ready")

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		log.Fatalf("Error verifying the user, %v", err)
	}
	log.Println("Credentials verified")

	// download all the recent tweets
	t := &twitter.UserTimelineParams{
		ScreenName:      "k8sblr",
		Count:           400,
		ExcludeReplies:  boolPtr(true),
		IncludeRetweets: boolPtr(true),
	}
	tweets, _, err := client.Timelines.UserTimeline(t)
	if err != nil {
		log.Fatal(err)
	}

	var IDs []string
	// filter all the tweets based on the date
	stopNowOnDifferentDate := false
	sameDateFound := false
	for _, tweet := range tweets {
		if strings.HasPrefix(strings.ToLower(tweet.CreatedAt), strings.ToLower(*date)) {
			sameDateFound = true
		}
		if sameDateFound {
			stopNowOnDifferentDate = true
		}
		if stopNowOnDifferentDate && !strings.HasPrefix(strings.ToLower(tweet.CreatedAt), strings.ToLower(*date)) {
			break
		}
		log.Println("Tweet text:", tweet.Text, "| Tweet Date:", tweet.CreatedAt, "| Tweet ID:", tweet.IDStr)
		IDs = append(IDs, tweet.IDStr)
	}
	IDs = reverseIDs(IDs)

	log.Println(generateBlog(IDs))
}

func reverseIDs(IDs []string) []string {
	var rev []string
	for i := len(IDs) - 1; i >= 0; i-- {
		rev = append(rev, IDs[i])
	}
	return rev
}

func boolPtr(b bool) *bool {
	return &b
}

func generateBlog(ids []string) string {
	ret := template
	for _, id := range ids {
		ret = fmt.Sprintf("%s\n{{< tweet %s >}}", ret, id)
	}
	return ret
}

var template = `
+++
author = "Suraj Deshmukh"
title = "Kubernetes Bangalore MONTH YEAR Event Report"			# CHANGE ME
description = "Event Report for Kubernetes Bangalore Meetup"
date = "2019-MM-DDT01:00:51+05:30" 								# CHANGE ME
categories = ["event_report"]
tags = ["kubernetes", "meetup"]
+++

The Kubernetes Bangalore Meetup was organized at ADD LOCATION on ADD DATE. The agenda for the meetup was following:

ADD AGENDA

The moments from Meetup:


`
