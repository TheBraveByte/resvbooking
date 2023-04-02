# Reservation and Booking Project

This repository contains the source code for a reservation and booking project.

- Built using Go version 1.17
- The [chi router](https://github.com/go-chi/chi) 
- The [SCS](https://github.com/alexedwards/scs/v2) session management package
- [Nosurf](https://github.com/justinas/nosurf) package

## Features

The project has the following features:

- User authentication and session management
- User registration and login
- Ability to search and filter available reservations
- Ability to make a reservation
- Ability to view and manage reservations
- Protection against CSRF attacks using  [Nosurf](https://github.com/justinas/nosurf) package
