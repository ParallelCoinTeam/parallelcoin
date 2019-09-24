# Cleaner

Cleaner simply just reads in a source file, arranges the the source in the following order, grouping all root level declarations:

- types
- const
- var
- func

It simply sorts them lexicographically inside these groups. 

Constants are left as groups because of iota increments, but variables and types are all split into one per declaration.

The function sorting sorts by *, receiver type and function name, in that order, and being that ASCII puts the capital letters first, the exported methods are higher up and unexported lower down, within each of the classes in the file.

A key purpose of this little tool is to assist with splitting and merging source files within packages so as to reduce the amount of scrolling required, and so logically connected things are grouped together rather than hodge podge. 

Variable, type and constant declarations are often tempting to mix between functions near where the data / type is used, but again, this can make it difficult to get a clear and quick gauge of the contents of a source file.

I personally have a rule if I have to scroll through more than two screenfuls the source is probably getting too long anyway (and probably is a lot copypasta), and as I attempted to clean up the source in this project repository, hoping I would finish before the second coming, this tool would have saved me so much time.

Long source files and huge APIs are a maintenance nightmare. Hopefully this will help you avoid that, or more quickly deal with a mess you inherited from some former C++ programmer.

A little known fact about Go's syntax is that receiver/parameter/return value blocks enclosed in brackets can be split with newlines if the final item inside them ends with a comma. At this point they are not fully split, only the receiver is split, parameter lines are broken at the start but none of the rest (I'm still a beginner with regexp). You might wonder why someone would want things formatted this way - in my opinion scanning down through text a list of items is more readable when its items are broken up per line. Hopefully in the future I will finish the function parameter section splitter but for now it is nice that at least parts of it are now automated.

## Known issues

- The sorting algorithm based on the ast-based sorter removes functions with anonymous receivers. This is addressed by preprocessing to add an anonymous name `_` for the receivers without an address. I am not sure why the AST parser deletes these anon receivers but if they are given this nominal non-address it works.

- Parenthesised var and type declarations are split (so they sort properly) but due to quirks in the `dst` library (a fork of go/ast that keeps whitespace and comments together with the code they refer to) it is not simple to break the group and properly format it. Manual removal of the braces will be required.