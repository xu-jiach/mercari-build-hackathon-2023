# Mini Mercari Web App

## Our goal is to offer these experiences to users:

  1. Basic mercari features 
  2. Powerful auto-writing with GPT
  3. Let's mercari in person "Face2Pay"

Here are the demos for the features

## Log in / Sign up

https://github.com/xu-jiach/mercari-build-hackathon-2023/assets/90857923/dd4cb806-3f63-4259-a815-41596b48f4c3

## Browsing, Filtering, Searching

https://github.com/xu-jiach/mercari-build-hackathon-2023/assets/90857923/667dbcbe-11fd-45c5-bb50-61337cf9bb99

## Listing

https://github.com/xu-jiach/mercari-build-hackathon-2023/assets/90857923/57855b26-ebb2-451b-bef9-de8a4fc8e822

If you register the passcode for in-person purchase, you can still check the passcode for the item! (Of course others can't!)

## In-person purchase "Face2Pay"
### 1. Owner: enter a passcode when you list the item

https://github.com/xu-jiach/mercari-build-hackathon-2023/assets/90857923/9912eb95-d1ed-4fc9-9131-a96b1150f49c

### 2. Buyer: check "I'm with the owner!" and enter the passcode that the only seller knows
The "I'm with the owner!" shows up only for the items that the seller allowed the in-person buying

https://github.com/xu-jiach/mercari-build-hackathon-2023/assets/90857923/3778c385-a8b9-46fc-ba66-f58e21f8131e

## Transaction History

https://github.com/xu-jiach/mercari-build-hackathon-2023/assets/90857923/4f63cd55-3c0e-4214-b465-e5402abbf39f



-----------

## Requirements

* [go](https://go.dev/)

## Getting started

### 1. Update environment values

1. Create `initialize` branch.

2. Run bellow command.

```shell
$ go run tools/setup.go -g [your github name] -t [your team id]

(e.g.) $ go run tools/setup.go -g yourname -t 16
```

3. Create a PR form `initialize` to `main`.

### 2. Launch services

See `backend/README.md` for backend service and `frontend/simple-mercari-web/README.md` for frontend service.

## How to run bench marker

* Open [dashboard](https://mercari-build-hackathon-2023-front-d3sqdyhc4a-uc.a.run.app/)
* Tap `RUN BENCHMARK` Button!

## What should we do first?

- First, stand up services and see logs both of backend and frontend services
  - For backend, you can see the logs on the terminal where the server is set up
  - For frontend, use Chrome Devtool to check
    - https://developer.chrome.com/docs/devtools/overview/
- Try to use mini Mercari and find problems that should not occur in the original Mercari service. For example:
  - When you check the item detail page of your listed items...?
  - When you try to buy items that exceed your available balance...?
  - When there are multiple users purchase an item at the same time...?
- The UI is quite simple and difficult to use
  - It looks inconvenient if there is no message indicating a request to the backend has failed
  - As the number of items increase, the UI is likely to become slow
  - Feel free to make the site more user-friendly like implementing features not implemented in actual Mercari web site

However, your changes must be made within the constraints of the bench marker. Please refer to backend/README.md for details.
