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

Also I just want to thank this blog article:
https://www.alexedwards.net/blog/working-with-cookies-in-go
