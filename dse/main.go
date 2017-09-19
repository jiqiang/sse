package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("10.10.1.167")
	cluster.Port = 9042
	cluster.Keyspace = "iepmaster"
	cluster.Timeout = 15 * time.Second
	cluster.ProtoVersion = 4

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// sql := `INSERT INTO event_current(enterprise_uid,event_bucket,resrc_uid, event_ts) VALUES (?, ?, ?, ?)`
	// err = session.Query(sql, "SiteLastMessage", "SiteLastMessage", "ENTERPRISE_bc7967ad-0bbb-4804-8145-f087955ffae1.SITE_5521b291-1e8e-37d4-9dd9-87d6311305c2", time.Now().UTC()).Exec()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var t time.Time
	iter := session.Query(`select event_ts from event_current limit 1`).Iter()
	for iter.Scan(&t) {
		d := time.Now().Sub(t)
		fmt.Println(t.String(), time.Now().UTC(), d.Minutes())
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
