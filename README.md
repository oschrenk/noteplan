# README

Unofficial companion app for [NotePlan](https://noteplan.co/)

## Usage

Print out today's task count

```
noteplan todo
Day, 2023-10-10, Open: 2
Day, 2023-10-10, Closed: 1
Week, 2023-W41, Open: 5
Week, 2023-W41, Closed: 3
```

Print out any date's task count

```
noteplan todo "last friday"
Day, 2023-10-06, Open: 2
Day, 2023-10-06, Closed: 1
Week, 2023-W40, Open: 6
Week, 2023-W40, Closed: 2
```

Print out any date's task count, only on the day

```
noteplan todo --day-only "last friday"
Day, 2023-10-06, Open: 2
Day, 2023-10-06, Closed: 1
```

Print out any date's task count, only during that week

```
noteplan todo --week-only "last friday"
Week, 2023-W40, Open: 6
Week, 2023-W40, Closed: 2
```

Print as json

```
noteplan todo --json "last friday"
{
  "day": {
    "iso": "2023-10-06",
    "open": 2,
    "closed": 1
  },
  "week": {
    "iso": "2023-W40",
    "open": 6,
    "closed": 2
  }
}
```

Fail on missing entries

```
noteplan todo --fail-fast "next year"
2024/01/10 14:52:44 sql: no rows in result set
```

## Installation

**Via Github**

```
git clone git@github.com:oschrenk/noteplan.git
cd noteplan

# installs to $GOBIN/noteplan
task install
```

## Known Issues

### "[App]" would like to access data from other apps" warning

`noteplan` accesses NotePlan3's "internal" data at `$HOME/Library/Containers/co.noteplan.NotePlan3/Data/Library/Application Support/co.noteplan.NotePlan3/Caches/note-cache.db` and triggers a Privacy & Security alert. When called from the terminal it might say "Alacritty would like to access data from other apps" warning".

This will happen every time you start a new process. A dialog will pop up asking for access. You can give your terminal "Systems Settings > Privacy & Security > Full Disk Access" to circumvent the warning, but it does mean giving broad access also to other tools invoked through the terminal.

I'm still trying to figure out a better way to access specific application data.

#### Background

This is because we are trying to access another Application's container, and we are crossinga a TCC (Transparenency, Consent and Control) boundary).

> TCC privileges fall into various groups, including:
>
> 1. Persistent, system wide
> 2. Persistent, per user
> 3. Transient, process lifetime
>    The privilege to access other apps containers falls into that third group.

See also [here](https://developer.apple.com/forums/thread/742147#:~:text=Yes.%20TCC%20privileges%20fall%20into,falls%20into%20that%20third%20group.)

Since the executing process is the terminal, we can leverage the permissions of the process and give it Full Disk Access.

#### sketchybar

In case of [sketchybar](https://github.com/FelixKratz/SketchyBar), we could give it full disk access, but since it's a long running process and will only ask once when you start or restart the service, the impact is not as big, if you leave it off (recommended) and press "Allow" once when sketchybar starts.
