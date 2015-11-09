React in go
===========

This is an experiment in rendering React JavaScript applications on the server,
using go for the backend. It builds upon the [otto](https://github.com/robertkrimen/otto)
JavaScript interpreter and my [ottoext](http://fknsrs.biz/p/ottoext) package to
provide just enough of a JavaScript environment to support rendering a React
application, provided that application uses very few browser APIs during its
critical render path.

NOTE: Right now, this project relies on several unmerged patches to otto's
interpreter, which fix bugs involved in cloning VM instances and parsing certain
expressions, and add support for source maps. Please build the server
component with `godep go build` until this is merged, or apply the following
otto PRs to your local copy of otto.

* #115 (`432b4361a78dd336691458dea33bda22ab158e7a`) _Fixed TypeError message
	when evaluating a.b["c"] where a.b is undefined_
* #134 (`11907dfb901ea7521123bb6dc1ac7a31cc11115d`) _add test for for cloning
	an object with a get/set property_
* #138 (`18c054f15ea76a2124dfbad05e0059e573241be1`) _add sourcemap support_
* #139 (`33ac8bd28a37753ea1e8649cadd4196be6e55e06`) _fix parsing of statements
	like `new a["b"]`_
