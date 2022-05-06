# goodsman2


- Please create `{repo}/config.yml` according model below.

- Set EncryptKey `""` if you do not need secure verification

``` yml
Base:
  HttpPort: 1926

App:
  AppID: "{Set your app id here.}"
  AppSecret: "{Set your app password here.}"
  EncryptKey: "{Set your encrypt key here}"

Mongo: 
  DBName: "goodsman2"
  Url: "{Set your DB connect url}"
```