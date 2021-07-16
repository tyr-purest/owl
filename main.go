package main

import (
	"github.com/tyr-purest/owl/broker"
	"github.com/tyr-purest/owl/transport"

	"github.com/davecgh/go-spew/spew"
	"github.com/kataras/iris/v12"
)

func main() {
	broker.Init()
	transport.NewTransport()
	app := iris.Default()

	app.Handle("GET", "/ws", func(c iris.Context) {
		session := "sessionID"
		node, _ := transport.NewNode(c.ResponseWriter(), c.Request())
		transport.Join(node, session)

		node.OnNotify = func(e *transport.NotifyEvent) {
			spew.Dump("OnNotify", e)
		}

		node.Notify(transport.Notify, &transport.NotifyEvent{
			Send:    node.ID,
			Session: session,
		})

		<-node.Ws.Close

		transport.Leave(session, node.ID)
	})

	app.Run(iris.Addr(":3000"), iris.WithoutInterruptHandler)
}
