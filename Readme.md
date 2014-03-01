### Simple 'script' to count the file descriptors of a process.

It uses pidof to look up the process and then checks /prod/[pid]/fd to get a
count of open file descriptors. It runs as a go/compiled binary so it can have
the SUID bit set, which is required to access the (access controlled)
/proc/[pid]/ data.

The name of the processes is hardcoded so it doesn't use any input that would
get fed to the system calls w/ root permision. You can set this 'hardcoded'
value at compile time by setting the PROCESS environment variable to the name
of the process you want to monitor. It will set the constant variable in the
code and add it to the name of the binary and package (so you can monitor
multiple things on the same system).

Eg. to create a binary for watching nginx..
    $ env PROCESS=nginx make
To create a simple debian package..
    $ env PROCESS=nginx make deb

See the Makefile for more.
