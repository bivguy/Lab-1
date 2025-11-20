package scheduler

import (
	m "github.com/bivguy/Comp412/models"
)

var DEPENDENCE_GRAPHS = []*DependenceGraph{
	func() *DependenceGraph {
		// a=1, b=2, c=3, d=4, e=5, f=6, g=7, h=8, i=9
		a := makeNode(1, "load")  // loadAI rARP,@a => r1
		b := makeNode(2, "add")   // add r1,r1 => r2
		c := makeNode(3, "load")  // loadAI rARP,@b => r3
		d := makeNode(4, "mult")  // mult r2,r3 => r4
		e := makeNode(5, "load")  // loadAI rARP,@c => r5
		f := makeNode(6, "mult")  // mult r4,r5 => r6
		g := makeNode(7, "load")  // loadAI rARP,@d => r7
		h := makeNode(8, "mult")  // mult r6,r7 => r8
		i := makeNode(9, "store") // storeAI r8 => @a

		gph := makeGraph(a, b, c, d, e, f, g, h, i)

		// connect all the DATA edges
		gph.ConnectNodes(b, a, m.DATA) // b -> a
		gph.ConnectNodes(d, b, m.DATA) // d -> b
		gph.ConnectNodes(d, c, m.DATA) // d -> c
		gph.ConnectNodes(f, d, m.DATA) // f -> d
		gph.ConnectNodes(f, e, m.DATA) // f -> e
		gph.ConnectNodes(h, f, m.DATA) // h -> f
		gph.ConnectNodes(h, g, m.DATA) // h -> g
		gph.ConnectNodes(i, h, m.DATA) // i -> h

		// connect all the SERIALIZATION edges
		gph.ConnectNodes(i, a, m.SERIALIZATION)
		gph.ConnectNodes(i, c, m.SERIALIZATION)
		gph.ConnectNodes(i, e, m.SERIALIZATION)
		gph.ConnectNodes(i, g, m.SERIALIZATION)

		return gph
	}(),
}
