# Tasks

## Service Level

Cron: cron info parsing

- [x] Read crontab content
- [x] Parse content into fields
- [ ] Translate each tasks to human readable
- [x] Index tasks
- [x] Get running server hostname
- [x] Get machine timezone

New

- [x] Support task change detection
- [ ] Delete task report once task has changed, treat it as a new task

Commando: Task executor

- [x] Command execution
- [x] Update status
- [x] Update result
- [x] Go func implements parallel
- [x] Status and Result query
- [ ] Show log?


### task execution solution

How to store the result? Use a file aggregation to represent execution report.
- Use `/tmp/goat/[id]/status` to store execution status
  - FSM: created -> running -> done, one direction
- Use `/tmp/goat/[id]/lastresult` to store execution result of last run, default to `Not Run` when first created
  - FSM: Not Run (first created) -> Pass or Fail, bi-direction because the task could run several times.
- Use `/tmp/goat/[id]/lastlog` to store execution log of last run, do not support query, only for debugging.



How to Run the command?
- Got command to run, set status to `created`
- Start a go func to run the command, and set the status to `running`
- Go func completed the command, update the status to `done`

### Limitation

This solution might return wrong info or run wrong job when cron jobs has any changes, for example, 
```
03 14 * * * python3 -u /home/example/thg/scripts/po_uid.py >> /tmp/po_uid.log 2>&1
#01 08 * * * python3 -u /home/example/thg/scripts/po_sd1.py >> /tmp/po_sd1.log 2>&1
04 14 * * * python3 -u /home/example/thg/scripts/po_ad2.py >> /tmp/po_ad2.log 2>&1
```
changes to 
```
03 14 * * * python3 -u /home/example/thg/scripts/po_uid.py >> /tmp/po_uid.log 2>&1
01 08 * * * python3 -u /home/example/thg/scripts/po_sd1.py >> /tmp/po_sd1.log 2>&1
04 14 * * * python3 -u /home/example/thg/scripts/po_ad2.py >> /tmp/po_ad2.log 2>&1
```
For viewing task info scenario,

The index will be changed. In first text, the task with index=2 is `po_ad2`, and in the later text, the task with index=2 is `po_sd1`. At this point, `po_sd1` takes the result of previous `po_ad2`.

For task execution scenario,

The task a user want to run is `po_ad2` but after the change, the actual command will run is `po_sd1`. And this is not the intention.

## API Level
- [x] View existing cron jobs in fields -> Cron
- [x] The view should return the latest content of crontab, this means n calls to open cron file to response n API calls
- [x] Trigger task execution -> Commando
- [x] View last triggered job result, pass/fail -> Commando


## Security Consideration
1. How to authenticate and authorize?
2. TLS support?