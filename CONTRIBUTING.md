# Contribution Guidelines

Very likely the following recommendations will be in disagreement with
many other Go developers, because they are sheep and write their code
in a way that is dictated to them rather than actually thinking about the
ergonomics and debugging considerations.

Here at Parallelcoin, we care about most especially bug free code and
simple to read and understand code. We are writing systems that handle
people's literal money and as such, conventions can go to hell if they
create risk for our users.

Any contributions that don't obey the following rules will be amended to
conform, or if they are extensive and large pieces of code, the changes
will be summarily be rejected. Follow these guidelines or go find another
project to work on. These guidelines came out of 18 months of intensive
work on this project and saved huge amounts of time already.

## Coding Style

### Text format

We think it is most antiquated and outdated to talk about 80 characters
as the standard width for a line of code. Even back when there was still
a lot of people working on dumb terminals, 120 characters quickly became
common because quite simply, 80 pixels is an artifact of old television
video standards, and while we discourage writing of excessively long
lines of code, 120 characters avoids a lot of linebreaks that disrupt
the flow while reading.

For this reason also, we don't want the bureaucratic excessively
verbose style commonly seen in Java and similar OOP languages. The simple
way to explain it is the smaller the scope, the smaller the name.
The exception to this is method receivers, which should also be short
due to being used a lot.

Also, when there is a list of items such as in a return value if the
signature would break two lines of length it probably needs to have
the values structured, but in case it's unavoidable then list them 
vertically for easier reading, for either or both parameters and return
values. 

### Named Returns

Many Go coders are much too fond of the `:=` operator and write code
that often creates several instances of the same `err` variable.

In the experience of the lead developer, one of the most pernicious and
difficult to find bugs is scope shadowing, where a scope has a variable
declared in an outer scope that is then re-declared in an inner scope.

By making all return variables named, redundant creation assignments can
be eliminated, for pointers and slices they need to be `make`d or 
`new`ed anyway but the `:=` is unnecessary.

Even if there is only one line of a return calling another function,
all return values should be given names in their signature to both
inform the consumer of the function what is returned and to 
reduce the chance of scope shadowing occurring.

Though it is unlikely to often happen, by avoiding remaking variables
that will be returned can potentially eliminate the need for heap
allocations and prevent the need for extending the stack.

Very occasionally the obverse applies, but more often ill-considered 
use of `:=` increases memory consumption and besides performance
considerations, the code is more readable when declarations are more 
explicit.

Some devs don't like naked returns because of long function bodies, but
in addition to encouraging the use of named and naked returns, functions
that are more than a screenful between signature, and the end of the 
function block, probably need to be refactored so that inner parts of 
the function are separately defined in a helper function, thus negating
this common justification for explicit, verbose return statements.

### Keep functions as simple as possible

Following on from the last paragraph of the foregoing, functions should
not generally be more than fits on an average SD display at standard 
font scale. The longer a function is, the greater chances of unintended
side effects appearing in the code, which is another type of bug 
practically universal outside of functional programming.

Side effects are sometimes required, but should be avoided whenever 
possible. Aiming to keep functions single-purpose and short helps
with debugging and eliminates bugs caused by unintended side effects.

### Use slog

The library used for logging in this project is 
https://github.com/stalker-loki/app/slog which provides a convenient
`Check` function which checks if the error is nil and if not prints
it to the logger and returns true if it is not nil. This function
prints at `debug` level, so in normal default configuration it is not
printed.

One of the most common constructs found in many functions is to attempt
some call and if it fails, to return. Many times coders explicitly
specify the values to return, but it is much tidier and easier to read
if instead the *named* return variables are assigned, and the
return is naked.

Concurrent programming, unlike pure procedural or imperative programming
is not well provided for by debuggers such as Go's Delve debugger, at 
least as far as standard interfaces to it go, are practically useless
for debugging highly concurrent code, and for this reason debugging
with Go really requires some kind of light-weight tracing which can
record the traces of multiple concurrent threads of execution.

```go
	case PubKeyTy:
		// A pay-to-pubkey script is of the form:  <pubkey> OP_CHECKSIG Therefore the pubkey is the first item on the 
		// stack. Skip the pubkey if it's invalid for some reason.
		requiredSigs = 1
		addr, err := util.NewAddressPubKey(pops[0].data, chainParams)
		if err == nil {
			addrs = append(addrs, addr)
		}
```

should be like this:

```go
	case PubKeyTy:
		// A pay-to-pubkey script is of the form:  <pubkey> OP_CHECKSIG Therefore the pubkey is the first item on the 
		// stack. Skip the pubkey if it's invalid for some reason.
		requiredSigs = 1
		var addr *util.AddressPubKey
		if addr, err = util.NewAddressPubKey(pops[0].data, chainParams); !slog.Check(err) {
			addrs = append(addrs, addr)
		}
```

This way it is always explicit the type of the variable, and the oft-used error variable there is one in every
scope, and only one. They are not used in select blocks anyway. As mentioned above, functions should not be
excessively long anyway for the context of return statements to be unclear.

The other thing is the logging, which here is then printed at the site of its appearance, where usually the error
in the code is found. Though several of the better Go coders use these in their lectures and tutorials somehow the
check/print log function has not become a part of the language, probably due to the general use of excessively verbose 
returns, and long functions, a hallmark of C++ programmers.

Where they are positioned matters too. There is no point in creating the variable long before it is going to be used
and with many conditions under which it may never have to be done. 

An exception to this might be in the creation of large complex APIs, but usually it is a return value and its type is
always explicitly shown, by using its conventional receiver variable name also makes its references more consistent.

Also, by making the declaration explicit, when it is repeated many times inside scopes of a common outer scope,
it can be easily seen that one variable can be created and used by several cases of a switch or strings of if/then/err
type blocks. In the case above the other cases use differing concrete types though they are all of an interface, thus
the declaration is in situ.

```go
	case MultiSigTy:
		// A multi-signature script is of the form:  <numsigs> <pubkey> <pubkey> <pubkey>... <numpubkeys> OP_CHECKMULTISIG Therefore the number of required signatures is the 1st item on the stack and the number of public keys is the 2nd to last item on the stack.
		requiredSigs = asSmallInt(pops[0].opcode)
		numPubKeys := asSmallInt(pops[len(pops)-2].opcode)
		// Extract the public keys while skipping any that are invalid.
		addrs = make([]util.Address, 0, numPubKeys)
		var addr *util.AddressPubKey
		for i := 0; i < numPubKeys; i++ {
			if addr, err = util.NewAddressPubKey(pops[i+1].data, chainParams); !slog.Check(err) {
				addrs = append(addrs, addr)
			}
		}
```

In the above you see the addr variable is created before the loop as it is shared by every iteration of the loop, which
is sequential and non-concurrent. In this way the variable is not incurring allocation overhead when it is being
written to before reading.

#### Log the error at the site

Furthermore, the slog logging library prints the code locations of log
entries. A total time waster for debugging using trace logging is when 
the log is printed several levels further up than the site of the error
which is time wasted tracing the error back to the source.

There is some efforts in the Go 2 plans to improve this tracability but
for the most part it doesn't address the fact that finding a bug several 
layers deeper than the site it is printed the programmer still has to 
find that location, and it becomes doubly complicated when the function
is an interface implementation, as by default editors will first take 
you back to the interface definition and *then* let you follow that back
to the actual executing code.

By printing the error *in* that code, this time is saved, meaning faster 
feedback cycles in attempts to fix the bug. 
