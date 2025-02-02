status page:
...endpoints...
/ui/status - view of all check status as timelines
/ui/checks - manage checks

status api:
...endpoints...
/api/checks
/api/checks/{ID}
/api/checks/{ID}/status
/api/checks/{ID}/check

status checker:
...

models:
check = {
  name; string
  command; string
  schedule; string
  enabled; string
  recoveryCommand; string
}

status = {
  ok; boolean
  logs; string
  recovered; boolean
}