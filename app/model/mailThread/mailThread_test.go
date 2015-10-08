package mailThread

import (
	"testing"

	"github.com/PradheepShrinivasan/mailbox/app/model/mail"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const mailBox = "Mailbox"
const threadsCollection = "MailThreads"

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

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestMailThreadInit(c *C) {

	var mailThread MailThread

	mailThread.Init()

	c.Assert(mailThread.Id, Not(Equals), bson.ObjectId(""))
	c.Assert(mailThread.ReferenceCount, Equals, 0)
	c.Assert(len(mailThread.Mails), Equals, 0)
}

func (s *MySuite) TestCreateEmptyMailThread(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	var mail mail.Mail
	mail.Init()
	threadId, retval := CreateNewMailThread(coll, &mail)
	c.Assert(retval, Equals, true)

	var mailThread MailThread
	err = coll.FindId(threadId).One(&mailThread)
	c.Assert(err, IsNil)
	c.Assert(mailThread.Id, Equals, threadId)
	c.Assert(mailThread.ReferenceCount, Equals, 0)
	c.Assert(len(mailThread.Mails), Equals, 1)

	// TODO figure out why nsec is not same :(
	//c.Assert(mailThread.Mails[0], DeepEquals, mail)
}

func (s *MySuite) TestCreateMailThread(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	mail := createMail("Subject1", "Mail text 1")
	threadId, retval := CreateNewMailThread(coll, &mail)
	c.Assert(retval, Equals, true)

	var mailThread MailThread
	err = coll.FindId(threadId).One(&mailThread)
	c.Assert(err, IsNil)
	c.Assert(mailThread.Id, Equals, threadId)
	c.Assert(mailThread.ReferenceCount, Equals, 16)
	c.Assert(len(mailThread.Mails), Equals, 1)
	// TODO figure out why nsec is not same :(
	//c.Assert(mailThread.Mails[0], DeepEquals, mail)
}

func (s *MySuite) TestAddMailToMailThread(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	mail := createMail("Subject1", "Mail text 1")
	newMail := createMail("Subject1", "NewMailText")
	threadId, _ := CreateNewMailThread(coll, &mail)
	retval := AddMailtoMailThread(coll, threadId, &newMail)
	c.Assert(retval, Equals, true)

	var mailThread MailThread
	err = coll.FindId(threadId).One(&mailThread)
	c.Assert(err, IsNil)
	c.Assert(mailThread.Id, Equals, threadId)
	c.Assert(mailThread.ReferenceCount, Equals, 16)
	c.Assert(len(mailThread.Mails), Equals, 2)
	// TODO figure out why nsec is not same :(
	//c.Assert(mailThread.Mails[0], DeepEquals, mail)
	//c.Assert(mailThread.Mails[1], DeepEquals, newMail)

}

func (s *MySuite) TestFindMailThreadWithEmptyThreadId(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	_, result := findMailThread(coll, bson.ObjectId(""))

	c.Assert(result.Error(), Equals, "ObjectId is empty")

}

func (s *MySuite) TestFindMailThreadWithNotFoundThreadId(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	_, result := findMailThread(coll, bson.NewObjectId())

	c.Assert(result.Error(), NotNil)
}

func (s *MySuite) TestRemoveMailThreadWithSingleReference(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	var mail mail.Mail
	mail.Init()
	mail.From = bson.NewObjectId()
	mail.Subject = "Hello world"
	mail.Mailtext = "Long mail text"
	mailThreadId, errr := CreateNewMailThread(coll, &mail)
	c.Assert(errr, Equals, true)
	retval := DeleteMailThread(coll, mailThreadId)

	c.Assert(retval, Equals, true)
	_, err = findMailThread(coll, mailThreadId)
	c.Assert(err, NotNil)

}

func (s *MySuite) TestRemoveMailThreadWithMultipleReference(c *C) {

	session, err := mgo.Dial("localhost")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB(mailBox).C(threadsCollection)

	mail := createMail("Subject1", "Some random text")
	mailThreadId, errr := CreateNewMailThread(coll, &mail)
	c.Assert(errr, Equals, true)
	retval := DeleteMailThread(coll, mailThreadId)

	c.Assert(retval, Equals, true)
	retrievedMailThread, err := findMailThread(coll, mailThreadId)
	c.Assert(err, IsNil)
	c.Assert(retrievedMailThread.Id, Equals, mailThreadId)
	// we know that we created 16 references as using createMail
	c.Assert(retrievedMailThread.ReferenceCount, Equals, 15)
}
