

*This readme is in WIP.*

## SHOULD READ FIRST!

This project is to fast demo features we expect with temporal and PRDs.

One key current assumption/design for this repo is that we base approval domain totally on temporal sdk, no other middle layer is used.
This may not be the final version, but a good start to test, for this is how temporal examples are organized and layered, and this is simple to get ideas of Temporal.
In this way, approvals can have config to change from the Frontend user's website, but can hardly draw flows freely which need code. When we delve more deep into temporal, another repo will probably be set as a demo for this different design.


## KEY PRD FEATURES

To find more in this link:

https://airwallex.atlassian.net/wiki/spaces/~625423894f1d57006a23dd99/pages/2619679918/WIP+Temporal+Startup+Kit



## Quick Experiment with commandlines

- install golang1.17
- docker-compose up the temporal platform
- cd app/worker;  go run main.go
- cd app/cmdstart;  go run main.go  (copy the workflowID of stdout)
- cd app/cmdquery; go run main.go -workflowID 3154512a-c925-44de-8309-bf215423f67c
- cd app/cmdsubmit; go run main.go -workflowID 3154512a-c925-44de-8309-bf215423f67c


## How to see temporal data

- expose the ports in docker-compose
- key engine data in pgsql: pgcli postgres://temporal:temporal@localhost:5432
- self-defined qureis from es: http://localhost:1358/?appname=*&url=http://localhost:9200&mode=edit


## Temporal Investigation

https://airwallex.atlassian.net/wiki/spaces/~625423894f1d57006a23dd99/pages/2631632801/Temporal+Investigation
