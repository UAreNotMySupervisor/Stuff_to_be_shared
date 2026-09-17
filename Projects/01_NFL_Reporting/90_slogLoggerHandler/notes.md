# Notes about the custom slog Handler

The provided code is my implementation of a slog.Handler used to store slog.records persistently in my Postgresql database. 

## Goal of this Implementation

The function Logger_init is used to create a new slog.Logger with two Handles. 

One for default console-logging during development. 

The second custom handler is used to store logs above a certain log-level in my postgresql database.

The returned slog.Logger Pointer is passed to the rest of controllers/services/repostiories.

## Some Background on the code

- The LogRepo instance is created at program start and passed as parameter when initializing the Logger
- The Log-Level is extracted from a Config file and also passed along as parameter

- This only shows a small extract of my NFL Reporting Project
- As key functionality is still missing, it is not yet uploaded to this repo
- I used additional comments, as this was my first slog.Hanlder implementation and i wanted to easily remember it at later stages...

<br>

## Why this Source Code?

- This code is not very complex and that is its strength!
- It solves my problem very efficiently while staying maintainable, expandable and easy to understand for someone new working on this project.
- This is a core principle I try to follow: Keep it simple where possible!