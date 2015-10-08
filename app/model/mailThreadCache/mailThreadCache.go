package mailThreadCache

import (
	"time"

	"github.com/PradheepShrinivasan/mailbox/app/model/mail"

	"gopkg.in/mgo.v2/bson"
)

type MailCache struct {
	MailId    bson.ObjectId `bson:"_id"`
	Read      bool          `bson:"read"`
	From      bson.ObjectId `bson:"from"`
	TimeStamp time.Time     `bson:"time"`
	Subject   string        `bson:"subject"`
}

func (cache *MailCache) Init(m *mail.Mail) {

	cache.MailId = m.Id
	cache.From = m.From
	cache.Subject = m.Subject
	cache.TimeStamp = m.TimeStamp
}

type MailThreadCache struct {
	ThreadId   bson.ObjectId `bson:"_id"` // this is the primary key
	MailCaches []MailCache   `bson:"mailCache"`
}

func (mTCache *MailThreadCache) AddMailToMailThreadCache(mTId bson.ObjectId, m *mail.Mail) {

	var mCache MailCache

	mCache.Init(m)
	mTCache.ThreadId = mTId
	mTCache.MailCaches = append(mTCache.MailCaches, mCache)
}

func CreateMailThreadCache(threadId bson.ObjectId, m *mail.Mail) MailThreadCache {

	var mTCache MailThreadCache
	mTCache.AddMailToMailThreadCache(threadId, m)
	return mTCache
}

func CreateMailCache(m *mail.Mail) MailCache {

	var mCache MailCache
	mCache.Init(m)
	return mCache
}
