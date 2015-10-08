package mailThreadCache

import (
	"testing"

	"github.com/PradheepShrinivasan/mailbox/app/model/mail"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

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

func (s *MySuite) TestmailacheInit(c *C) {

	mail := createMail("my subject", "body")

	var myCache MailCache
	myCache.Init(&mail)

	c.Assert(myCache.MailId, Equals, mail.Id)
	c.Assert(myCache.Read, Equals, false)
	c.Assert(myCache.From, Equals, mail.From)
	c.Assert(myCache.Subject, Equals, mail.Subject)
	c.Assert(myCache.TimeStamp, Equals, mail.TimeStamp)
}

func (s *MySuite) TestAddMailToMailThreadCache(c *C) {

	mail := createMail("my subject", "body")
	mailThreadId := bson.NewObjectId()

	var mailThreadCache MailThreadCache
	mailThreadCache.AddMailToMailThreadCache(mailThreadId, &mail)

	var myCache MailCache
	myCache.Init(&mail)
	c.Assert(mailThreadCache.ThreadId, Equals, mailThreadId)
	c.Assert(mailThreadCache.MailCaches[0], DeepEquals, myCache)
}
