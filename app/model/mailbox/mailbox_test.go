package mailbox

import (
	"testing"

	"github.com/PradheepShrinivasan/mailbox/app/model/mail"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const mailBox = "Mailbox"
const userCache = "userCache"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	session *mgo.Session
	coll    *mgo.Collection
	mail    mail.Mail
	mail1   mail.Mail
}

func (s *MySuite) SetUpSuite(c *C) {

	var err error
	s.session, err = mgo.Dial("localhost")
	c.Assert(err, IsNil)
	s.coll = s.session.DB(mailBox).C(userCache)
	s.mail = createMail("Subject1", "my messages")
	s.mail1 = createMail("Subject2", "my messages2")

}

func (s *MySuite) TearDownSuite(c *C) {

	s.session.Close()
}

var _ = Suite(&MySuite{})

func createMail(subject string, body string) (mail mail.Mail) {

	var i int

	mail.Init()

	mail.From = bson.NewObjectId()
	for i = 0; i < 5; i++ {
		mail.To = append(mail.To, bson.NewObjectId())
	}

	for i = 0; i < 5; i++ {
		mail.Cc = append(mail.Cc, bson.NewObjectId())
	}

	for i = 0; i < 5; i++ {
		mail.Bcc = append(mail.Bcc, bson.NewObjectId())
	}

	mail.Subject = subject

	mail.Mailtext = body

	return mail
}

func (s *MySuite) TestAddNewThreadToInbox(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)

	ret := mailbox.AddNewThreadToInbox(s.coll, threadId, &s.mail)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 1)
	c.Assert(len(mBox.Sent), Equals, 0)
	c.Assert(len(mBox.Important), Equals, 0)
	c.Assert(len(mBox.Deleted), Equals, 0)

	//c.Assert(mBox.Inbox[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
}

func (s *MySuite) TestAddNewThreadToSent(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)

	ret := mailbox.AddNewThreadToSent(s.coll, threadId, &s.mail)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 0)
	c.Assert(len(mBox.Sent), Equals, 1)
	c.Assert(len(mBox.Important), Equals, 0)
	c.Assert(len(mBox.Deleted), Equals, 0)
	//c.Assert(mBox.Sent[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
}

func (s *MySuite) TestAddNewThreadToImportant(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)

	ret := mailbox.AddNewThreadToImportant(s.coll, threadId, &s.mail)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 0)
	c.Assert(len(mBox.Sent), Equals, 0)
	c.Assert(len(mBox.Important), Equals, 1)
	c.Assert(len(mBox.Deleted), Equals, 0)
	//c.Assert(mBox.Important[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
}

func (s *MySuite) TestAddNewThreadToDeleted(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)

	ret := mailbox.AddNewThreadToDeleted(s.coll, threadId, &s.mail)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 0)
	c.Assert(len(mBox.Sent), Equals, 0)
	c.Assert(len(mBox.Important), Equals, 0)
	c.Assert(len(mBox.Deleted), Equals, 1)
	//c.Assert(mBox.Deleted[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
}

func (s *MySuite) TestAddNewMailToInbox(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)
	mailbox.AddNewThreadToInbox(s.coll, threadId, &s.mail)

	ret := mailbox.AddMailToInbox(s.coll, threadId, &s.mail1)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 1)
	c.Assert(len(mBox.Inbox[0].MailCaches), Equals, 2)
	c.Assert(len(mBox.Sent), Equals, 0)
	c.Assert(len(mBox.Important), Equals, 0)
	c.Assert(len(mBox.Deleted), Equals, 0)

	// c.Assert(mBox.Inbox[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
	// c.Assert(mBox.Inbox[1], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail1))

}

func (s *MySuite) TestAddNewMailToSent(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)
	mailbox.AddNewThreadToSent(s.coll, threadId, &s.mail)

	ret := mailbox.AddMailToSent(s.coll, threadId, &s.mail1)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 0)
	c.Assert(len(mBox.Sent), Equals, 1)
	c.Assert(len(mBox.Sent[0].MailCaches), Equals, 2)
	c.Assert(len(mBox.Important), Equals, 0)
	c.Assert(len(mBox.Deleted), Equals, 0)

	// c.Assert(mBox.Sent[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
	// c.Assert(mBox.Sent[1], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail1))

}

func (s *MySuite) TestAddNewMailToImportant(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)
	mailbox.AddNewThreadToImportant(s.coll, threadId, &s.mail)

	ret := mailbox.AddMailToImportant(s.coll, threadId, &s.mail1)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 0)
	c.Assert(len(mBox.Sent), Equals, 0)
	c.Assert(len(mBox.Important), Equals, 1)
	c.Assert(len(mBox.Important[0].MailCaches), Equals, 2)
	c.Assert(len(mBox.Deleted), Equals, 0)

	// c.Assert(mBox.Sent[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
	// c.Assert(mBox.Sent[1], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail1))

}

func (s *MySuite) TestAddNewMailToDeleted(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)
	mailbox.AddNewThreadToDeleted(s.coll, threadId, &s.mail)

	ret := mailbox.AddMailToDeleted(s.coll, threadId, &s.mail1)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(len(mBox.Inbox), Equals, 0)
	c.Assert(len(mBox.Sent), Equals, 0)
	c.Assert(len(mBox.Important), Equals, 0)
	c.Assert(len(mBox.Deleted), Equals, 1)
	c.Assert(len(mBox.Deleted[0].MailCaches), Equals, 2)

	// c.Assert(mBox.Sent[0], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail))
	// c.Assert(mBox.Sent[1], DeepEquals, mailThreadCache.CreateMailThreadCache(threadId, &s.mail1))

}

func (s *MySuite) TestIncrementUnreadEmails(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	s.coll.Insert(mailbox)
	mailbox.AddNewThreadToDeleted(s.coll, threadId, &s.mail)

	ret := mailbox.IncrementUnreadEmail(s.coll, 1)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(mBox.NumberOfUnreadMails, Equals, 1)
}

func (s *MySuite) TestDecrementUnreadEmails(c *C) {

	var mailbox MailBox
	threadId := bson.NewObjectId()
	mailbox.UserId = bson.NewObjectId()
	mailbox.NumberOfUnreadMails = 2
	s.coll.Insert(mailbox)
	mailbox.AddNewThreadToDeleted(s.coll, threadId, &s.mail)

	ret := mailbox.DecrementUnreadEmail(s.coll, 1)

	mBox, _ := FindMailBox(s.coll, mailbox.UserId)
	c.Assert(ret, Equals, true)
	c.Assert(mBox.NumberOfUnreadMails, Equals, 1)
}
