# twittermap

```
[batch]
                    tweet       S3
                      |          ^
                      v          | json
                    |lambda -> lambda |
cloudwatch alarm -> |  step function  |

[web]
s3 html -> S3 json
```
