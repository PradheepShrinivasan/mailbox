package mail

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Mail struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	From bson.ObjectId   `bson:"from,omitempty"`
	To   []bson.ObjectId `bson:"To,omitempty"`
	Cc   []bson.ObjectId `bson:"cc,omitempty"`
	Bcc  []bson.ObjectId `bson:"bcc,omitempty"`

	TimeStamp time.Time `bson:"time"`

	Subject  string `bson:"subject"`
	Mailtext string `bson:"mailtext"`
}

func (mail *Mail) Init() {

	mail.Id = bson.NewObjectId()
	mail.TimeStamp = time.Now()
}

func (mail *Mail) GetReferenceCount() int {

	var numberOfReference = 0

	if mail.From != bson.ObjectId("") {
		numberOfReference++
	}

	numberOfReference += len(mail.To) + len(mail.Cc) + len(mail.Bcc)
	return numberOfReference
}
