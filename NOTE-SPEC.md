# ScheduleMe

## Overview

ScheduleMe is an open scheduling platform that allows guests to schedule events on a user's calendar. Users can set specific times of availability and create customizable event types. This functionality enables seamless and efficient planning, with users and guests being able to coordinate their schedules in real-time.

## Features

- User authentication (signup and login)
- Creation of event types with granularities of 15, 30, and 60 minutes
- Availability setting based on user preference (e.g., Monday-Friday, 9AM-5PM)
- Synchronization of availability with linked calendar
- Collection of email and a note during guest schedule

## Components

- **User Authentication**
  - OAuth signup and login
- **Event Creation and Management**
  - Creation of event types with specific titles and granularities
  - Ability to enable, disable, or delete created events
- **User Event View**
  - Private display of events for management
  - Editing functionality for created events
- **Guest Event View**
  - Public endpoint listing enabled events for a given user
- **Calendar Synchronization**
  - Fetch remote calendar and merge timeslots for a requested event type
- **Event Signup**
  - HTTP post guest email and time block
  - Creation of event in the user's calendar
  - Sending of confirmation email

## API Endpoints

- User Authentication
  - `/signup`
  - `/oauth-callback`
  - `/signin`
  - `/signout`
- Event Management
  - GET Single (Public/Private): `/{user_id}/events/{id}`
  - GET Index (Public/Private): `/{user_id}/events`
  - POST Create: `/{user_id}/events`
  - DELETE Delete: `/{user_id}/events/{id}`
- Time-slot Management
  - GET Time-slots: `/{user_id}/timeslots/{event_id}?{start=<s>}&{end=<e>}&{page=<p>}`
  - POST time-slots: `/{user_id}/timeslots/{event_id}`

## Development Tools

- Database migrations with `go-migrate`

## Tech Stack

- For the first iteration (MVP), the focus will be on functionality rather than aesthetic design. Only the essential UI components will be developed.
- We aim to stick to the Go standard library (template, SQL, HTTP), excluding complex OAuth implementation.
- CSS and JavaScript frameworks will not be used. CSS will be utilized sparingly and only when necessary.
