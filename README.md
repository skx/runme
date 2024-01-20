# runme - A simple markdown extractor

This is a simple project which allows you to list, extract, and execute commands from markdown files.

Given the following code-block you'll see three things:

```/bin/bash test
uptime
```

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



## Listing Blocks

By default all blocks found in the README.md file will be shown, for example:

```bash
$ ./runme
Shell:/bin/bash  Name:test
uptime

Shell:/usr/bin/python3  Name:home
import os
print(os.environ['HOME'])

..
```

Note here that only blocks with **both** a shell and a name are listed?

You can filter the output to only show blocks with a given shell, or name via the flags:

* `runme --name test`
  * Show only blocks with the given name.
* `runme --shell bash`
  * Show only bash-scripts.



## Executing Blocks

To execute a matching block, or set of blocks, add `--run` to your command-line argument:

```
$ ./runme --shell /usr/bin/python3 --run
/home/skx
```

As you might guess we work by writing the given block(s) to a temporary file, and executing it.

Add `--keep` to see the names of the temporary file(s) created.
