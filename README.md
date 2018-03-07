# twittermap

```
[batch]
                    tweet
                      |
                      v
                    |lambda -> lambda | -> S3 json
cloudwatch alarm -> |  step function  |

[web]
s3 html -> S3 json
```
