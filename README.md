# The Sis50 House Site

## What is this site?

This is a simple site created for the house Sis50 in the Netherlands to coordinate cooking schedules, share funny memes, and create a community shopping list. This site's backend is written in Go, and the design is focused on a mobile-first experience.

## Features

(I could have used GitHub Issues for this, but now I already wrote them down)

(In-between stages I will wait for user feedback and add requested features)

### **Feature Set** (Replacement Stage)
- [ ] A basic calendar for each person in the house.
- [ ] The current day being highlighted on the calendar
- [ ] Each calendar cell has 6 possible stages:
    - [ ] Present
    - [ ] Maybe present
    - [ ] Present but cannot cook
    - [ ] Maybe cooking
    - [ ] Cooking
    - [ ] Not there
- [ ] A basic shopping list with 3 options:
    - [ ] Remove
    - [ ] Edit
    - [ ] Add
- [ ] Notification / Notices Page
    - [ ] Ability to post new Announcement with bold highlighting

### **Feature Set** (User Customisation Stage)

- [ ] The ability to register for new users
- [ ] The ability to login to be a user
- [ ] The ability to post a message as that user
- [ ] A basic admin account
- [ ] A basic admin panel
- [ ] A basic help page
- [ ] The ability to set a background photo

### **Feature Set** (Improved Calendar Stage)
- [ ] There will be 2 calendars, one for the current week, another for the next week
- [ ] There will be dates associated with each calendar
- [ ] The next week calendar replaces the current week calendar once the weeks change
- [ ] Add the ability to add, edit and remove annotations from the calendar, on each cell
- [ ] More advanced admin panel controls for the calendar

- **Feature Set** (Recipe Stage)
- [ ]

- **Feature Set** (Minesweeper Stage)
- [ ]

- **Feature Set** (Pull up Leaderboard)
- [ ]

- **Feature Set** ("Avalex" Stage)
- [ ]

- **Feature Set** (Generic HTTP API Stage)
- [ ]

- **Feature Set** (Terminal Program Stage)
- [ ]

- **Feature Set** (Android Program Stage)
- [ ]


## Technical Information

### Interacting with Web sockets

> note that currently web sockets simply just send strings separated using $, this will eventually be replaced simply with JSON

structure of a web socket request is as follows:

COMMAND$further information...

Available commands and the further structure are:

- toggle
    - toggle$(state$)week$person$day
- addnote / deletenote
    - addnote$content$week$person$day
- post
    - post$content
// Going to be making a move over to JSON right about now, so these are the new definitions: (more for self-reference)

message:

```json
{
    command: ...,
    content: ...
}
```

toggle:

```json
{
    command: "toggle",
    currentState: ("E"|"O"|"X"|"?"),
    week: ("current"|"next"),
    person: ("milan"| "rick"| "robert"|"youri")
    day: ... all the days as listed above

}
```

### The Database


```sql
CREATE TABLE days (
    week VARCHAR(8) NOT NULL,
    person VARCHAR(16) NOT NULL,
    day VARCHAR(16) NOT NULL,
    state INTEGER NOT NULL,
    PRIMARY KEY (week, person, day)
);

CREATE TABLE berichte (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content VARCHAR(512)
);

CREATE TABLE notes (
    id INTEGER PRIMARY KEY,
    week VARCHAR(8) NOT NULL,
    person VARCHAR(16) NOT NULL,
    day VARCHAR(16) NOT NULL,
    content VARCHAR(512) NOT NULL,
    FOREIGN KEY (week, person, day) REFERENCES days(week, person, day)
);

```

Also I just want to thank this blog article:
https://www.alexedwards.net/blog/working-with-cookies-in-go
