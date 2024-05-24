# The Sis50 House Site

## Contents

- [What is this Site?](#What-is-this-Site)

- [Features](#Features)

- [Technical Information](#Technical-Information)

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
- [x] The ability to login to be a user
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

- [x] The calendar will be focused on one day, where the days will be shifting, with you being able to see the 2 previous days and 4 of the upcoming days.
- [ ] There will be dates associated with each calendar
- [ ] Add the ability to add, edit and remove annotations from the calendar, on each cell
- [ ] More advanced admin panel controls for the calendar

- **Feature Set** (Recipe Stage)
- [ ] Todo...

- **Feature Set** (Minesweeper Stage)
- [ ] Todo...

- **Feature Set** (Pull up Leaderboard)
- [ ] Todo...

- **Feature Set** ("Avalex" Stage)
- [ ] Todo...

- **Feature Set** (Generic HTTP API Stage)
- [ ] Todo...

- **Feature Set** (Terminal Program Stage)
- [ ] Todo...

- **Feature Set** (Android Program Stage)
- [ ] Todo...


## Technical Information

### Installing for Yourself

This section is a todo. You can probably figure it out, there is a make file and the first thing in the go code are the flags you can use to pass things for SSL.

Also I just want to thank this blog article:
https://www.alexedwards.net/blog/working-with-cookies-in-go
