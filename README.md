# jsonconf

A minimalistic JSON config reader.

No nasty struct definitions in client code, just one-line retrieving values by
keys. All numeric values in JSON configs are converted to int64 type.
JSON arrays are returned as []string and []int64.

Usage:


```
conf = jsonconf.ReadFile("config.json")
port, ok := conf("tcp.port", 80).(int)
...
host, ok := conf("tcp.host", "microsoft.com").(string)
...
clusters, ok := conf("clusters", []string{}{"default"}).([]string)
...
// or the same without ", ok" if you are sure about your config, e.g.:
name := conf("client.name", "noname")
```

Copyright (c) 2016 Eugene Korenevsky

License: MIT

