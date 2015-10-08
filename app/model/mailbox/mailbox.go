package mailbox

import (
	"log"

	"github.com/PradheepShrinivasan/mailbox/app/model/mail"
	"github.com/PradheepShrinivasan/mailbox/app/model/mailThreadCache"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MailBox struct {
	UserId bson.ObjectId `bson:"_id"`

	NumberOfUnreadMails int `bson:"numberOfUnreadMails"`

	Inbox     []mailThreadCache.MailThreadCache `bson:"Inbox"`
	Sent      []mailThreadCache.MailThreadCache `bson:"Sent"`
	Important []mailThreadCache.MailThreadCache `bson:"Important"`
	Deleted   []mailThreadCache.MailThreadCache `bson:"Deleted"`
}

func FindMailBox(s *mgo.Collection, userId bson.ObjectId) (mbox MailBox, ret bool) {

	err := s.FindId(userId).One(&mbox)
	if err != nil {
		log.Fatal("cannot find the mailbox of user due to error " + err.Error())
		return mbox, false
	}
	return mbox, true
}

func (mbox *MailBox) storeMailThread(s *mgo.Collection, folder string, threadId bson.ObjectId, m *mail.Mail) bool {

	mTCache := mailThreadCache.CreateMailThreadCache(threadId, m)
	err := s.UpdateId(mbox.UserId, bson.M{"$push": bson.M{folder: mTCache}})
	if err != nil {
		log.Fatal("error message is " + err.Error())
		return false
	} else {
		return true
	}

}

func (mbox *MailBox) storeMail(s *mgo.Collection, folder string, threadId bson.ObjectId, m *mail.Mail) bool {

	// the folder id here is inbox.id or Important.id etc
	// in short we are trying to generate this command but generalizing it
	//db.getCollection("userCache").update({"_id": ObjectId("55ef74610057a57bb1000013"),"Deleted._id":ObjectId("55ef74610057a57bb1000012")}, {"$push": {"Deleted.$.mailCache":{"_read":false}}})

	folderId := folder + "._id"
	folderToUpdate := folder + ".$." + "mailCache"
	mCache := mailThreadCache.CreateMailCache(m)

	err := s.Update(bson.M{"_id": mbox.UserId, folderId: threadId}, bson.M{"$push": bson.M{folderToUpdate: mCache}})
	if err != nil {
		log.Fatal("error message is " + err.Error())
		return false
	} else {
		return true
	}

}

func (mbox *MailBox) IncrementUnreadEmail(s *mgo.Collection, val int) bool {

	err := s.UpdateId(mbox.UserId, bson.M{"$inc": bson.M{"numberOfUnreadMails": val}})
	if err != nil {
		log.Fatal("error message is " + err.Error())
		return false
	} else {
		return true
	}
	return true
}

func (mbox *MailBox) DecrementUnreadEmail(s *mgo.Collection, val int) bool {

	val *= -1
	return mbox.IncrementUnreadEmail(s, val)
}

func (mbox *MailBox) AddNewThreadToInbox(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMailThread(s, "Inbox", threadId, m)
}

func (mbox *MailBox) AddNewThreadToSent(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMailThread(s, "Sent", threadId, m)
}

func (mbox *MailBox) AddNewThreadToImportant(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMailThread(s, "Important", threadId, m)
}

func (mbox *MailBox) AddNewThreadToDeleted(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMailThread(s, "Deleted", threadId, m)
}

func (mbox *MailBox) AddMailToInbox(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMail(s, "Inbox", threadId, m)
}

func (mbox *MailBox) AddMailToSent(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMail(s, "Sent", threadId, m)
}

func (mbox *MailBox) AddMailToImportant(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMail(s, "Important", threadId, m)
}

func (mbox *MailBox) AddMailToDeleted(s *mgo.Collection, threadId bson.ObjectId, m *mail.Mail) bool {
	return mbox.storeMail(s, "Deleted", threadId, m)
}
