# Resolved - B2B ticketing system

This is the repository for my first Go project.

It's my playground for learning Go. I decided to build a B2B ticketing system as I am working with similar tool in my job. The goal is to have a fully functional application where customers can create ticket both via Email and client-dedicated website (Client Portal). Internal users can also update/create tickets both via Email or internal app. 

I am using PostgreSQL as my database, JQuery and clean JS for the frontend, but the goal is to change it to React/Vue/Angular, but I am still deciding which framework I want to learn. 

App will be dockerized. Ideally I want to include multiple microservices (email processor, automations, SLAs).


## Current functionalities of the app:
* Creating users, groups, tickets, products, customers, contacts
* Associating Tickets with Group, Contact, Customer and Product
* Adding Users to Groups
* Displaying only Tickets from User's Groups or all Tickets
* Associating Contacts with Customers
* Managing User permissions - grant/remove ability to Edit/Create other Objects
* User authorization with JWT token
* Custom Fields - fields defined by admin users that are displayed in tickets. Each value is ID based. Custom fields values are stored in JSONB field in PostgreSQL database.


## To Do:
* Change frontend to React/Vue/Angular (still deciding)
* Dockerize the app
* Remove unused parts of the code
* Password hashing
* Forms validation
* API (both REST and GraphQL for learning purpose)
* Email handling (still figuring out how to handle it - thinking about connecting to user mailbox via IMAP and processing received emails, but might end up using something else)
* Add ticket automation
* Assign user to group
* Email notification about updates in tickets
* Improve permission check - I think that the current solution isn't the best
* Improve project structure - can't store everything in the handlers file
* Remove redundancy - a lot of my functions do the same thing, it can be improved. For example DB connection or multiple SQL queries.
* Add search function
* Add chat function - WebSocket?
* Add ticket types
* Ticket views defined by ticket type (different custom fields, different pick list options based on ticket type)
