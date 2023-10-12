# The Sis50 House Site

## What is this site?

This is a simple site created for the house Sis50 in the Netherlands to coordinate cooking schedules, share funny memes, and create a community shopping list. This site's backend is written in Go, and the design is focused on a mobile-first experience.

## Features

- [x] Weekly Cooking Table for each Weekday
    - [x] Has four options for what people have one that 
        - [x] You can switch between options simply by clicking
        - [x] May add annotations for extra things that should happen
            - [ ] Should be able to edit these annotations
            - [ ] Should be synced between users
        - [ ] Should be able to set specific meal (that can call up a recipe)
    - [ ] That resets at the end of every week
    - [ ] People have the ability to add small annotations to each day in case of special occurrences
    - [ ] A small recipe section where people can pull up recipes of what will be cooked today if wanted
    - [x] A next week calendar, where people can pre-plan / show when they are available next week
        - [ ] Automatic switching, so that the next-week schedule becomes the current weeks'
- [ ] Simple Shopping list 
    - [ ] Ability to add items to the shopping list
    - [ ] Easily delete items off the shopping list
    - [ ] Edit Items on the shopping list 
        - [ ] Have this edit feature be fancy
- [ ] Small Mini-game
- [ ] Daily Comic
    - [ ] Have the server decide on a daily comic (XKCD)
        - [ ] Transfer this comics title and link to the original
- [ ] Small Announcement Section
    - [ ] Max limit of 5 announcements, the oldest one gets discarded

- [ ] Input sanitization
- [ ] A dark mode theme with automatic detection as well
- [ ] A password-protected Admin panel
    - [ ] With the ability to reset respective weeks completely


- [ ] The ability for a user to change the background on their client side
    - [ ] This being saved
- [ ] Ability for user to be able to remove headers (to create unobstructed a background)
    - [ ] all headers should have an outline that semi-guarantee visibility on any background
- [ ] Ability for users to remove otherwise visible backgrounds from the main elements of the site

- [ ] Individual User authentication
    - [ ] That persists for a month
    - [ ] Only allows users to edit their respective column


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
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week VARCHAR(8) NOT NULL,
    person VARCHAR(16) NOT NULL,
    day VARCHAR(16) NOT NULL,
    content VARCHAR(512) NOT NULL,
    FOREIGN KEY (week, person, day) REFERENCES days(week, person, day)
);

```
