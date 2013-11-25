Gossh Http Agent
================

This package will create a binary Http Agent that can be used to interact with gossh.

Examples

Run cmd on multiple hosts
```
$ curl -s 'http://localhost:7654/ssh/localhost,localhost?cmd=echo+Hello+World' | python -mjson.tool
{
    "Responses": [
        {
            "Duration": "88ms", 
            "Hostname": "localhost", 
            "Response": {
                "Code": 0, 
                "Stderr": "", 
                "Stdout": "Hello World\n"
            }
        }, 
        {
            "Duration": "88ms", 
            "Hostname": "localhost", 
            "Response": {
                "Code": 0, 
                "Stderr": "", 
                "Stdout": "Hello World\n"
            }
        }
    ]
}

```

Failure to login
```
$ curl -s 'http://localhost:7654/ssh/localhost?cmd=echo+Hello+World&user=root&option=PreferredAuthentications=publickey' | python -mjson.tool
{
    "Responses": [
        {
            "Duration": "56ms", 
            "Hostname": "localhost", 
            "Response": {
                "Code": 255, 
                "Stderr": "Permission denied (publickey,keyboard-interactive).\r\n", 
                "Stdout": ""
            }
        }
    ]
}

```
