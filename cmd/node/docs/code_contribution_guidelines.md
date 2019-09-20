### Table of Contents

1. [Overview](#Overview)<br />
2. [Minimum Recommended Skillset](#MinSkillset)<br />
3. [Required Reading](#ReqReading)<br />
4. [Development Practices](#DevelopmentPractices)<br />
   4.1. [Share Early, Share Often](#ShareEarly)<br />
   4.2. [Testing](#Testing)<br />
   4.3. [Code Documentation and Commenting](#CodeDocumentation)<br />
   4.4. [Model Git Commit Messages](#ModelGitCommitMessages)<br />
5. [Code Approval Process](#CodeApproval)<br />
   5.1 [Code Review](#CodeReview)<br />
   5.2 [Rework Code (if needed)](#CodeRework)<br />
   5.3 [Acceptance](#CodeAcceptance)<br />
6. [Contribution Standards](#Standards)<br />
   6.1. [Contribution Checklist](#Checklist)<br />
   6.2. [Licensing of Contributions](#Licensing)<br />

<a name="Overview"></a>

### 1. Overview

Developing cryptocurrencies is an exciting endeavor that touches a wide variety of areas such as wire protocols, peer-to-peer networking, databases, cryptography, language interpretation (transaction scripts), RPC, and websockets. They also represent a radical shift to the current fiscal system and as a result provide an opportunity to help reshape the entire financial system. There are few projects that offer this level of diversity and impact all in one code base.

However, as exciting as it is, one must keep in mind that cryptocurrencies represent real money and introducing bugs and security vulnerabilities can have far more dire consequences than in typical projects where having a small bug is minimal by comparison. In the world of cryptocurrencies, even the smallest bug in the wrong area can cost people a significant amount of money. For this reason, the pod suite has a formalized and rigorous development process which is outlined on this page.

We highly encourage code contributions, however it is imperative that you adhere to the guidelines established on this page.

<a name="MinSkillset"></a>

### 2. Minimum Recommended Skillset

The following list is a set of core competencies that we recommend you possess before you really start attempting to contribute code to the project. These are not hard requirements as we will gladly accept code contributions as long as they follow the guidelines set forth on this page. That said, if you don't have the following basic qualifications you will likely find it quite difficult to contribute.

- A reasonable understanding of bitcoin at a high level (see the [Required Reading](#ReqReading) section for the original white paper)

- Experience in some type of C-like language. Go is preferable of course.

- An understanding of data structures and their performance implications

- Familiarity with unit testing

- Debugging experience

- Ability to understand not only the area you are making a change in, but also the code your change relies on, and the code which relies on your changed code

Building on top of those core competencies, the recommended skill set largely depends on the specific areas you are looking to contribute to. For example, if you wish to contribute to the cryptography code, you should have a good understanding of the various aspects involved with cryptography such as the security and performance implications.

<a name="ReqReading"></a>

### 3. Required Reading

- [Effective Go](http://golang.org/doc/effective_go.html) - The entire pod suite follows the guidelines in this document. For your code to be accepted, it must follow the guidelines therein.

- [Original Satoshi Whitepaper](http://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&ved=0CCkQFjAA&url=http%3A%2F%2Fbitcoin.org%2Fbitcoin.pdf&ei=os3VUuH8G4SlsASV74GoAg&usg=AFQjCNEipPLigou_1MfB7DQjXCNdlylrBg&sig2=FaHDuT5z36GMWDEnybDJLg&bvm=bv.59378465,d.b2I) - This is the white paper that started it all. Having a solid foundation to build on will make the code much more comprehensible.

<a name="DevelopmentPractices"></a>

### 4. Development Practices

Developers are expected to work in their own trees and submit pull requests when they feel their feature or bug fix is ready for integration into the master branch.

<a name="ShareEarly"></a>

### 4.1 Share Early, Share Often

We firmly believe in the share early, share often approach. The basic premise of the approach is to announce your plans **before** you start work, and once you have started working, craft your changes into a stream of small and easily reviewable commits.

This approach has several benefits:

- Announcing your plans to work on a feature **before** you begin work avoids duplicate work

- It permits discussions which can help you achieve your goals in a way that is consistent with the existing architecture

- It minimizes the chances of you spending time and energy on a change that might not fit with the consensus of the community or existing architecture and potentially be rejected as a result

- Incremental development helps ensure you are on the right track with regards to the rest of the community

- The quicker your changes are merged to master, the less time you will need to spend rebasing and otherwise trying to keep up with the main code base

<a name="Testing"></a>

### 4.2 Testing

One of the major design goals of all core pod packages is to aim for complete test coverage. This is financial software so bugs and regressions can cost people real money. For this reason every effort must be taken to ensure the code is as accurate and bug-free as possible. Thorough testing is a good way to help achieve that goal.

Unless a new feature you submit is completely trivial, it will probably be rejected unless it is also accompanied by adequate test coverage for both positive and negative conditions. That is to say, the tests must ensure your code works correctly when it is fed correct data as well as incorrect data (error paths).

Go provides an excellent test framework that makes writing test code and checking coverage statistics straight forward. For more information about the test coverage tools, see the [golang cover blog post](http://blog.golang.org/cover).

A quick summary of test practices follows:

- All new code should be accompanied by tests that ensure the code behaves correctly when given expected values, and, perhaps even more importantly, that it handles errors gracefully

- When you fix a bug, it should be accompanied by tests which exercise the bug to both prove it has been resolved and to prevent future regressions

<a name="CodeDocumentation"></a>

### 4.3 Code Documentation and Commenting

Comments have a way of turning into lies during development, you should not expect readers to depend on it. Much more important is that names are meaningful, they do not take up excessive space, and comments are only necessary when the meaning of the code needs clarification.

- At a minimum every function must be commented with its intended purpose and any assumptions that it makes

  - Function comments must always begin with the name of the function per [Effective Go](http://golang.org/doc/effective_go.html)

  - Function comments should be complete sentences since they allow a wide variety of automated presentations such as [godoc.org](https://godoc.org)

  - Comments should be brief, function type signatures should be informative enough that the comment is for clarification. Comments are not tested by the compiler, and can obscure the intent of the code if the code is opaque in its semantics.

  - Comments will be parsed by godoc and excess vertical space usage reduces the readability of code, so there is no sane reason why the comments (and indeed, in documents such as this) should be manually split into lines. That's what word wrap is for.

  - The general rule of thumb is to look at it as if you were completely unfamiliar with the code and ask yourself, would this give me enough information to understand what this function does and how I'd probably want to use it?

  - Detailed information in comments should be mainly in the type definitions. Meaningful names in function parameters are more important than silly long complicated comments and make the code harder to read where the clarity is most needed.

  - If you need to write a lot of comments about code you probably have not written it well.

  - Variable and constant names should be informative, and where obvious, brief.

  - If the function signature is longer than 80 characters in total, you should change the parameters to be a structured variable, the structure will explain the parameters better and more visually attractive than a function call with more than 5 parameters.

  - If you use a constant value more than a few times, and especially, in more than a few source files, you should give it a meaningful name and place it into separate folders that allows you to avoid circular dependencies. It is better to define a type independently, and create an alias in the implementation as methods must be defined locally to the type. However, if you need to access the fields of the type, the definition needs to be isolated separate from the implementation, otherwise you almost certainly will run into a circular dependency that will block compilation.

  - The best place for detailed information is a separate `doc.go` file, where the comment that appears before the package name appears at the very top in the Godoc output, and in the structure and type definitions for exported types. Functions should not be the place to put this, as it interferes with readability, and scatters the information, at the same time.

<a name="ModelGitCommitMessages"></a>

### 4.4 Model Git Commit Messages

This project prefers to keep a clean commit history with well-formed commit messages. This section illustrates a model commit message and provides a bit of background for it. This content was originally created by Tim Pope and made available on his website, however that website is no longer active, so it is
being provided here.

Here’s a model Git commit message:

```
Short (50 chars or less) summary of changes

More detailed explanatory text, if necessary.  Wrap it to about 72
characters or so.  In some contexts, the first line is treated as the
subject of an email and the rest of the text as the body.  The blank
line separating the summary from the body is critical (unless you omit
the body entirely); tools like rebase can get confused if you run the
two together.

Write your commit message in the present tense: "Fix bug" and not "Fixed
bug."  This convention matches up with commit messages generated by
commands like git merge and git revert.

Further paragraphs come after blank lines.

- Bullet points are okay, too
- Typically a hyphen or asterisk is used for the bullet, preceded by a
  single space, with blank lines in between, but conventions vary here
- Use a hanging indent
```

Prefix the summary with the subsystem/package when possible. Many other projects make use of the code and this makes it easier for them to tell when something they're using has changed. Have a look at [past commits](https://git.parallelcoin.io/dev/pod/commits/master) for examples of commit messages.

<a name="CodeApproval"></a>

### 5. Code Approval Process

This section describes the code approval process that is used for code contributions. This is how to get your changes into pod.

<a name="CodeReview"></a>

### 5.1 Code Review

All code which is submitted will need to be reviewed before inclusion into the master branch. This process is performed by the project maintainers and usually other committers who are interested in the area you are working in as well.

##### Code Review Timeframe

The timeframe for a code review will vary greatly depending on factors such as the number of other pull requests which need to be reviewed, the size and complexity of the contribution, how well you followed the guidelines presented on this page, and how easy it is for the reviewers to digest your commits. For example, if you make one monolithic commit that makes sweeping changes to things in multiple subsystems, it will obviously take much longer to review. You will also likely be asked to split the commit into several smaller, and hence more manageable, commits.

Keeping the above in mind, most small changes will be reviewed within a few days, while large or far reaching changes may take weeks. This is a good reason to stick with the [Share Early, Share Often](#ShareOften) development practice outlined above.

##### What is the review looking for?

The review is mainly ensuring the code follows the [Development Practices](#DevelopmentPractices) and [Code Contribution Standards](#Standards). However, there are a few other checks which are generally performed as follows:

- The code is stable and has no stability or security concerns
- The code is properly using existing APIs and generally fits well into the overall architecture
- The change is not something which is deemed inappropriate by community consensus

<a name="CodeRework"></a>

### 5.2 Rework Code (if needed)

After the code review, the change will be accepted immediately if no issues are found. If there are any concerns or questions, you will be provided with feedback along with the next steps needed to get your contribution merged with master. In certain cases the code reviewer(s) or interested committers may help you rework the code, but generally you will simply be given feedback for you to make the necessary changes.

This process will continue until the code is finally accepted.

<a name="CodeAcceptance"></a>

### 5.3 Acceptance

Once your code is accepted, it will be integrated with the master branch. Typically it will be rebased and fast-forward merged to master as we prefer to keep a clean commit history over a tangled weave of merge commits. However,regardless of the specific merge method used, the code will be integrated with the master branch and the pull request will be closed.

Rejoice as you will now be listed as a [contributor](https://git.parallelcoin.io/dev/pod/graphs/contributors)!

<a name="Standards"></a>

### 6. Contribution Standards

<a name="Checklist"></a>

### 6.1. Contribution Checklist

- [&nbsp;&nbsp;] All changes are Go version 1.11 compliant

- [&nbsp;&nbsp;] The code being submitted is commented according to the [Code Documentation and Commenting](#CodeDocumentation) section

- [&nbsp;&nbsp;] For new code: Code is accompanied by tests which exercise both the positive and negative (error paths) conditions (if applicable)

- [&nbsp;&nbsp;] For bug fixes: Code is accompanied by new tests which trigger the bug being fixed to prevent regressions

- [&nbsp;&nbsp;] Any new logging statements use an appropriate subsystem and logging level

- [&nbsp;&nbsp;] Code has been formatted with `go fmt`

- [&nbsp;&nbsp;] Running `go test` does not fail any tests

- [&nbsp;&nbsp;] Running `go vet` does not report any issues

- [&nbsp;&nbsp;] Running [golint](https://github.com/golang/lint) does not report any **new** issues that did not already exist

<a name="Licensing"></a>

### 6.2. Licensing of Contributions

All contributions must be licensed with the [ISC license](https://git.parallelcoin.io/dev/pod/blob/master/LICENSE). This is the same license as all of the code in the pod suite.
