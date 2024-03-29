# The Sis50 House Site

## What is this site?

This is a simple site created for the house Sis50 in the Netherlands to coordinate cooking schedules, share funny memes, and create a community shopping list. This site's backend is written in Go, and the design is focused on a mobile-first experience.

## Features

(I could have used GitHub Issues for this, but now I already wrote them down)

(In-between stages I will wait for user feedback and add requested features)

### **Feature Set** (Replacement Stage)
- [x] A basic calendar for each person in the house.
- [x] The current day being highlighted on the calendar
- [x] Each calendar cell has 6 possible stages:
    - [x] Present
    - [x] Maybe present
    - [x] Present but cannot cook
    - [x] Maybe cooking
    - [x] Cooking
    - [x] Not there
- [x] Calendar will be saved
- [x] A basic shopping list with 3 options:
    - [x] Remove
    - [x] Edit
    - [x] Add
- [x] Notification / Notices Page
    - [x] Ability to post new Announcement 
    - [x] Ability to have bold highlighting
    - [x] Ability to have italics
- [x] Addition of a basic favicon for the website
~- [ ] A link to the local supermarket's site~
- [x] Add mobile design (basic)
- [x] There is a basic help page

### **Feature Set** (User Customisation Stage)

- [ ] The ability to register for new users
- [ ] The ability to login to be a user
- [ ] The ability to post a message as that user
- [ ] A basic admin account
- [ ] A basic admin panel
- [x] A basic help page
- [x] The ability to set a background photo
- [ ] The ability to modify the basic colours of the page

### **Feature Set** (Shopping Streamling Stage)
- [ ] CSV Download for Shopping and when it happened
... other things too

### **Feature Set** (Improved Calendar Stage)
- [ ] The calendar will be focused on one day, where the days will be shifting, with you being able to see the 2 previous days and 4 of the upcoming days.
- [ ] There will be dates associated with each calendar
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
