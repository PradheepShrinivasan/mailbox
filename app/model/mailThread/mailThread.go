package mailThread

import (
	"errors"

	"github.com/PradheepShrinivasan/mailbox/app/model/mail"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MailThread struct {
	Id             bson.ObjectId `bson:"_id"`
	ReferenceCount int           `bson:"referenceCount"`
	Mails          []mail.Mail   `bson:"mails"`
}

// Call this
func (mailThread *MailThread) Init() {

	mailThread.Id = bson.NewObjectId()

}

// Create a new Mail thread
func CreateNewMailThread(s *mgo.Collection, mail *mail.Mail) (bson.ObjectId, bool) {

	var newMailThread MailThread

	newMailThread.Init()
	newMailThread.ReferenceCount = mail.GetReferenceCount()

	newMailThread.Mails = append(newMailThread.Mails, *mail)

	err := s.Insert(newMailThread)
	if err != nil {
		// TODO do a logging to track everything
		panic("Failed to insert the mail" + err.Error())
		return newMailThread.Id, false
	}
	return newMailThread.Id, true
}

func AddMailtoMailThread(s *mgo.Collection, threadId bson.ObjectId, mail *mail.Mail) bool {

	/* TODO increment the reference count*/
	err := s.UpdateId(threadId, bson.M{"$push": bson.M{"mails": mail}})
	if err != nil {
		return false
	}
	return true
}

func updateMailThread(s *mgo.Collection, mailThread MailThread) bool {

	err := s.UpdateId(mailThread.Id, mailThread)
	if err != nil {
		// TODO loging for stats
		return false
	}
	return true
}

func findMailThread(s *mgo.Collection, threadId bson.ObjectId) (mailThread MailThread, err error) {

	if threadId == bson.ObjectId("") {
		err = errors.New("ObjectId is empty")
		return mailThread, err
	}

	err = s.FindId(threadId).One(&mailThread)
	return mailThread, err
}

func DeleteMailThread(s *mgo.Collection, threadId bson.ObjectId) (retval bool) {

	mailThread, err := findMailThread(s, threadId)
	if err != nil {
		// log stats here
		return false
	}

	if mailThread.ReferenceCount == 1 {

		err = s.RemoveId(mailThread.Id)
		if err != nil {
			//TODO do a logging to track everything
			panic("Failed to remo the mail" + err.Error())
			return false
		}
		return true
	}

	mailThread.ReferenceCount -= 1

	return updateMailThread(s, mailThread)
}

func GetMailThread(s *mgo.Collection, threadId bson.ObjectId) (mailThread MailThread, err error) {

	return findMailThread(s, threadId)
}
