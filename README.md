# runme - A simple markdown extractor

This is a simple project which allows you to list, extract, and execute commands from markdown files.

Given the following code-block you'll see three things:

```/bin/bash test
uptime
```

> **NOTE**: If you're viewing this file on github you'll need to select the [RAW View](https://raw.githubusercontent.com/skx/runme/master/README.md) to see the markup.

The block has:

* A shell defined.
   * `/bin/bash` in this case
* A name defined.
   * `test` in this case.
* A command, or series of commands to execute.
   * `uptime` in this case.

Another block might use python3:

```/usr/bin/python3 home
import os
print(os.environ['HOME'])
```

Finally to test running all blocks, as one script, we can have another shell block:

```/bin/bash whoami
id
```



## Listing Blocks

`runme` will process `README.md` in the current directory, if it exists, otherwise you will be expected to provide the name(s) of the file(s) to process.  In all the later examples I've specified that filename explicitly, just to avoid confusion.

By default all blocks found will be shown, for example:

```bash
Shell:/bin/bash  Name:test
uptime

Shell:/usr/bin/python3  Name:home
import os
print(os.environ['HOME'])

Shell:/bin/bash  Name:whoami
id

```

Notice here that only blocks with **both** a shell and a name are listed.

You can filter the output to only show blocks with a given shell, or name via the flags:

* `runme --name test README.md`
  * Show only blocks with the given name.
* `runme --shell bash README.md`
  * Show only scripts using the given shell.
  * **NOTE** To ease real-life usage we use "contains" here, rather than equals as we do for the name.
    * This means "bash" matches "/bin/bash", and "sh" will match both "/bin/bash" and "/bin/sh".



## Executing Blocks

To execute a matching block, or set of blocks, add `--run` to your command-line argument:

```
$ ./runme --shell /usr/bin/python3 --run README.md
/home/skx
```

As you might guess we work by writing the given block(s) to a temporary file, and executing it.

Add `--keep` to see the names of the temporary file(s) created.

You can also use `--join` to join the contents of all matching blocks to a single file.  For example
this document has two shell-blocks configured to use `/bin/bash`.  We can see these like so:

```
$ ./runme --shell bash README.md
Shell:/bin/bash  Name:test
uptime

Shell:/bin/bash  Name:whoami
id
```

If we wanted to we could run these commands as one shell-script, and keep the output we see the expected content:

```
$ ./runme --shell bash --join --run --keep README.md
wrote to /tmp/rm1802985270
 12:44:56 up 47 days, 18:36,  1 user,  load average: 0.27, 0.24, 0.19
uid=1000(skx) gid=1000(skx) groups=1000(skx),24(cdrom),25(floppy),27(sudo),29(audio),30(dip),44(video),46(plugdev),108(netdev),113(bluetooth),114(lpadmin),118(scanner),133(uml-net),999(docker)

$ cat /tmp/rm1802985270




#!/bin/bash

uptime

id
```
