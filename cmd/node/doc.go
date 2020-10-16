/*Package node is a full-node Parallelcoin implementation written in Go.

The default options are sane for most users. This means pod will work 'out of the box' for most users. However, there
are also a wide variety of flags that can be used to control it.

The following section provides a usage overview which enumerates the flags. An interesting point to note is that the
long form of all of these options ( except -C/--configfile and -D --datadir) can be specified in a configuration file
that is automatically parsed when pod starts up. By default, the configuration file is located at ~/.pod/pod. conf on
POSIX-style operating systems and %LOCALAPPDATA%\pod\pod. conf on Windows. The -D (--datadir) flag, can be used to
override this location.

NAME:
   pod node - start parallelcoin full node

USAGE:
   pod node [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

COMMANDS:
     dropaddrindex  drop the address search index
     droptxindex    drop the address search index
     dropcfindex    drop the address search index

GLOBAL OPTIONS:
   --help, -h  show help
*/
package node
