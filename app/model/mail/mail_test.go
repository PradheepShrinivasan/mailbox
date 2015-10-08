package mail

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func createMail(subject string, body string) (mail Mail) {

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

func (s *MySuite) TestMailInit(c *C) {

	var mail Mail

	mail.Init()

	c.Assert(mail.Id, NotNil)
	c.Assert(mail.From, Equals, bson.ObjectId(""))
	c.Assert(len(mail.To), Equals, 0)
	c.Assert(len(mail.Cc), Equals, 0)
	c.Assert(len(mail.Bcc), Equals, 0)
	c.Assert(mail.TimeStamp, NotNil)
	c.Assert(len(mail.Subject), Equals, 0)
	c.Assert(len(mail.Mailtext), Equals, 0)

}

func (s *MySuite) TestGetReferenceWithNoElements(c *C) {

	var mail Mail

	mail.Init()

	c.Assert(mail.GetReferenceCount(), Equals, 0)
}

func (s *MySuite) TestGetReferenceWithFrom(c *C) {

	var mail Mail

	mail.Init()
	mail.From = bson.NewObjectId()

	c.Assert(mail.GetReferenceCount(), Equals, 1)
}

func (s *MySuite) TestGetReferenceWithSingleTo(c *C) {

	var mail Mail

	mail.Init()
	mail.To = append(mail.To, bson.NewObjectId())

	c.Assert(len(mail.To), Equals, 1)
	c.Assert(mail.GetReferenceCount(), Equals, 1)
}

func (s *MySuite) TestGetReferenceWithMultipleTo(c *C) {

	var mail Mail

	mail.Init()
	mail.To = append(mail.To, bson.NewObjectId())
	mail.To = append(mail.To, bson.NewObjectId())

	c.Assert(len(mail.To), Equals, 2)
	c.Assert(mail.GetReferenceCount(), Equals, 2)
}

func (s *MySuite) TestGetReferenceWithSingleCc(c *C) {

	var mail Mail

	mail.Init()
	mail.Cc = append(mail.Cc, bson.NewObjectId())

	c.Assert(len(mail.Cc), Equals, 1)
	c.Assert(mail.GetReferenceCount(), Equals, 1)
}

func (s *MySuite) TestGetReferenceWithMultipleCc(c *C) {

	var mail Mail

	mail.Init()
	mail.Cc = append(mail.Cc, bson.NewObjectId())
	mail.Cc = append(mail.Cc, bson.NewObjectId())

	c.Assert(len(mail.Cc), Equals, 2)
	c.Assert(mail.GetReferenceCount(), Equals, 2)
}

func (s *MySuite) TestGetReferenceWithSingleBcc(c *C) {

	var mail Mail

	mail.Init()
	mail.Bcc = append(mail.Bcc, bson.NewObjectId())

	c.Assert(len(mail.Bcc), Equals, 1)
	c.Assert(mail.GetReferenceCount(), Equals, 1)
}

func (s *MySuite) TestGetReferenceWithMultipleBcc(c *C) {

	var mail Mail

	mail.Init()
	mail.Bcc = append(mail.Bcc, bson.NewObjectId())
	mail.Bcc = append(mail.Bcc, bson.NewObjectId())

	c.Assert(len(mail.Bcc), Equals, 2)
	c.Assert(mail.GetReferenceCount(), Equals, 2)
}

func (s *MySuite) TestGetReference(c *C) {

	var mail Mail

	mail.Init()
	mail.From = bson.NewObjectId()
	mail.To = append(mail.To, bson.NewObjectId())
	mail.To = append(mail.To, bson.NewObjectId())
	mail.Cc = append(mail.Cc, bson.NewObjectId())
	mail.Cc = append(mail.Cc, bson.NewObjectId())
	mail.Bcc = append(mail.Bcc, bson.NewObjectId())

	c.Assert(mail.From, NotNil)
	c.Assert(len(mail.To), Equals, 2)
	c.Assert(len(mail.Cc), Equals, 2)
	c.Assert(len(mail.Bcc), Equals, 1)
	c.Assert(mail.GetReferenceCount(), Equals, 6)

}
