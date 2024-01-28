Go+ Classfiles
=====

Rob Pike once said that if he could only introduce one feature to Go, he would choose `interface` instead of `goroutine`. `classfile` is as important to Go+ as `interface` is to Go.

In the design philosophy of Go+, we do not recommend `DSL` (Domain Specific Language). But `SDF` (Specific Domain Friendliness) is very important. The Go+ philosophy about `SDF` is:

```
Don't define a language for specific domain.
Abstract domain knowledge for it.
```

Go+ introduces `classfile` to abstract domain knowledge.


### classfile: Unit Test

Go+ has built-in support for a classfile to simplify unit testing. This classfile has the file suffix `_test.gox`.

Suppose you have a function named `foo`:

```go
func foo(v int) int {
	return v * 2
}
```

Then you can create a `foo_test.gox` file to test it (see [unit-test/foo_test.gox](testdata/unit-test/foo_test.gox)):

```go
if v := foo(50); v != 100 {
	t.error "foo(50) ret: ${v}"
} else {
	t.log "test foo(50) ok"
}

t.run "foo -10", t => {
	if foo(-10) != -20 {
		t.fatal "foo(-10) != -20"
	}
}

t.run "foo 0", t => {
	if foo(0) != 0 {
		t.fatal "foo(0) != 0"
	}
}
```

You don't need to define a series of `TestXXX` functions like Go, just write your test code directly.

If you want to run a subtest case, use `t.run`.


### yap: Yet Another Go/Go+ HTTP Web Framework

This classfile has the file suffix `_yap.gox`.

Before using `yap`, you need to add it to `go.mod` by using `go get`:

```sh
go get github.com/goplus/yap@latest
```

Find `require github.com/goplus/yap` statement in `go.mod` and add `//gop:class` at the end of the line:

```go.mod
require github.com/goplus/yap v0.7.2 //gop:class
```

#### Router and Parameters

demo in Go+ classfile ([hello_yap.gox](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_hello/hello_yap.gox)):

```coffee
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}
handle "/", ctx => {
	ctx.html `<html><body>Hello, <a href="/p/123">Yap</a>!</body></html>`
}

run ":8080"
```

#### Static files

Static files server demo in Go+ classfile ([staticfile_yap.gox](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_static/staticfile_yap.gox)):

```coffee
static "/foo", FS("public")
static "/"

run ":8080"
```

#### YAP Template

demo in Go+ classfile ([blog_yap.gox](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_blog/blog_yap.gox), [article_yap.html](https://github.com/goplus/yap/blob/v0.7.2/demo/classfile_blog/yap/article_yap.html)):

```coffee
get "/p/:id", ctx => {
	ctx.yap "article", {
		"id": ctx.param("id"),
	}
}

run ":8080"
```

### yaptest: HTTP Test Framework

This classfile has the file suffix `_ytest.gox`.

Before using `yaptest`, you need to add `yap` to `go.mod`:

```go.mod
require github.com/goplus/yap v0.7.2 //gop:class
```

Then you can create a `example_ytest.gox` file to test your HTTP server:

```coffee
host "https://example.com", "http://localhost:8080"
testauth := oauth2("...")

run "urlWithVar", => {
	id := "123"
	get "https://example.com/p/${id}"
	ret
	echo "code:", resp.code
	echo "body:", resp.body
}

run "matchWithVar", => {
	code := Var(int)
	id := "123"
	get "https://example.com/p/${id}"
	ret code
	echo "code:", code
	match code, 200
}

run "postWithAuth", => {
	id := "123"
	title := "title"
	author := "author"
	post "https://example.com/p/${id}"
	auth testauth
	json {
		"title":  title,
		"author": author,
	}
	ret 200 # match resp.code, 200
	echo "body:", resp.body
}

run "matchJsonObject", => {
	title := Var(string)
	author := Var(string)
	id := "123"
	get "https://example.com/p/${id}"
	ret 200
	json {
		"title":  title,
		"author": author,
	}
	echo "title:", title
	echo "author:", author
}
```

### spx: A Go+ 2D Game Engine for STEM education

This classfile has the file suffix `.spx`. It is the earliest classfile in the world.

Before using `spx`, you need to add it to `go.mod` by using `go get`:

```sh
go get github.com/goplus/spx@latest
```

It's also a built-in classfile so you don't need append `//gop:class` after `require github.com/goplus/spx`.

#### tutorial/01-Weather

![Screen Shot1](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/1.jpg) ![Screen Shot2](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/2.jpg)

Through this example you can learn how to listen events and do somethings.

Here are some codes in [Kai.spx](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/Kai.spx):

```coffee
onStart => {
	say "Where do you come from?", 2
	broadcast "1"
}

onMsg "2", => {
	say "What's the climate like in your country?", 3
	broadcast "3"
}

onMsg "4", => {
	say "Which seasons do you like best?", 3
	broadcast "5"
}
```

We call `onStart` and `onMsg` to listen events. `onStart` is called when the program is started. And `onMsg` is called when someone calls `broadcast` to broadcast a message.

When the program starts, Kai says `Where do you come from?`, and then broadcasts the message `1`. Who will recieve this message? Let's see codes in [Jaime.spx](https://github.com/goplus/spx/blob/main/tutorial/01-Weather/Jaime.spx):

```coffee
onMsg "1", => {
	say "I come from England.", 2
	broadcast "2"
}

onMsg "3", => {
	say "It's mild, but it's not always pleasant.", 4
	# ...
	broadcast "4"
}
```

Yes, Jaime recieves the message `1` and says `I come from England.`. Then he broadcasts the message `2`. Kai recieves it and says `What's the climate like in your country?`.

The following procedures are very similar. In this way you can implement dialogues between multiple actors.

#### tutorial/02-Dragon

![Screen Shot1](https://github.com/goplus/spx/blob/main/tutorial/02-Dragon/1.jpg)

Through this example you can learn how to define variables and show them on the stage.

Here are all the codes of [Dragon](https://github.com/goplus/spx/blob/main/tutorial/02-Dragon/Dragon.spx):

```coffee
var (
	score int
)

onStart => {
	score = 0
	for {
		turn rand(-30, 30)
		step 5
		if touching("Shark") {
			score++
			play chomp, true
			step -100
		}
	}
}
```

We define a variable named `score` for `Dragon`. After the program starts, it moves randomly. And every time it touches `Shark`, it gains one score.

How to show the `score` on the stage? You don't need write code, just add a `stageMonitor` object into [assets/index.json](https://github.com/goplus/spx/blob/main/tutorial/02-Dragon/assets/index.json):

```json
{
  "zorder": [
    {
      "type": "stageMonitor",
      "target": "Dragon",
      "val": "getVar:score",
      "color": 15629590,
      "label": "score",
      "mode": 1,
      "x": 5,
      "y": 5,
      "visible": true
    }
  ]
}
```

#### tutorial/03-Clone

![Screen Shot1](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/1.png)

Through this example you can learn:
* Clone sprites and destory them.
* Distinguish between sprite variables and shared variables that can access by all sprites.

Here are some codes in [Calf.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/Calf.spx):

```coffee
var (
	id int
)

onClick => {
	clone
}

onCloned => {
	gid++
	...
}
```

When we click the sprite `Calf`, it receives an `onClick` event. Then it calls `clone` to clone itself. And after cloning, the new `Calf` sprite will receive an `onCloned` event.

In `onCloned` event, the new `Calf` sprite uses a variable named `gid`. It doesn't define in [Calf.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/Calf.spx), but in [main.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/main.spx).


Here are all the codes of [main.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/main.spx):

```coffee
var (
	Arrow Arrow
	Calf  Calf
	gid   int
)

run "res", {Title: "Clone and Destory (by Go+)"}
```

All these three variables in [main.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/main.spx) are shared by all sprites. `Arrow` and `Calf` are sprites that exist in this project. `gid` means `global id`. It is used to allocate id for all cloned `Calf` sprites.

Let's back to [Calf.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/Calf.spx) to see the full codes of `onCloned`:

```coffee
onCloned => {
	gid++
	id = gid
	step 50
	say id, 0.5
}
```

It increases `gid` value and assigns it to sprite `id`. This makes all these `Calf` sprites have different `id`. Then the cloned `Calf` moves forward 50 steps and says `id` of itself.

Why these `Calf` sprites need different `id`? Because we want destory one of them by its `id`.

Here are all the codes in [Arrow.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/Arrow.spx):

```coffee
onClick => {
	broadcast "undo", true
	gid--
}
```

When we click `Arrow`, it broadcasts an "undo" message (NOTE: We pass the second parameter `true` to broadcast to indicate we wait all sprites to finish processing this message).

All `Calf` sprites receive this message, but only the last cloned sprite finds its `id` is equal to `gid` then destroys itself. Here are the related codes in [Calf.spx](https://github.com/goplus/spx/blob/main/tutorial/03-Clone/Calf.spx):

```coffee
onMsg "undo", => {
	if id == gid {
		destroy
	}
}
```

#### tutorial/04-Bullet

![Screen Shot1](https://github.com/goplus/spx/blob/main/tutorial/04-Bullet/1.jpg)

Through this example you can learn:
* How to keep a sprite following mouse position.
* How to fire bullets.

It's simple to keep a sprite following mouse position. Here are some related codes in [MyAircraft.spx](https://github.com/goplus/spx/blob/main/tutorial/04-Bullet/MyAircraft.spx):


```coffee
onStart => {
	for {
		# ...
		setXYpos mouseX, mouseY
	}
}
```

Yes, we just need to call `setXYpos mouseX, mouseY` to follow mouse position.

But how to fire bullets? Let's see all codes of [MyAircraft.spx](https://github.com/goplus/spx/blob/main/tutorial/04-Bullet/MyAircraft.spx):

```coffee
onStart => {
	for {
		wait 0.1
		Bullet.clone
		setXYpos mouseX, mouseY
	}
}
```

In this example, `MyAircraft` fires bullets every 0.1 seconds. It just calls `Bullet.clone` to create a new bullet. All the rest things are the responsibility of `Bullet`.

Here are all the codes in [Bullet.spx](https://github.com/goplus/spx/blob/main/tutorial/04-Bullet/Bullet.spx):

```coffee
onCloned => {
	setXYpos MyAircraft.xpos, MyAircraft.ypos+5
	show
	for {
		wait 0.04
		changeYpos 10
		if touching(Edge) {
			destroy
		}
	}
}
```

When a `Bullet` is cloned, it calls `setXYpos MyAircraft.xpos, MyAircraft.ypos+5` to follow `MyAircraft`'s position and shows itself (the default state of a `Bullet` is hidden). Then the `Bullet` moves forward every 0.04 seconds and this is why we see the `Bullet` is flying.

At last, when the `Bullet` touches screen `Edge` or any enemy (in this example we don't have enemies), it destroys itself.

These are all things about firing bullets.
