
### Contribution Guidelines

Loki, the primary author of this code has somewhat unconventional ideas about everything including Go, developed out of the process of building this application and the use of compact logging syntax with visual recognisability as a key debugging technique for mostly multi-threaded code, and so, tools have been developed and conventions are found throughout the code and increasingly as all the dusty parts get looked at again.

So, there's a few things that you should do, which Loki thinks you will anyway realise that should be the default anyway. He will stop speaking in the third person now.

1. e is error

being a verb, `err` is a stolen name. But it also is 3 times as long. It is nearly universal among programmers that i, usually also j and k are iterator variables, and commonly also x, y, z, most cogently relevant when they are coordinates. So use e, not err

2. Use the logger, E.Chk for serious conditions or while debugging, non-blocking failures use D.Chk 

It's not difficult, literally just copy any log.go file out of this repo, change the package name and then you can use the logger, which includes these handy check functions that some certain Go Authors use in several of their lectures and texts, and actually probably the origin of the whole idea of making the error type a pointer to string, which led to the Wrap/Unwrap interface which I don't use. 

4. Use if statements with ALL calls that return errors or booleans

And declare their variables with var statements, unless there is no more uses of the `e` again before all paths of execution return.

5. Prefer to return early

The `.Chk` functions return true if not nil, and so without the `!` in front, the most common case, the content of the following `{}` are the error handling code. 

In general, and especially when the process is not idempotent (changed order breaks the process), which will be most of the time with several processes in a sequence, especially in Gio code which is naturally a bit wide if you write it to make it easily sliced, you want to keep the success line of execution as far left as possible.

Sometimes the negative condition is ignored, as there is a retry or it is not critical, and for these cases use `!E.Chk` and put the success path inside its if block.

6. In the Gio code of the wallet, take advantage of the fact that you can break a line after the dot `.` operator and as you will see amply throughout the code, as it allows items to be repositioned and added with minimal fuss.

Deep levels of closures and stacked variables of any kind tend to lead quickly to a lot of nasty red lines in one's IDE. The fluent method chaining pattern is used because it is far more concise than the raw type definition followed by closure attached by a dot operator, but since it would be valid and pass `gofmt` unchanged to put the whole chain in there (assuming it somehow had no anonymous functions in it).

7. Use the logger

This is being partly repeated as it's very important. Regardless of programmer opinions about whether a debugger is a better tool than a log viewer, note that while it is not fully implemented, `pod` already contains a protocol to aggregate child process logs invisibly from the terminal through an IPC, and logging is one means to enabling auditability of code. 

So long as logs concern themselves primarily with metadata information and only expose data in `trace` level (with the `T` error type) and put the really heavy stuff like printing tree walks over thousands of nodes or other similarly very short operations, put them inside closures with `T.C(func()string{})`

8. `gofmt` sort the imports, and avoid whenever possible a package name different from a folder, unless you can't avoid an import name conflict. A notable example is uber/atomic. 

9. 120 characters wide

We are living in the 21st century. It was already common as this century started that monitors were big enough to show 120 or more characters wide, 80 characters is just too narrow. Well, it would be even better if I figured out a way to make the Gio code not stack out to the right so deep but it's comfortable in 120 unless I have put too many things into one closure.

10. Always put doc comments on exported functions. Put them on struct fields and variables also in the case that detailed information isn't already available elsewhere related to the item. Keep the doc comments to one line unless you really need to explain something that needs to be visible in the documentation. For the most part this means libraries and for the most part a lot of that is in independent repositories.
