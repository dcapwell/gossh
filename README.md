Gossh
====

Go library and CLIs for running ssh and scp concurrently over many hosts.

For the CLIs, arguments try to stick as close as possible to their linux versions.

```
gossh host-[1..100].example.com \
  -n 40 \                         # 40 concurrent ssh sessions at a time
  -l techops \                    # ssh as user techops
  -i /home/techops/.ssh/id_dsa \  # with identity from path
  -t \                            # add a tty session
  -o ServerAliveInterval=10 \     # options to pass to ssh
  -o ConnectTimeout=2 \
  -o UserKnownHostsFile=/dev/null \
  -o StrictHostKeyChecking=no \
  date                            # command to run on all hosts
```

```
goscp host-[1..100].example.com \
  -n 40 \                         # 40 concurrent ssh sessions at a time
  -l techops \                    # ssh as user techops
  -i /home/techops/.ssh/id_dsa \  # with identity from path
  -o ServerAliveInterval=10 \     # options to pass to ssh
  -o ConnectTimeout=2 \
  -o UserKnownHostsFile=/dev/null \
  -o StrictHostKeyChecking=no \
  /etc/hosts                      # local src file to copy over
  /tmp/admins-hosts               # location where to save file on remote hosts
```
