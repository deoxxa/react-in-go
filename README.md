React in go
===========

This is an experiment in rendering React JavaScript applications on the server,
using go for the backend. It builds upon the [otto](https://github.com/robertkrimen/otto)
JavaScript interpreter and my [ottoext](http://fknsrs.biz/p/ottoext) package to
provide just enough of a JavaScript environment to support rendering a React
application, provided that application uses very few browser APIs during its
critical render path.

NOTE: Right now, this project relies on an unmerged patch to otto's interpreter,
which fixes a bug involved in cloning VM instances. Please build the server
component with `godep go build` until this is merged, or apply otto PR #134
(commit `a9a9ffc43b4176b63426c44ba8447c50f83f4e85`) to your local copy of otto.
