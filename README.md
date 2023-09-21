# The Sis50 House Site

## What is this site?

This is a simple site created for the house Sis50 in the netherlands to coordinate cooking schedules, share funny memes, and create a community shopping list. This site's backend is written in Go, and the design is focused on a mobile-first experience.

## Features

- [x] Weekly Cooking Table for each Weekday
    - [x] Has four options for what people have one that 
        - [x] You can switch between options simply by clicking
        - [x] May add annotations for extra things that should happen
            - [ ] Should be able to edit these annotations
            - [ ] Should be synced between users
        - [ ] Should be able to set specific meal (that can call up a recipe)
    - [ ] That resets at the end of every week
    - [ ] People have the ability to add small annotations to each day in case of special occurences
    - [ ] A small recipe section where people can pull up recipes of what will be cooked today if wanted
    - [x] A next week calendar, where people can pre-plan / show when they are available next week
        - [ ] Automatic switching, so that the next-week schedule becomes the current weeks'
- [ ] Simple Shopping list 
    - [ ] Ability to add items to the shopping list
    - [ ] Easily delete items off the shopping list
    - [ ] Edit Items on the shopping list 
        - [ ] Have this edit feature be fancy
- [ ] Small Minigame
- [ ] Daily Comic
    - [ ] Have the server decide on a daily comic (XKCD)
        - [ ] Transfer this comics title and link to the original
- [ ] Small Announcement Section
    - [ ] Max limit of 5 announcements, oldest one gets discarded

- [ ] Input sanitisation
- [ ] A dark mode theme with automatic detection as well
- [ ] A password-protected Admin panel
    - [ ] With the ability to reset respective weeks completely


- [ ] The ability for a user to change the background on their client side
    - [ ] This being saved

- [ ] Individual User authentication
    - [ ] That persists for a month
    - [ ] Only allows users to edit their respective column


## Technical Information

### Interacting with Websockets

> note that currently websockets simply just send strings separated using $, this will eventually be replaced simply with JSON

structure of a websocket request is as follows:

COMMAND$further information...

available commands and the further structure are:

- toggle
    - toggle$(state$)week$person$day
- addnote / deletenote
    - addnote$content$week$person$day
- post
    - post$content
// gonna be making a move over to JSON right about now, so these are the new defintions: (more for self-reference)

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

