package txsub

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultSubmissionList(t *testing.T) {
	ctx := context.Background()

	Convey("submissionList (The default openSubmissionList implementation)", t, func() {
		list := NewDefaultSubmissionList(time.Hour)
		realList := list.(*SubmissionList)
		hashes := []EnvelopeInfo{
			{ContentHash: "0000000000000000000000000000000000000000000000000000000000000000"},
			{ContentHash: "0000000000000000000000000000000000000000000000000000000000000001"},
		}

		listeners := []chan fullResult{
			make(chan fullResult, 1),
			make(chan fullResult, 1),
		}

		Convey("Add()", func() {
			Convey("adds an entry to the submission list when a new hash is used", func() {
				So(list.Add(ctx, &hashes[0], listeners[0]), ShouldBeNil)
				sub := realList.submissions[hashes[0].ContentHash]
				So(sub.Envelope.ContentHash, ShouldEqual, hashes[0].ContentHash)
				So(sub.SubmittedAt, ShouldHappenWithin, 1*time.Second, time.Now())

				// drop the send side of the channel by casting to listener
				var l listener = listeners[0]
				So(sub.Listeners[0], ShouldEqual, l)
			})

			Convey("adds an listener to an existing entry when a hash is used with a new listener", func() {
				So(list.Add(ctx, &hashes[0], listeners[0]), ShouldBeNil)
				sub := realList.submissions[hashes[0].ContentHash]
				st := sub.SubmittedAt
				<-time.After(20 * time.Millisecond)
				So(list.Add(ctx, &hashes[0], listeners[1]), ShouldBeNil)

				// increases the size of the listener
				So(len(sub.Listeners), ShouldEqual, 2)
				// doesn't update the submitted at time
				So(st == sub.SubmittedAt, ShouldEqual, true)
			})

			Convey("panics when the listener is not buffered", func() {
				So(func() { So(list.Add(ctx, &hashes[0], make(listener)), ShouldBeNil) }, ShouldPanic)
			})

			Convey("errors when the provided hash is not 64-bytes", func() {
				wrongEnvelope := EnvelopeInfo{
					ContentHash: "123",
				}
				err := list.Add(ctx, &wrongEnvelope, listeners[0])
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Finish()", func() {
			So(list.Add(ctx, &hashes[0], listeners[0]), ShouldBeNil)
			So(list.Add(ctx, &hashes[0], listeners[1]), ShouldBeNil)
			r := fullResult{
				Result: Result{
					Hash: hashes[0].ContentHash,
				},
			}
			So(list.Finish(ctx, r), ShouldBeNil)

			Convey("writes to every listener", func() {
				r1, ok1 := <-listeners[0]
				So(r1, ShouldResemble, r)
				So(ok1, ShouldBeTrue)

				r2, ok2 := <-listeners[1]
				So(r2, ShouldResemble, r)
				So(ok2, ShouldBeTrue)
			})

			Convey("removes the entry", func() {
				_, ok := realList.submissions[hashes[0].ContentHash]
				So(ok, ShouldBeFalse)
			})

			Convey("closes every listener", func() {
				<-listeners[0]
				_, more := <-listeners[0]
				So(more, ShouldBeFalse)

				<-listeners[1]
				_, more = <-listeners[1]
				So(more, ShouldBeFalse)
			})

			Convey("works when the noone is waiting for the result", func() {
				err := list.Finish(ctx, r)
				So(err, ShouldBeNil)
			})

		})

		Convey("Clean()", func() {
			So(list.Add(ctx, &hashes[0], listeners[0]), ShouldBeNil)
			<-time.After(200 * time.Millisecond)
			So(list.Add(ctx, &hashes[1], listeners[1]), ShouldBeNil)
			left := list.Clean(ctx, 200*time.Millisecond)

			So(left, ShouldEqual, 1)

			Convey("removes submissions older than the maxAge provided", func() {
				_, ok := realList.submissions[hashes[0].ContentHash]
				So(ok, ShouldBeFalse)
			})

			Convey("leaves submissions that are younger than the maxAge provided", func() {
				_, ok := realList.submissions[hashes[1].ContentHash]
				So(ok, ShouldBeTrue)
			})

			Convey("closes any cleaned listeners", func() {
				So(len(listeners[0]), ShouldEqual, 1)
				<-listeners[0]
				select {
				case _, stillOpen := <-listeners[0]:
					So(stillOpen, ShouldBeFalse)
				default:
					panic("cleaned listener is still open")
				}
			})
		})

		Convey("List() works as expected", func() {
			So(len(list.Pending()), ShouldEqual, 0)
			So(list.Add(ctx, &hashes[0], listeners[0]), ShouldBeNil)
			So(len(list.Pending()), ShouldEqual, 1)
			So(list.Add(ctx, &hashes[1], listeners[1]), ShouldBeNil)
			So(len(list.Pending()), ShouldEqual, 2)
		})
	})
}
