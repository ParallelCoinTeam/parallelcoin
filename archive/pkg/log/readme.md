# logg

This is a very simple, but practical library for 
logging in applications. Currenty in process of
superseding the logi library, once the pipe logger 
is converted, as it is integrated into the block
sync progress logger.

To use it, create a file in a package you want
to add the logger to with the name `log.go`, and then:

```bash
go run ./pkg/logg/deploy/.
```

from the root of the github.com/p9c/pod repository
and it replicates its template with the altered
package name to match the folder name - so you need
to adhere to the rule that package names and the folder
name are the same.

The library includes functions to toggle the filtering,
highlight and filtering sets while it is running.