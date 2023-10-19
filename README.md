## Artifacts upload and Download Report 
Get Report of Artifacts upload and Download in a specified amount of days. 
#### Script Inputs: 
Update input.json file. (Dont change the name of the file)
* Access Token.   
* Jfrog Base url.
* Number of Days.

**Example:**:
```
{
    "token": "<Your Token>",
    "jfrogbaseurl": "https://myart.jfrog.io",
    "odays": "30d"
}
```

#### Create a build 
`go build -o reportgen`

#### Script Execution: 
`./reportgen` 

#### Script Output Example: 
Below is the output 
**Output:** 
```
Report Based on Repo
+--------------------------------+---------------+----------------+-----------------+--------------+
|              NAME              | TOTALDOWNSIZE | TOTALDOWNLOADS | TOTALUPLOADSIZE | TOTALUPLOADS |
+--------------------------------+---------------+----------------+-----------------+--------------+
| Total                          |      60858898 |            559 |       771859527 |         2064 |
| nagag-ldocker                  |             0 |              0 |       629084916 |           30 |
| maven-remote-cache             |      56030992 |            510 |        21725246 |          510 |
| sup002-swampup-maven-dev-local |             0 |              0 |        68854914 |            3 |
| sup002-swampup-npm-dev-local   |             0 |              0 |        22479783 |            5 |
| jfrogpipelines                 |       1475802 |              3 |        11352572 |           76 |
| dev-docker-local               |             0 |              0 |         9983219 |           14 |
| maven-snap-local               |             0 |              0 |         5147706 |           26 |
| nagag-maven                    |         15888 |              5 |         2622741 |           24 |
| jfrog-automation               |       2375706 |              4 |             110 |            1 |
| nagag-terraform-remote-cache   |        959622 |             36 |          399374 |           39 |
| dev-npm-local                  |             0 |              0 |           81014 |           12 |
| artifactory-pipe-info          |           888 |              1 |           66258 |          150 |
| nagag-terraformbe              |             0 |              0 |           42137 |           18 |
| artifactory-build-info         |             0 |              0 |           19537 |            2 |
| jfrog-usage-logs               |             0 |              0 |               0 |         1154 |
+--------------------------------+---------------+----------------+-----------------+--------------+
Report Based on Users
+---------------------------------------+---------------+----------------+-----------------+--------------+
|                 NAME                  | TOTALDOWNSIZE | TOTALDOWNLOADS | TOTALUPLOADSIZE | TOTALUPLOADS |
+---------------------------------------+---------------+----------------+-----------------+--------------+
| Total                                 |      60858898 |            559 |       771859527 |          896 |
| admin                                 |             0 |              0 |       629084749 |           29 |
| fabienl@jfrog.com                     |             0 |              0 |        91334697 |            6 |
| workshop_admin                        |      59382208 |            555 |        30038032 |          622 |
| token:jfpip@ZHNM7EVEA0E               |        111410 |              3 |        11093785 |          199 |
| token:jfds@01g8xgwtgt8kv211h54hjw03nx |             0 |              0 |         9983219 |           13 |
| token:jfpip@IfUydYUCLxM               |       1365280 |              1 |               0 |            0 |
| token:jfpip@mk6DryW9iPY               |             0 |              0 |          325045 |           27 |
+---------------------------------------+---------------+----------------+-----------------+--------------+
```
